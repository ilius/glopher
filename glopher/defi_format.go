package glopher

// format of the definition
type DefiFormat int8

var (
	DefiFormatPlain = DefiFormat('m') // plain text
	DefiFormatHTML  = DefiFormat('h') // html
	DefiFormatXDXF  = DefiFormat('x') // xdxf
)
