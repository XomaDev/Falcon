package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pion/webrtc/v3"
)

const RENDEZVOUS = "https://rendezvous.appinventor.mit.edu/rendezvous/"

const INIT_FORM = "(begin (require <com.google.youngandroid.runtime>) (process-repl-input -1 (begin (define-syntax protect-enum   (lambda (x)     (syntax-case x ()       ((_ enum-value number-value)         (if (< com.google.appinventor.components.common.YaVersion:BLOCKS_LANGUAGE_VERSION 34)           #'number-value           #'enum-value)))))(clear-current-form))))"

func main() {
	sixDigitCode := "g"
	hasher := sha1.New()
	hasher.Write([]byte(sixDigitCode))
	hashedBytes := hasher.Sum(nil)

	rendezvousKey := hex.EncodeToString(hashedBytes)
	pollUrl := RENDEZVOUS + rendezvousKey

	for {
		resp, err := http.Get(pollUrl)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode == 200 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				resp.Body.Close()
				panic(err)
			}
			resp.Body.Close()

			if len(body) > 0 {
				fmt.Println("Companion checked in! Received payload:")
				whenCompanionAvailable(sixDigitCode, rendezvousKey, body)
				break
			} else {
				fmt.Println("Got 200 OK, but no data yet. Retrying...")
			}

		} else {
			fmt.Printf("Awaiting companion... (Status: %s)\n", resp.Status)
			resp.Body.Close()
		}
		time.Sleep(2 * time.Second)
	}

}

// Root structure for the JSON
type Config struct {
	IceServers  []IceServer `json:"iceservers"`
	Rendezvous2 string      `json:"rendezvous2"`
	IPAddr      string      `json:"ipaddr"`
}

