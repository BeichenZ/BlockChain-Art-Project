package shared

type TaskPayload struct {
	SenderID       int
	SenderAddr     string
	DestPoint      PointStruct
	SendlogMessage []byte
}
type TaskDescisionPayload struct {
	SenderID       int
	SenderAddr     string
	Descision      bool // true -> robot will do other robots task
	SendlogMessage []byte
}
