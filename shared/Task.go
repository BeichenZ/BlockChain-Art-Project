package shared

type TaskPayload struct {
	SenderID         int
	ListOfDirections []Coordinate
	SendlogMessage   []byte
}
