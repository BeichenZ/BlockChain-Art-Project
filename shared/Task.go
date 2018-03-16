package shared

type TaskPayload struct {
	SenderID         int
	DestPoint PointStruct
	SendlogMessage   []byte
}
