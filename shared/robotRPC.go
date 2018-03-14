package shared

type RobotRPC struct {
	PiRobot *Robot
}

func (robotRPC *RobotRPC) ReceiveMap(senderMap *Map, reply *int) error {
	return nil
}

func (robotRPC *RobotRPC) ReceiveTask(senderTask *Task, reply *int) error {
	return nil
}
