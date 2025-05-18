package analysis

type NameResolver struct {
	Procedures map[string]Procedure
}

type Procedure struct {
	Name       string
	Parameters []string
	Returning  bool
}
