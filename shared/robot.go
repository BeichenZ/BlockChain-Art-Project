package shared

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"net/rpc"
	"os"
	"time"
	"github.com/DistributedClocks/GoVector/govec"
)

const XMIN = "xmin"
const XMAX = "xmax"
const YMIN = "ymix"
const YMAX = "ymax"
const EXRADIUS = 6

type RobotStruct struct {
	CurrTask          TaskPayload
	RobotID           int // hardcoded
	RobotIP           string
	RobotListenConn   *rpc.Client
	RobotNeighbours	  []Neighbour
	RMap              Map
	CurPath           Path
	CurLocation       PointStruct // TODO need this? check
	ReceivedTask      []string // change this later
	JoiningSig   chan bool
	BusySig      chan bool
	WaitingSig   chan bool
	FreeSpaceSig chan bool
	WallSig      chan bool
	RightWallSig chan bool
	LeftWallSig  chan bool
	WalkSig      chan bool
	Logger       *govec.GoLog
}

type Robot interface {
	SendMyMap()
	MergeMaps(neighbourMaps []Map) error
	Explore() error //make a step base on the robat's current path
	GetMap() Map
	SendFreeSpaceSig()
}

var robotStruct RobotStruct

// FN: this robot sends map and ID to its neighbours
func (r *RobotStruct) SendMyMap() {
	return
}

func (r *RobotStruct) SendFreeSpaceSig() {
	fmt.Println("got here")
	r.FreeSpaceSig <- true
}

//error is not nil when the task queue is empty
// FN: Return list of destination points for each node in the network (one point for each node)
//     This robots destination point is placed at the beginning
//		TODO: comment : whats the error for?
func (r *RobotStruct) TaskCreation() ([]PointStruct, error) {

	xmin := r.FindMapExtrema(XMIN)
	xmax := r.FindMapExtrema(XMAX)
	ymin := r.FindMapExtrema(YMIN)
	ymax := r.FindMapExtrema(YMAX)

	center := PointStruct{Point: Coordinate{float64((xmax - xmin) / 2), float64((ymax - ymin) / 2)}}

	DestNum := len(r.RobotNeighbours) + 1

	DestPoints := FindDestPoints(DestNum, center)

	// move DestpointForMe to beginning of list
	DestPointForMe := r.FindClosestDest(DestPoints)

	tempEle := DestPoints[0]
	for idx, value := range DestPoints {
		if value == DestPointForMe {
			DestPoints[0] = value
			DestPoints[idx] = tempEle
			break
		}
	}

	return DestPoints, nil

}
// TODO: comment: yo why isnt this a switch statement?
func (r *RobotStruct) FindMapExtrema(e string) float64 {

	if e == XMAX {
		var xMax float64 = math.MinInt64
		for _, point := range r.RMap.ExploredPath {
			if xMax < point.Point.X {
				xMax = point.Point.X
			}
		}
		return xMax
	} else if e == XMIN {
		var xMin float64 = math.MaxFloat64
		for _, point := range r.RMap.ExploredPath {
			if xMin > point.Point.X {
				xMin = point.Point.X
			}
		}
		return xMin
	} else if e == YMAX {
		var yMax float64 = math.MinInt64
		for _, point := range r.RMap.ExploredPath {
			if yMax < point.Point.Y {
				yMax = point.Point.Y
			}
		}
		return yMax
	} else {
		var yMin float64 = math.MaxFloat64
		for _, point := range r.RMap.ExploredPath {
			if yMin > point.Point.Y {
				yMin = point.Point.Y
			}
		}
		return yMin
	}
}
// FN: Find destination point that will require the least amound of energy to go to
func (r *RobotStruct) FindClosestDest(lodp []PointStruct) PointStruct {
	dist := math.MaxFloat64
	var rdp PointStruct
	for _, dp := range lodp {
		del := DistBtwnTwoPoints(r.CurLocation, dp)
		if del < dist {
			dist = del
			rdp = dp
		}
	}

	return rdp
}

