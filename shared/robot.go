package shared

import (
	"fmt"
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

func (r *RobotStruct) TakeRandomStep() {
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

//error is not nil when the task queue is empty
func (r *RobotStruct) TakeNextTask() (Path, error) {
	//r.NextStep = r.CurPath.ListOfPCoordinates //need to convert ListOfPCoordinates to the directions
	newTask := Path{}
	return newTask, nil
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

			if len(r.CurPath.ListOfPCoordinates) == 0 {
				newTask, err := r.TakeNextTask()
				if err != nil {
					fmt.Println("error generating task")
				}
				r.NextStep = Coordinate{newTask.ListOfPCoordinates[0].Point.X, newTask.ListOfPCoordinates[0].Point.Y}
				r.CurPath = newTask
				// DISPLAY task with GPIO
			}

			select {
			case <-r.FreeSpaceSig:
				r.UpdateLocation(true)
				r.TakeOneStep()
				// Display task with GPIO
			case <-r.WallSig:
				r.UpdateLocation(false)
				r.TakeOneStep()
				// Display task with GPIO
			}

		}
	}
}

func (r *RobotStruct) TakeOneStep() {
	r.CurPath.ListOfPCoordinates = r.CurPath.ListOfPCoordinates[1:]
}

//pointkind: true => freespace, false  => wall
func (r *RobotStruct) UpdateLocation(pointKind bool) {

	newLocation := PointStruct{
		Point: Coordinate{
			X: r.CurLocation.Point.X + r.NextStep.X,
			Y: r.CurLocation.Point.Y + r.NextStep.Y,
		},
		PointKind:     pointKind,
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
