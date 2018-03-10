package shared

import "net/rpc"

type RobotStruct struct {
	RobotID         uint // hardcoded
	RobotIP         string
	RobotListenConn *rpc.Client
	RMap            Map
	CurPath         Path
}

type Robot interface {
	SendMyMap(rID uint, rMap Map)
	MergeMaps(neighbourMaps []Map)
}

func (r *RobotStruct) MergeMaps(neighbourMaps []Map) error {
	newMap := Map{}
	for _, robotMap := range neighbourMaps {
		for _, coordinate := range robotMap.ExploredPath {
			for _, newCor := range newMap.ExploredPath {
				if (newCor.Point.X == coordinate.Point.X) && (newCor.Point.Y == coordinate.Point.Y) {
					if coordinate.TraversedTime > newCor.TraversedTime {
						newCor.Point.X = coordinate.Point.X
						newCor.Point.Y = coordinate.Point.Y
					}
				} else {
					newMap.ExploredPath = append(newMap.ExploredPath, coordinate)
				}
			}
		}
	}
	return nil
}
