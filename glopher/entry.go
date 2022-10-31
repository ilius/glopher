package glopher

type Entry struct {
	Error      error
	Word       string
	Defi       string
	AltWord    []string
	AltDefi    []string
	DefiFormat DefiFormat
	IsInfo     bool
}
