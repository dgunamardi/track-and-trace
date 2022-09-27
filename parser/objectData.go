package parser

type ObjectData interface {
	PopulateWithMap(record map[string]string)
	IsValid() bool
	GetId() string
}
