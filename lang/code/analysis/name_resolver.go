package analysis

type NameResolver struct {
	Procedures        map[string]*Procedure
	ComponentTypesMap map[string]string // Button1 -> Button
	ComponentNameMap  map[string][]string
}

type Procedure struct {
	Name       string
	Parameters []string
	Returning  bool
}
