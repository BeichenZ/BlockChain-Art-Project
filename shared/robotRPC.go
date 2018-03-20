package shared

import "fmt"

type RobotRPC struct {
	PiRobot *RobotStruct
}

func (robotRPC *RobotRPC) ReceiveMap(senderMap *Map, reply *int) error {
	return nil
}

func (robotRPC *RobotRPC) ReceiveTask(senderTask *TaskPayload, reply *int) error {
	fmt.Println(senderTask.SendlogMessage)
	robotRPC.PiRobot.Logger.UnpackReceive("Receiving Message", senderTask.SendlogMessage, TaskPayload{})
	return nil
}

// TODO
func (robotRPC *RobotRPC) ReceiveTaskDecsionResponse(senderTaskDecision *TaskDescisionPayload, reply *int) error {
	return nil
}
