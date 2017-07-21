package glopher

type Entry struct {
	Word       string
	AltWord    []string
	Defi       string
	AltDefi    []string
	DefiFormat DefiFormat
	IsInfo     bool
	Error      error
}
