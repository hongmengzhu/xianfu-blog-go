package templatePg

// Delims 分隔符
// Delims represents a set of Left and Right delimiters for HTML template rendering.
type Delims struct {
	// Left delimiter, defaults to {{.
	Left string
	// Right delimiter, defaults to }}.
	Right string
}

// NewDelims 默认分隔符
func NewDelims() Delims {
	return Delims{Left: "{{", Right: "}}"}
}
