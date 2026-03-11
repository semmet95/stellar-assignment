package asset

// TODO: generate all the models using OAPI
type Asset struct {
	Name string
	// TODO: custom types to represent a device type
	Type string
	ID   string
	// TODO: custom types to represent supported protocols
	ConnProtocol string
	// Q: Not sure about the map value type
	RegisterMap map[string]uint16
}
