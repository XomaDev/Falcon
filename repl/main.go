package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pion/webrtc/v3"
)

func main() {
	var code string
	fmt.Print("Enter code: ")
	fmt.Scan(&code)

	repl := NewRepl(code, DefaultRendezvous, 60, onConnect, onDisconnect, onMessageReceived)
	if err := repl.Connect(); err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

func onDisconnect(graceful bool) {
	// TODO: We gotta do nothing for now
	fmt.Println("Companion disconnected.")
}

func onConnect(c *webrtc.DataChannel) {
	fmt.Println("Companion connected!")
	// TODO:
	//  Later link up a file and lookup for changes and update it
	testYail := "(begin (require <com.google.youngandroid.runtime>) (process-repl-input -1 (begin (define-syntax protect-enum   (lambda (x)     (syntax-case x ()       ((_ enum-value number-value)         (if (< com.google.appinventor.components.common.YaVersion:BLOCKS_LANGUAGE_VERSION 34)           #'number-value           #'enum-value)))))(clear-current-form)))) (begin (require <com.google.youngandroid.runtime>) (process-repl-input -1 (begin (try-catch (let ((attempt (delay (set-form-name \"Screen1\")))) (force attempt)) (exception java.lang.Throwable 'notfound))))) (begin (require <com.google.youngandroid.runtime>) (process-repl-input -1 (begin (do-after-form-creation (set-and-coerce-property! 'Screen1 'ActionBar #t 'boolean) (set-and-coerce-property! 'Screen1 'AppName \"hahah\" 'text) (set-and-coerce-property! 'Screen1 'ScreenOrientation \"unspecified\" 'text) (set-and-coerce-property! 'Screen1 'ShowListsAsJson #t 'boolean) (set-and-coerce-property! 'Screen1 'Sizing \"Responsive\" 'text) (set-and-coerce-property! 'Screen1 'Title \"Screen1\" 'text)) (add-component Screen1 com.google.appinventor.components.runtime.Button Button1 (set-and-coerce-property! 'Button1 'Text \"Text for Button1\" 'text) ) (add-component Screen1 com.google.appinventor.components.runtime.Button Button2 (set-and-coerce-property! 'Button2 'Text \"Text for Button2\" 'text) ) ))) (begin (require <com.google.youngandroid.runtime>) (process-repl-input -1 (begin (init-runtime)))) (begin (require <com.google.youngandroid.runtime>) (process-repl-input -1 (begin (call-Initialize-of-components 'Screen1 'Button1 'Button2))))"
	c.SendText(testYail)
}

func onMessageReceived(msg webrtc.DataChannelMessage) {
	fmt.Printf("Message received: %s\n", msg.Data)
}
