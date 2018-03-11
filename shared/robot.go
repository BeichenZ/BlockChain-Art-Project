package shared

import (
	"net/rpc"
)

type RobotStruct struct {
	RobotID         uint // hardcoded
	RobotIP         string
	RobotListenConn *rpc.Client
	RMap            Map
	CurPath         Path
	CurLocation     PointStruct
	JoiningSig      chan bool
	BusySig         chan bool
	WaitingSig      chan bool
	FreeSpaceSig    chan bool
	WallSig         chan bool
}

type Robot interface {
	SendMyMap(rID uint, rMap Map)
	MergeMaps(neighbourMaps []Map) error
	Explore() error//make a step base on the robat's current path
	GetMap() Map
}

var robotStruct RobotStruct

func (r *RobotStruct) SendMyMap(rId uint, rMap Map) {
	return
}

func (r *RobotStruct) Explore() error {
	for {
		select {
		case <-r.JoiningSig:
			// TODO do joining thing
		case <-r.BusySig:
			// TODO do busy thing
			// TODO merge map here?
		case <-r.WaitingSig:
			// TODO do waiting thing
		default:
			select {
			case <-r.FreeSpaceSig:
				// TODO what happens after a button is pressed
			case <-r.WallSig:
				// TODO update when hitting a wall
			default:
				// TODO walk around randomly 1 unit
			}

		}
	}
}

// Assuming same coordinate system, and each robot has difference ExploredPath
func (r *RobotStruct) MergeMaps(neighbourMaps []Map) error {
	newMap := r.RMap

	for _, robotMap := range neighbourMaps {
		for _, coordinate := range robotMap.ExploredPath {
			if len(newMap.ExploredPath) == 0 {
				r.RMap.ExploredPath = append(r.RMap.ExploredPath, coordinate)
			} else {

				for _, newCor := range newMap.ExploredPath {
					if (newCor.Point.X == coordinate.Point.X) && (newCor.Point.Y == coordinate.Point.Y) {
						if coordinate.TraversedTime > newCor.TraversedTime {
							newCor.Point.X = coordinate.Point.X
							newCor.Point.Y = coordinate.Point.Y
						}
					} else {
						r.RMap.ExploredPath = append(r.RMap.ExploredPath, coordinate)
						r.RMap.FrameOfRef = r.RobotID
					}
				}
			}
		}
	}
	return nil
}

func (r *RobotStruct) GetMap() Map {
	return r.RMap
}

func InitRobot(rID uint, initMap Map) Robot {
	robotStruct.RobotID = rID
	robotStruct.RMap = initMap
	return &robotStruct
}

func (r *RobotStruct) SetCurrentLocation(location PointStruct) {
	r.CurLocation = location
}
