package serializer

type Format string

const (
	undefined Format = "undefined"
	JSON             = "JSON"
	XML              = "XML"
	CSV              = "CSV"
)