// Nested structure for the ICE servers array
type IceServer struct {
	Server   string `json:"server"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func whenCompanionAvailable(key string, rendezvousKey string, response []byte) {
	var cfg Config
	if err := json.Unmarshal(response, &cfg); err != nil {
		panic(err)
	}
	var iceServers []webrtc.ICEServer
	for _, serverInfo := range cfg.IceServers {
		iceServers = append(iceServers, webrtc.ICEServer{
			URLs:       []string{serverInfo.Server},
			Username:   serverInfo.Username,
			Credential: serverInfo.Password,
		})
	}
	config := webrtc.Configuration{ICEServers: iceServers}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		panic(err)
	}

	dataChannel.OnOpen(func() {
		fmt.Printf("Data channel has opened!\n")
	})

	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Got datachannel message: %s\n", msg.Data)
	})

	// rendezvousKey is the SHA1 hash from the previous stage
	signalingKey := rendezvousKey + "-s"
	println("signalling key: " + signalingKey)
	rendezvous2URL := cfg.Rendezvous2 // The URL from the initial response

	nonce := rand.Intn(10000) + 1
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return // Indicates candidate gathering is complete
		}

		// Construct the JSON payload to send
		payload := map[string]interface{}{
			"key":       signalingKey,
			"webrtc":    true,
			"nonce":     nonce, // A random number to prevent replays
			"candidate": c.ToJSON(),
		}

		// Convert payload to JSON bytes
		jsonBytes, _ := json.Marshal(payload)

		// POST jsonBytes to rendezvous2URL
		// ... (your http.Post implementation here) ...
		fmt.Printf("Sent ICE Candidate: %s\n", string(jsonBytes))
	})

	// 1. Create the Offer
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	// 2. Set Local Description (and start gathering ICE candidates)
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	// 3. Construct the JSON payload for the offer
	payload := map[string]interface{}{
		"key":       key + "-s", // The key ending in "-s"
		"webrtc":    true,
		"offer":     offer,
		"nonce":     rand.Intn(10000) + 1,
		"candidate": nil,
	}
	jsonBytes, _ := json.Marshal(payload)

	println("Rendezvous posting offer")
	println(rendezvous2URL)
	println(string(jsonBytes))
	// 4. POST the offer to the signaling server
	resp, err := http.Post(rendezvous2URL, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		panic(err)
	}
	println(resp.StatusCode)

	fmt.Println("Sent SDP Offer. Now waiting for Answer...")

	receiverURL := rendezvous2URL + key + "-r"
	startPoll(peerConnection, receiverURL)

	// after channel opened

	// --- Assume 'dataChannel' is your *webrtc.DataChannel object ---
	// You can do this right inside the OnOpen callback for a first test,
	// or anywhere in your code after the OnOpen event has fired.

	// A simple YAIL command to show a toast message on the screen.
	// Note the use of backslashes to escape the quotes inside the string.

	dataChannel.SendText(INIT_FORM)

	yailCode := "(begin\n  (require <android.util.Log>)\n  (android.util.Log:i \"YAIL\" \"fuck me cat\"))\n"

	fmt.Printf("Sending YAIL: %s\n", yailCode)

	// Send the YAIL string over the data channel
	err = dataChannel.SendText(yailCode)
	if err != nil {
		panic(err)
	}

	// Set this callback up when you first create your peerConnection object.
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

		switch s {
		case webrtc.PeerConnectionStateFailed:
			// The connection has failed unexpectedly. It's time to clean up.
			fmt.Println("Connection Failed! Closing peer connection.")
			// You may want to attempt to restart the connection process here.
			peerConnection.Close()

		case webrtc.PeerConnectionStateDisconnected:
			// This is a temporary state. The connection might come back.
			// For example, if the user walks through a Wi-Fi dead spot.
			// You can wait a few seconds here to see if it recovers to "Connected".
			// If it stays disconnected for too long, it will eventually move to "Failed".
			fmt.Println("Connection Disconnected. Waiting to see if it recovers...")

		case webrtc.PeerConnectionStateClosed:
			// This means the connection was closed gracefully, either by you
			// calling peerConnection.Close() or the remote peer closing it.
			fmt.Println("Connection Closed. Cleanup complete.")
		}
	})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func startPoll(peerConnection *webrtc.PeerConnection, receiverURL string) {

	// --- Assume these are already defined from the previous steps ---
	// peerConnection *webrtc.PeerConnection
	// rendezvousKey string
	// response      YourInitialResponseStruct // Contains rendezvous2 URL
	// -------------------------------------------------------------

	// Construct the receiver URL

	// ... (imports and setup from before) ...

	// --- NEW: Create a buffer for pending candidates ---
	var pendingCandidates []webrtc.ICECandidateInit

	// --- This flag is the key to the logic ---
	haveSetAnswer := false

	// Poll for a maximum of ~60 seconds
	for i := 0; i < 60; i++ {
		if peerConnection.ConnectionState() == webrtc.PeerConnectionStateConnected {
			fmt.Println("PeerConnection is CONNECTED!")
			break
		}
		fmt.Printf("Polling companion... Attempt %d/60 (State: %s)\n", i+1, peerConnection.ConnectionState().String())

		resp, err := http.Get(receiverURL)
		if err != nil {
			fmt.Printf("Polling error: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if resp.StatusCode == 200 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				resp.Body.Close()
				panic(err)
			}
			resp.Body.Close()

			if len(body) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}

			var hunks []map[string]interface{}
			if err := json.Unmarshal(body, &hunks); err != nil {
				fmt.Printf("Error unmarshalling hunks: %v\n", err)
				continue
			}

			for _, hunk := range hunks {
				// ... (nonce check) ...

				// Check for the SDP Answer
				if answerData, ok := hunk["offer"].(map[string]interface{}); ok && !haveSetAnswer {
					fmt.Println("Received SDP Answer!")
					answerBytes, _ := json.Marshal(answerData)
					var answer webrtc.SessionDescription
					json.Unmarshal(answerBytes, &answer)

					// Set the remote description
					if err := peerConnection.SetRemoteDescription(answer); err != nil {
						panic(err) // This is a fatal error if it fails
					}
					haveSetAnswer = true // --- SET THE FLAG ---

					// --- NEW: Process any buffered candidates ---
					fmt.Printf("Processing %d pending ICE candidates...\n", len(pendingCandidates))
					for _, pendingCandidate := range pendingCandidates {
						if err := peerConnection.AddICECandidate(pendingCandidate); err != nil {
							panic(err)
						}
					}
					// Clear the buffer
					pendingCandidates = nil

				}

				// Check for an ICE Candidate
				if candidateData, ok := hunk["candidate"].(map[string]interface{}); ok {
					candidateBytes, _ := json.Marshal(candidateData)
					var candidate webrtc.ICECandidateInit
					json.Unmarshal(candidateBytes, &candidate)

					// --- NEW: Check the flag before adding ---
					if !haveSetAnswer {
						fmt.Println("Received an early ICE candidate. Buffering it.")
						pendingCandidates = append(pendingCandidates, candidate)
					} else {
						fmt.Println("Received ICE Candidate. Adding it now.")
						if err := peerConnection.AddICECandidate(candidate); err != nil {
							panic(err)
						}
					}
				}
			}
		} else {
			resp.Body.Close()
		}

		time.Sleep(1 * time.Second)
	}

	// ... (final connection check) ...
}
