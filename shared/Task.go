package shared

type TaskPayload struct {
	SenderID       int
	DestPoint      PointStruct
	SendlogMessage []byte
	Energy         int
}
type TaskDescisionPayload struct {
	SenderID  int
	Descision bool // true -> robot will do other robots task
}
