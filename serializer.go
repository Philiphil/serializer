package serializer

// Serializer is the main serializer struct
type Serializer struct {
	Format Format
}

// NewSerializer creates a new instance of Serializer
func NewSerializer(format Format) *Serializer {
	return &Serializer{Format: format}
}
