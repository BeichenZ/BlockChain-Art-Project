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
	// CurPath        []Coordinate // TODO: yo micheal here uncomment, n delete the whole struct
	CurLocation       Coordinate // TODO why isn't type coordinate instead?
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

	center := Coordinate{float64((xmax - xmin) / 2), float64((ymax - ymin) / 2)}

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
		del := DistBtwnTwoPoints(dp.Point, r.CurLocation)
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
			if err != nil {
				fmt.Println("error generating task")
			}
			var newPath Path
			if len(dpts) == 1 {
				newPath = CreatePathBetweenTwoPoints(r.CurLocation, dpts[0].Point)
			} else {
				return CodeError("Explore() > 1 destination point returned when it should have no neighbours")
			}
			r.CurPath = newPath
			// DISPLAY task with GPIO
		}

		fmt.Println("\nWaiting for signal to proceed.....")

		select {
		case <-r.FreeSpaceSig:
			fmt.Println("FreeSpaceSig received")
			r.UpdateMap(FreeSpace)
			r.UpdateCurLocation()
			r.UpdatePath()
			// r.SetCurrentLocation()
			// r.TookOneStep() //remove the first element from r.CurPath.ListOfPCoordinates

			// Display task with GPIO
		case <-r.WallSig:
			r.UpdateMap(Wall)
			r.ModifyPathForWall()
			// Display task with GPIO
		case <- r.RightWallSig:
			r.UpdateMap(RightWall)
		case <- r.LeftWallSig:
			r.UpdateMap(LeftWall)
		case <-r.JoiningSig:
			fmt.Println("join sig received")
			// TODO follow procedure to ensure all neighbours are valid to be in the network
			newNeighbour := Neighbour{
				Addr: "8080",
				NID: 1,
			}
			r.RobotNeighbours = append(r.RobotNeighbours, newNeighbour)
		case <-r.BusySig: // TODO whole thing
			fmt.Println("busy sig received")
			// Exchange my map with neighbours
			// Wait till maps from all neighbours are recevied
			// Merge my map with neighbours
			// Create tasks for current robot network
			tasks, _ := r.TaskCreation()
			// Allocate tasks to current robot network
			r.AllocateTaskToNeighbours(tasks)
			// Wait for tasks from each neighbour
			// Respond to each task given by my fellow robots
			// Agree with everyone in the network of who assigned the task
			//		- YES --> set newTaskthreshold thing, create new path based on new task
			//		- NO --> handle case ?
			// set busysig off
			// procede with new task
		case <-r.WaitingSig: // TODO
			// keep pinging the neighbour that is within it's communication radius
				// if neighbour in busy state
					// YES -> keep pinging
					// NO -> - turn WaitingSig off
					//		 - turn JoingingSig on
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
// FN: Removes the just traversed coordinate (first element in the Path list)
func (r *RobotStruct) UpdatePath() {
	r.CurPath.ListOfPCoordinates = r.CurPath.ListOfPCoordinates[1:]
}

//update explored point in map:
// pointkind: 1 - freespace
// 			  2 - wall at current coordinate
// 			  3 - right bumper wall
// 			  4 - left bumper wall
func (r *RobotStruct) UpdateMap(b Button) error {

	var justExploredPoint PointStruct

	switch b {
		case FreeSpace: {
			justExploredPoint.Point.X = r.CurLocation.X + r.CurPath.ListOfPCoordinates[0].Point.X
			justExploredPoint.Point.Y = r.CurLocation.Y + r.CurPath.ListOfPCoordinates[0].Point.Y
			justExploredPoint.PointKind = true
			justExploredPoint.Traversed = true
			justExploredPoint.TraversedTime = time.Now().Unix()

			break
	}
		case Wall:{
			justExploredPoint.Point.X = r.CurLocation.X + r.CurPath.ListOfPCoordinates[0].Point.X
			justExploredPoint.Point.Y = r.CurLocation.Y + r.CurPath.ListOfPCoordinates[0].Point.Y
			justExploredPoint.PointKind = false
			justExploredPoint.Traversed = true
			justExploredPoint.TraversedTime = time.Now().Unix()

			break
}
		case RightWall:{
			justExploredPoint.Point.X = r.CurLocation.X + r.CurPath.ListOfPCoordinates[0].Point.X + 1
			justExploredPoint.Point.Y = r.CurLocation.Y + r.CurPath.ListOfPCoordinates[0].Point.Y
			justExploredPoint.PointKind = false
			justExploredPoint.Traversed = true
			justExploredPoint.TraversedTime = time.Now().Unix()

			break
	}
		case LeftWall:{
			justExploredPoint.Point.X = r.CurLocation.X + r.CurPath.ListOfPCoordinates[0].Point.X - 1
			justExploredPoint.Point.Y = r.CurLocation.Y + r.CurPath.ListOfPCoordinates[0].Point.Y
			justExploredPoint.PointKind = true
			justExploredPoint.Traversed = true
			justExploredPoint.TraversedTime = time.Now().Unix()

			break
}
	default:
		fmt.Println("UpdateMap () Found incorrect type of wall -- CODE INCORRECT")
		return CodeError("UpdateMap () Found incorrect type of wall")

}

	oldcoor, exists := r.RMap.ExploredPath[justExploredPoint.Point]
	if exists {
		oldcoor.TraversedTime = justExploredPoint.TraversedTime
		oldcoor.Traversed = justExploredPoint.Traversed
		oldcoor.PointKind = justExploredPoint.PointKind
	}

	return nil
}

// Assuming same coordinate system, and each robot has difference ExploredPath
func (r *RobotStruct) MergeMaps(neighbourMaps []Map) error {
	refToOriginalMap := r.RMap

	for _, neighbourRobotMap := range neighbourMaps {

		if len(refToOriginalMap.ExploredPath) == 0 {
			r.RMap.ExploredPath = neighbourRobotMap.ExploredPath
		} else {
			neighbourExploredPath := neighbourRobotMap.ExploredPath

			for neighbourCoordinate, neighbourPointInfo := range neighbourExploredPath {
				if currentPointInfo, ok := r.RMap.ExploredPath[neighbourCoordinate]; ok &&
					currentPointInfo.TraversedTime < neighbourPointInfo.TraversedTime  {

					r.RMap.ExploredPath[neighbourCoordinate] = neighbourPointInfo
					continue
				}
				r.RMap.ExploredPath[neighbourCoordinate] = neighbourPointInfo
			}

		}



		//for neighbourCoordinate, neighbourPointStruct := range neighbourRobotMap.ExploredPath {
		//	if len(refToOriginalMap.ExploredPath) == 0 {
		//		r.RMap.ExploredPath[neighbourCoordinate] = neighbourPointStruct
		//	} else {
		//		for origCor, origPointStruct  := range refToOriginalMap.ExploredPath {
		//			if (origCor.X == neighbourCoordinate.X) && (origCor.Y == neighbourCoordinate.Y) {
		//
		//				var updatePointStruct PointStruct
		//				var updatedCoordinate Coordinate
		//
		//				if neighbourPointStruct.TraversedTime > origPointStruct.TraversedTime {
		//					updatePointStruct = neighbourPointStruct
		//					updatedCoordinate = neighbourCoordinate
		//				} else {
		//					updatePointStruct = origPointStruct
		//					updatedCoordinate = origCor
		//				}
		//
		//				r.RMap.ExploredPath[updatedCoordinate] = updatePointStruct
		//			} else {
		//				r.RMap.ExploredPath[newCor] = newPointStruct
		//				r.RMap.FrameOfRef = r.RobotID
		//			}
		//		}
		//	}
		//}
	}
	return nil
}

func (r *RobotStruct) GetMap() Map {
	return r.RMap
}

func (r *RobotStruct) SetCurrentLocation() {
	r.CurLocation = r.CurPath.ListOfPCoordinates[0].Point
}
func (r *RobotStruct) UpdateCurLocation() {

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
