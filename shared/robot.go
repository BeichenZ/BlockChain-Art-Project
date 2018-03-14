package shared

import (
	"fmt"
	"net/rpc"
	"time"
)

const XMIN = "xmin"
const XMAX = "xmax"
const YMIN = "ymix"
const YMAX = "ymax"

type RobotStruct struct {
	RobotID           uint // hardcoded
	RobotIP           string
	RobotListenConn   *rpc.Client
	RobotNeighbourNum int
	RMap              Map
	CurPath           Path
	CurLocation       PointStruct
	NeighboursAddr    []string
	ReceivedTask      []string // change this later
	//CurrentStep        	Coordinate
	JoiningSig   chan bool
	BusySig      chan bool
	WaitingSig   chan bool
	FreeSpaceSig chan bool
	WallSig      chan bool
	WalkSig      chan bool
}

type Robot interface {
	SendMyMap(rID uint, rMap Map)
	MergeMaps(neighbourMaps []Map) error
	Explore() error //make a step base on the robat's current path
	GetMap() Map
	SendFreeSpaceSig()
}

var robotStruct RobotStruct

func (r *RobotStruct) SendMyMap(rId uint, rMap Map) {
	return
}

func (r *RobotStruct) SendFreeSpaceSig() {
	fmt.Println("got here")
	r.FreeSpaceSig <- true
}

//error is not nil when the task queue is empty
func (r *RobotStruct) TaskCreation() (Path, error) {
	newTask := Path{}
	xmin := r.findMapExtrema(XMIN)
	xmax := r.findMapExtrema(XMAX)
	ymin := r.findMapExtrema(YMIN)
	ymax := r.findMapExtrema(YMAX)

	if r.RobotNeighbourNum == 0 {

	}

	//hard coded the newTask for testing purpose
	for i := 0; i < 5; i++ {
		pointToAdd := PointStruct{Coordinate{float64(i), float64(i + 1)}, true, 0, false}
		newTask.ListOfPCoordinates = append(newTask.ListOfPCoordinates, pointToAdd)
	}

	return newTask, nil
}

func (r *RobotStruct) findMapExtrema(e string) float64 {

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
			// TODO exchange tasks
		case <-r.WaitingSig:
			// TODO do waiting thing
		default:

			if len(r.CurPath.ListOfPCoordinates) == 0 {
				newTask, err := r.TaskCreation()
				if err != nil {
					fmt.Println("error generating task")
				}
				r.CurPath = newTask
				// DISPLAY task with GPIO
			}

			fmt.Println("\nWaiting for signal to proceed.....")

			select {
			case <-r.FreeSpaceSig:

				r.UpdateMap(true)
				r.SetCurrentLocation()
				r.TookOneStep() //remove the first element from r.CurPath.ListOfPCoordinates
				//r.UpdateCurrentStep()

				// Display task with GPIO
			case <-r.WallSig:
				r.UpdateMap(false)
				// Change wall path
				r.ModifyPathForWall()
				//r.UpdateCurrentStep()
				// Display task with GPIO
			}

		}
	}
}

//func (r *RobotStruct) UpdateCurrentStep() {
//	r.CurrentStep = Coordinate{X:r.CurPath.ListOfPCoordinates[0].Point.X, Y: r.CurPath.ListOfPCoordinates[0].Point.Y}
//}
func (r *RobotStruct) ModifyPathForWall() {

	wallCoor := r.CurPath.ListOfPCoordinates[0]
	tempList := r.CurPath.ListOfPCoordinates
	//tempList := make([]Coordinate, 0)
	for i, c := range tempList {
		if wallCoor == c {
			continue
		}
		r.CurPath.ListOfPCoordinates = r.CurPath.ListOfPCoordinates[i:]
		break
	}
}

func (r *RobotStruct) TookOneStep() {
	r.CurPath.ListOfPCoordinates = r.CurPath.ListOfPCoordinates[1:]
}

//pointkind: true => freespace, false  => wall
func (r *RobotStruct) UpdateMap(pointKind bool) {

	newLocation := PointStruct{
		Point: Coordinate{
			X: r.CurLocation.Point.X + r.CurPath.ListOfPCoordinates[0].Point.X,
			Y: r.CurLocation.Point.Y + r.CurPath.ListOfPCoordinates[0].Point.Y,
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
	robotStruct.FreeSpaceSig = make(chan bool)

	//JoiningSig      chan bool
	//BusySig         chan bool
	//WaitingSig      chan bool
	//FreeSpaceSig    chan bool
	//WallSig         chan bool
	//WalkSig         chan bool
	return &robotStruct
}

func (r *RobotStruct) SetCurrentLocation() {
	r.CurLocation = r.CurPath.ListOfPCoordinates[0]
}

func CheckExist(coordinate PointStruct, cooArr []PointStruct) (bool, int) {
	for i, point := range cooArr {
		if point.Point.X == coordinate.Point.X && point.Point.Y == coordinate.Point.Y {
			return true, i
		}
	}
	return false, -1
}

func (r *RobotStruct) WaitForEnoughTaskFromNeighbours() {
WaitingForEnoughTask:
	for {
		if len(r.ReceivedTask) == len(r.NeighboursAddr) {
			fmt.Println("waiting for my neighbours to send me tasks")
			// choose task
			// r.CurPath = something
			// should enter default Roaming state, aka don't need to do anything
			break WaitingForEnoughTask
		}
	}
}

func (r *RobotStruct) AllocateTaskToNeighbours() {
	for _, neighbourRoboAddr := range r.NeighboursAddr {
		neighbourClient, error := rpc.Dial("tcp", neighbourRoboAddr)
		alive := false
		// Here I send my robot the task
		err := neighbourClient.Call("RobotRPC.ReceiveTask", alive, &alive)
	}
}
