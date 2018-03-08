package shared

import "net/rpc"

type RobotStruct struct {
	RobotID         uint // hardcoded
	RobotIP         string
	RobotListenConn *rpc.Client
	RMap            Map
	CurPath         Path
}

type Robot interface{
	SendMyMap(rID uint, rMap Map)
}
