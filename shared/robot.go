package shared

import (
	"math/rand"
	"net/rpc"
	"time"
)

type RobotStruct struct {
	RobotID         uint // hardcoded
	RobotIP         string
	RobotListenConn *rpc.Client
	RMap            Map
	CurPath         Path
	CurLocation     PointStruct
	NextStep        Coordinate
	JoiningSig      chan bool
	BusySig         chan bool
	WaitingSig      chan bool
	FreeSpaceSig    chan bool
	WallSig         chan bool
	WalkSig         chan bool
}

type Robot interface {
	SendMyMap(rID uint, rMap Map)
	MergeMaps(neighbourMaps []Map) error
	Explore() error //make a step base on the robat's current path
	GetMap() Map
}

var robotStruct RobotStruct

func (r *RobotStruct) SendMyMap(rId uint, rMap Map) {
	return
}

func (r *RobotStruct) Explore() error {
	for {
		time.Sleep(time.Millisecond * time.Duration(1000))
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

				//TODO update its own map, check if the step was taken before
				newLocation := PointStruct{
					Point: Coordinate{
						X: r.CurLocation.Point.X + r.NextStep.X,
						Y: r.CurLocation.Point.Y + r.NextStep.Y,
					},
					PointKind:     true,
					TraversedTime: time.Now().Unix(),
					Traversed:     true,
				}

				exist, index := CheckExist(newLocation, r.RMap.ExploredPath)

				if exist {
					oldcoor := &(r.RMap.ExploredPath[index])
					oldcoor.Point.X = newLocation.Point.X
					oldcoor.Point.Y = newLocation.Point.Y
					oldcoor.TraversedTime = newLocation.TraversedTime
					oldcoor.Traversed = newLocation.Traversed
					oldcoor.PointKind = newLocation.PointKind
				} else {
					r.RMap.ExploredPath = append(r.RMap.ExploredPath, newLocation)
				}
				r.WalkSig <- true
			case <-r.WallSig:
				// TODO update when hitting a wall
			case <-r.WalkSig:
				// TODO walk around randomly by 1 unit
				randomNum := rand.Intn(3)
				if randomNum == 0 {
					// go north
					r.NextStep = Coordinate{0.0, 1.0}
				} else if randomNum == 1 {
					// go east
					r.NextStep = Coordinate{1.0, 0.0}
				} else if randomNum == 2 {
					// go south
					r.NextStep = Coordinate{0.0, -1.0}
				} else {
					// go west
					r.NextStep = Coordinate{-1.0, 0.0}
				}
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

func CheckExist(coordinate PointStruct, cooArr []PointStruct) (bool, int) {
	for i, point := range cooArr {
		if point.Point.X == coordinate.Point.X && point.Point.Y == coordinate.Point.Y {
			return true, i
		}
	}
	return false, -1
}