func (r *RobotStruct) RespondToButtons() error {
	// This function listen to GPIO
	for {
		fmt.Println(" Press j to send JoinSig \n Press b to send BusySig \n Press w to send WaitSig \n Press s to send WalkSig \n Press o to send WallSig")
		buf := bufio.NewReader(os.Stdin)
		signal, err := buf.ReadByte()
		if err != nil {
			fmt.Println(err)
		}
		command := string(signal)
		if command == "j" {
			r.JoiningSig <- true
		} else if command == "b" {
			r.BusySig <- true
		} else if command == "w" {
			r.WaitingSig <- true
		} else if command == "s" {
			r.FreeSpaceSig <- true
		} else if command == "o" {
			r.WallSig <- true
		}
	}
}
// TODOOOOO
func (r *RobotStruct) Explore() error {
	for {
		if len(r.CurPath.ListOfPCoordinates) == 0 {
			dpts, err := r.TaskCreation()
			var newPath Path
			if len(dpts) == 1 {
				//TODO
				newPath = CreatePathBetweenTwoPoints(r.CurLocation, dpts[0])
			} else {
				// send task to neighbours
			}
			if err != nil {
				fmt.Println("error generating task")
			}
			r.CurPath = newPath
			// DISPLAY task with GPIO
		}

		fmt.Println("\nWaiting for signal to proceed.....")

		select {
		case <-r.FreeSpaceSig:
			fmt.Println("FreeSpaceSig received")
			r.UpdateMap(true)
			r.SetCurrentLocation()
			r.TookOneStep() //remove the first element from r.CurPath.ListOfPCoordinates

			// Display task with GPIO
		case <-r.WallSig:
			r.UpdateMap(false)
			// Change wall path
			r.ModifyPathForWall()
			// Display task with GPIO
		case <- r.RightWallSig:
			// TODO
		case <- r.LeftWallSig:
			// TODO
		case <-r.JoiningSig:
			// TODO do joining thing
			newNeighbour := Neighbour{
				Addr: "8080",
				NID: 1,
			}
			r.RobotNeighbours = append(r.RobotNeighbours, newNeighbour)
			tasks, _ := r.TaskCreation()
			r.AllocateTaskToNeighbours(tasks)
			fmt.Println("join sig received")
		case <-r.BusySig:
			// TODO do busy thing
			// TODO merge map here?
			// TODO exchange tasks
			tasks, _ := r.TaskCreation()
			r.AllocateTaskToNeighbours(tasks)
			fmt.Println("busy sig received")
		case <-r.WaitingSig:
			// TODO do waiting thing
		}
	}
}

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

//update explored point in map:
// pointkind: 1 - freespace
// 			  2 - wall at current coordinate
// 			  3 - right bumper wall
// 			  4 - left bumper wall
func (r *RobotStruct) UpdateMap(b Button) {

	var exploredPoint PointStruct

	switch b {
		case FreeSpace: {
			exploredPoint.Point.X = r.CurLocation.Point.X + r.CurPath.ListOfPCoordinates[0].Point.X
			exploredPoint.Point.Y = r.CurLocation.Point.Y + r.CurPath.ListOfPCoordinates[0].Point.Y
			exploredPoint.PointKind = true
			exploredPoint.Traversed = true
			exploredPoint.TraversedTime = time.Now().Unix()

		break
	}
		case Wall:{
		break
}
		case RightWall:{
		break
	}
		case LeftWall:{
		break
}
	default:
		fmt.Println("UpdateMap () Found incorrect type of wall")

}

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

	// Check if the current location has been traversed already
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

func (r *RobotStruct) SetCurrentLocation() {
	r.CurLocation = r.CurPath.ListOfPCoordinates[0]
}

func (r *RobotStruct) WaitForEnoughTaskFromNeighbours() {
WaitingForEnoughTask:
	for {
		if len(r.ReceivedTask) == len(r.RobotNeighbours) {
			fmt.Println("waiting for my neighbours to send me tasks")
			// choose task
			// r.CurPath = something
			// should enter default Roaming state, aka don't need to do anything
			break WaitingForEnoughTask
		}
	}
}

func (r *RobotStruct) AllocateTaskToNeighbours(ldp []PointStruct) {
	ldpn := ldp[1:]
	rand.Seed(time.Now().UnixNano())
	for _, robotNeighbour := range r.RobotNeighbours {
		dpn := ldpn[rand.Intn(len(ldpn))]
		removeElFromlist(dpn, &ldpn)
		fmt.Println(robotNeighbour)
		// fmt.Println(neighbourRoboAddr)
		messagepayload := []byte("Sending to my number with ID:" + robotNeighbour.Addr)
		finalsend := r.Logger.PrepareSend("Sending Message", messagepayload)
		task := &TaskPayload{
			SenderID:         r.RobotID,
			DestPoint: 		  dpn,
			SendlogMessage:   finalsend,
		}
		fmt.Println("AllocateTaskToNeighbours() ")
		fmt.Println(task)
		// TESTING UNCOMMENT
		neighbourClient, err := rpc.Dial("tcp", robotNeighbour.Addr)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%+v", neighbourClient)
		alive := false
		// Here I send my robot the task
		err = neighbourClient.Call("RobotRPC.ReceiveTask", task, &alive)
		if err != nil {
			fmt.Println(err)
		}
		// TESTING UNCOMMENT
	}
}

func InitRobot(rID int, initMap Map, logger *govec.GoLog) *RobotStruct {
	newRobot := RobotStruct{
		RobotID:           rID,
		RobotNeighbours:   []Neighbour{},
		RMap:              initMap,
		JoiningSig:        make(chan bool),
		BusySig:           make(chan bool),
		WaitingSig:        make(chan bool),
		FreeSpaceSig:      make(chan bool),
		WallSig:           make(chan bool),
		RightWallSig:      make(chan bool),
		LeftWallSig:	   make(chan bool),
		WalkSig:           make(chan bool),
		Logger:            logger,
	}
	// newRobot.CurPath.ListOfPCoordinates = append(newRobot.CurPath.ListOfPCoordinates, shared.PointStruct{PointKind: true})
	return &newRobot
}
