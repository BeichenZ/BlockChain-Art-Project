package shared

import (
	"bufio"
	"strconv"
	//"bytes"
	//"encoding/gob"
	"fmt"
	//"io/ioutil"
	"math"
	"math/rand"
	"net/rpc"
	"os"
	"time"

	"github.com/DistributedClocks/GoVector/govec"
	"github.com/fatih/set"
	//	"encoding/json"
	"bytes"
	"encoding/gob"
	"net"
	"sync"
	hd44780 "../raspberryPiGo/go-hd44780"

	bgpio "../gpio"
)

const (
	FrontObstacleButton_Pin uint = 5
	FrontEmptyButton_Pin    uint = 6
	LeftObstacleButton_Pin  uint = 13
	RightObstacleButton_Pin uint = 19
)

const XMIN = "xmin"
const XMAX = "xmax"
const YMIN = "ymix"
const YMAX = "ymax"
const EXRADIUS = 6
const TIMETOJOINSECONDUNIT = 10
const TIMETOJOIN = TIMETOJOINSECONDUNIT * time.Second
const INITENERGY = 3000

var DEFAULTPATH = []PointStruct{SOUTH, SOUTH, SOUTH, WEST, WEST, WEST, NORTH, NORTH, NORTH, EAST, EAST, EAST, EAST}

type JoiningInfo struct {
	joiningTime      time.Time
	firstTimeJoining bool
}

type RobotLog struct {
	CurrTask    TaskPayload
	CurPath     Path
	RMap        Map
	CurLocation Coordinate
	REnergy		int
}

type RobotStruct struct {
	CurrTask           TaskPayload
	PossibleNeighbours *set.Set
	RobotID            int // hardcoded
	RobotIP            string
	RobotEnergy        int
	RobotListenConn    *rpc.Client
	RobotNeighbours RobotNeighboursMutex
	RMap            Map
	CurPath         Path
	CurLocation           Coordinate
	ReceivedTasks         []TaskPayload // change this later
	ReceivedTasksResponse []TaskDescisionPayload
	JoiningSig            chan Neighbour
	BusySig               chan bool
	WaitingSig            chan bool
	FreeSpaceSig          chan bool
	WallSig               chan bool
	RightWallSig          chan bool
	LeftWallSig           chan bool
	WalkSig               chan bool
	Logname               string
	Logger                *govec.GoLog
	State                 RobotMutexState
	joinInfo              JoiningInfo
	exchangeFlag          CanExchangeInfoWithRobots
}

type Robot interface {
	SendMyMap()
	MergeMaps(neighbourMaps []Map) error
	Explore() error //make a step base on the robat's current path
	GetMap() Map
	SendFreeSpaceSig()
}

type CanExchangeInfoWithRobots struct {
	sync.RWMutex
	flag bool
}

type RobotMutexState struct {
	sync.RWMutex
	rState RobotState
}

type RobotNeighboursMutex struct {
	sync.RWMutex
	rNeighbour map[int]Neighbour
}

var robotStruct RobotStruct

// FN: this robot sends map and ID to its neighbours
func (r *RobotStruct) SendMyMap() {
	return
}

type RobotState int

const (
	ROAM RobotState = iota
	JOIN RobotState = iota
	BUSY RobotState = iota
)

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

	center := Coordinate{Round(float64((xmax - xmin) / 2)), Round(float64((ymax - ymin) / 2))}
	center.X = Round(center.X)
	center.Y = Round(center.Y)
	//r.RobotNeighbours.Lock()
	r.RemoveDeadNeighbours()
	DestNum := len(r.RobotNeighbours.rNeighbour) + 1
	//r.RobotNeighbours.Unlock()
	//fmt.Println("DESTNum is ")
	//fmt.Println(DestNum)
	//fmt.Println(r.RobotNeighbours)

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

func (r *RobotStruct) FindMapExtrema(e string) float64 {

	if e == XMAX {
		var xMax float64 = math.MinInt64
		for _, point := range r.RMap.ExploredPath {
			if xMax < point.Point.X {
				xMax = point.Point.X
			}
		}

		if len(r.RMap.ExploredPath) == 0 {
			return 0.0
		}

		return Round(xMax)
	} else if e == XMIN {
		var xMin float64 = math.MaxFloat64
		for _, point := range r.RMap.ExploredPath {
			if xMin > point.Point.X {
				xMin = point.Point.X
			}
		}

		if len(r.RMap.ExploredPath) == 0 {
			return 0.0
		}
		return Round(xMin)
	} else if e == YMAX {
		var yMax float64 = math.MinInt64
		for _, point := range r.RMap.ExploredPath {
			if yMax < point.Point.Y {
				yMax = point.Point.Y
			}
		}

		if len(r.RMap.ExploredPath) == 0 {
			return 0.0
		}
		return Round(yMax)
	} else {
		var yMin float64 = math.MaxFloat64
		for _, point := range r.RMap.ExploredPath {
			if yMin > point.Point.Y {
				yMin = point.Point.Y
			}
		}

		if len(r.RMap.ExploredPath) == 0 {
			return 0.0
		}
		return Round(yMin)
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

// FN: Depending on which button is pressed, set the according signal only if it is in the ROAM state
func (r *RobotStruct) RespondToButtons() error {
	// This function listen to GPIO
	for {
		fmt.Println(" Press j to send JoinSig \n Press b to send BusySig \n Press w to send WaitSig \n Press f to send WalkSig \n Press u to send WallSig \n Press c to send RightWallSig \n Press k to send LeftWallSig")
		buf := bufio.NewReader(os.Stdin)
		signal, err := buf.ReadByte()
		if err != nil {
			fmt.Println(err)
		}
		command := string(signal)
		//TODO:Check for current state before sending signals. No signals should be sent during busy state
		if command == "j" {

			r.JoiningSig <- Neighbour{
				Addr:                ":8080",
				NID:                 1,
				NMap:                RandomMapGenerator(),
				NeighbourCoordinate: Coordinate{4.0, 5.0},
			}

		} else if command == "b" {
			r.BusySig <- true
		} else if command == "f" {
			r.FreeSpaceSig <- true
		} else if command == "w" {
			r.WallSig <- true
		} else if command == "r" {
			r.RightWallSig <- true
		} else if command == "l" {
			r.LeftWallSig <- true
		}
	}
}

func (r *RobotStruct) Explore() error {
	fmt.Printf("1 Explore() start of explore. Robot ID %+v\n", r.RobotID)

	lcd := hd44780.NewGPIO4bit()
	if err := lcd.Open();err != nil {
		panic("Cannot OPen lcd:"+err.Error())
	}
	defer lcd.Close()

	for {

		if len(r.CurPath.ListOfPCoordinates) == 0 {
			dpts, err := r.TaskCreation()
			fmt.Println("0. The new destination point fresh from TaskCreation() ", dpts[0].Point)
			dpts[0].Point.X = dpts[0].Point.X + r.CurLocation.X
			dpts[0].Point.Y = dpts[0].Point.Y + r.CurLocation.Y
			fmt.Println("The new destination point following offset addition", dpts[0].Point)
			if err != nil {
				fmt.Println("error generating task")
			}
			var newPath Path
			// TESTING
			newPath = CreatePathBetweenTwoPoints(r.CurLocation, dpts[0].Point)
			// TESTING

			/**** SHOULDN'T BE COMMENTED OUT JUST FOR TESTING
			if len(dpts) == 1 {
				newPath = CreatePathBetweenTwoPoints(r.CurLocation, dpts[0].Point)
			} else {
				fmt.Println("Explore() > 1 destination point returned when it should have no neighbours")
				return CodeError("Explore() > 1 destination point returned when it should have no neighbours")
			}
			*/
			r.CurPath = newPath

			fmt.Println("1. Explore() current path ", r.CurPath)

			r.WriteToLog()

		}
		var dir string


		fmt.Println("CHECKING FOR THE FIRST DIRECTION")
		switch r.CurPath.ListOfPCoordinates[0].Point {
		case WEST.Point:
			dir = "WEST"
			fmt.Println("LCD should display WEST")
			break;
		case EAST.Point:
			dir = "EAST"
			fmt.Println("LCD should display EAST")
			break;
		case SOUTH.Point:
			dir = "SOUTH"
			fmt.Println("LCD should display South")
			break;
		case NORTH.Point:
			dir = "NORTH"
			fmt.Println("LCD should display NORTH")
			break;
		default:
			fmt.Println("Current path direction incorrect")
			break
		}
		lcd.DisplayLines(dir)
		fmt.Println(" 2 Explore() \nWaiting for signal to proceed.....")


		if r.RobotEnergy==0 {
			lcd.Display(1,  "Energy Level is 0")
			os.Exit(0)
		}

		select {
		case <-r.FreeSpaceSig:
			fmt.Println("FreeSpaceSig received")
			r.UpdateMap(FreeSpace)
			fmt.Println("My current task: ", r.CurrTask)
			fmt.Println("Cur location before: ", r.CurLocation)
			fmt.Println("Cur path before: ", r.CurPath)
			r.UpdateCurLocation()
			fmt.Println("Cur location after: ", r.CurLocation)
			r.UpdatePath()
			fmt.Println("Cur path after: ", r.CurPath)
			// update LCD

			// Display task with GPIO
		case <-r.WallSig:
			fmt.Println("front wall sig received")
			r.UpdateMap(Wall)
			r.ModifyPathForWall()
			// Display task with GPIO
		case <-r.RightWallSig:
			fmt.Println("right wall sig received")
			r.UpdateMap(RightWall)
		case <-r.LeftWallSig:
			fmt.Println("left wall sig received")
			r.UpdateMap(LeftWall)
		case <-r.BusySig: // TODO whole thing
			fmt.Println("3 Explore() busy sig received. Robot ID %+v Robot state: %+v", r.RobotID, r.State)

			//listOfNeighbourMaps :=  make([]Map, len(r.RobotNeighbours))

			fmt.Println("THE CURRENT MAP IS BEFORE MERGING")
			fmt.Println(r.RMap)
			// fmt.Println("Getting the maps from the neighbour.................")
			r.RobotNeighbours.Lock()
			for k, nei := range r.RobotNeighbours.rNeighbour {
				neighbourMap := Map{}
				client, err := rpc.Dial("tcp", nei.Addr)
				if err != nil {
					fmt.Println("4 Explore() ", err)
					delete(r.RobotNeighbours.rNeighbour, k)
					continue
				}
				// This robot recevies maps from its neighbour
				messagepayload := 1
				finalsend := r.Logger.PrepareSend("I'm request map from my neighbour: "+nei.Addr, messagepayload)

				requestMapPayload := RequestMapPayloadStruct{
					arbitaryPayload:          false,
					requestMapSendlogMessage: finalsend,
				}
				err = client.Call("RobotRPC.ReceiveMap", requestMapPayload, &neighbourMap)
				if err != nil {
					fmt.Println("5 Explore() RIP neighbour - Going to delete this")
					delete(r.RobotNeighbours.rNeighbour, k)
					client.Close()
					continue
				}
				client.Close()
				fmt.Printf("Receive map from %s \n", nei.Addr)
				fmt.Println(neighbourMap)

				nei.NMap = neighbourMap
				//listOfNeighbourMaps = append(listOfNeighbourMaps, neighbourMap)
			}
			r.RobotNeighbours.Unlock()

			fmt.Println()
			fmt.Println("Retrieved the map. Start merging..........")
			fmt.Println()

			//logging
			fmt.Println()
			fmt.Println("The CURRENT ROBOT's id is")
			fmt.Println(r.RobotID)

			fmt.Println("THE CURRENT MAP IS")
			fmt.Println(r.RMap)

			fmt.Println("The current robot state is")
			fmt.Println(r.State)
			fmt.Println()

			//r.MergeMaps(listOfNeighbourMaps)
			r.MergeMaps()
			fmt.Println()
			fmt.Println("Map after merged is ")
			fmt.Println(r.RMap)

			fmt.Println("Finished Merging")
			fmt.Println()

			//// Exchange my map with neighbours
			//// Wait till maps from all neighbours are recevied
			//// Merge my map with neighbours
			//// Create tasks for current robot network
			tasks, _ := r.TaskCreation()

			fmt.Println()
			fmt.Println("The following is the list of tasks created by ", r.RobotIP)
			for _, t := range tasks {
				fmt.Println(t)
			}
			fmt.Println()

			//// Allocate tasks to current robot network
			//r.CurPath = CreatePathBetweenTwoPoints(r.CurLocation, tasks[0].Point)
			//// r.CurrTask = tasks[0]
			//fmt.Println("tasks length is")
			//fmt.Println(len(tasks))
			//
			//fmt.Println("number of neighbour is ")
			//fmt.Println(len(r.RobotNeighbours)
			r.TaskAllocationToNeighbours(tasks[1:])

			//
			//// Wait for tasks from each neighbour
			fmt.Println("Done allocating tasks for neighbours")
			fmt.Println("My neighbours are ")
			fmt.Println(r.RobotNeighbours)
			r.RobotNeighbours.Lock()
			fmt.Println(len(r.RobotNeighbours.rNeighbour))
			r.RobotNeighbours.Unlock()
			//rawRobotNeighbour, _:= json.MarshalIndent(r.RobotNeighbours, "", "")
			//fmt.Println(string(rawRobotNeighbour))

			r.WaitForEnoughTaskFromNeighbours()
			//// Choose task based with the lowest ID including its own
			fmt.Println("Done waiting for tasks from my neighbours")
			taskToDo := r.PickTaskWithLowestID(tasks[0])
			taskToDo.DestPoint.Point.X = taskToDo.DestPoint.Point.X + r.CurLocation.X
			taskToDo.DestPoint.Point.Y = taskToDo.DestPoint.Point.Y + r.CurLocation.Y
			r.CurPath = CreatePathBetweenTwoPoints(r.CurLocation, taskToDo.DestPoint.Point)
			r.CurrTask = taskToDo
			fmt.Println("The task I am going to dooooo ----> Sending ID", taskToDo.SenderID, "=>", taskToDo.DestPoint)
			fmt.Println()

			//// Respond to each task given by my fellow robots
			r.RespondToNeighoursAboutTask(taskToDo)

			// Wait for neighbours response
			fmt.Println("Done responding to task, going to wait for my neighbour to respond to my task")
			r.WaitForNeighbourTaskResponse()
			fmt.Println("Done getting response from all neighbour")

			fmt.Println("CALLING UPDATE UpdateStateForNewJourney")
			r.CurrTask = taskToDo
			r.CurPath = CreatePathBetweenTwoPoints(r.CurLocation, taskToDo.DestPoint.Point)
			r.WriteToLog()
			r.UpdateStateForNewJourney()
			//fmt.Println("I am going to sleep now")
			//time.Sleep(10*time.Minute)

		}
	}
}

func (r *RobotStruct) ModifyPathForWall() {
	fmt.Println("ModifyPathForWall()() Pressed wall")
	wallCoor := r.CurPath.ListOfPCoordinates[0]
	tempList := r.CurPath.ListOfPCoordinates
	i := 0
	//tempList := make([]Coordinate, 0)
	for j, c := range tempList {
		i = j
		if wallCoor == c {
			continue
		}
		break
	}

	r.CurPath.ListOfPCoordinates = r.CurPath.ListOfPCoordinates[i+1:]
	if len(r.CurPath.ListOfPCoordinates) == 0 {
		r.CurPath.ListOfPCoordinates = DEFAULTPATH
		fmt.Println("changed path to default path")
		r.ModifyPathForWall()
	}

	r.WriteToLog()
}

// FN: Removes the just traversed coordinate (first element in the Path list)
func (r *RobotStruct) UpdatePath() {
	r.CurPath.ListOfPCoordinates = r.CurPath.ListOfPCoordinates[1:]
	r.WriteToLog()
}

//update explored point in map:
// pointkind: 1 - freespace
// 			  2 - wall at current coordinate
// 			  3 - right bumper wall
// 			  4 - left bumper wall
func (r *RobotStruct) UpdateMap(b Button) error {

	var justExploredPoint PointStruct

	switch b {
	case FreeSpace:
		{
			justExploredPoint.Point.X = r.CurLocation.X + r.CurPath.ListOfPCoordinates[0].Point.X
			justExploredPoint.Point.Y = r.CurLocation.Y + r.CurPath.ListOfPCoordinates[0].Point.Y
			justExploredPoint.PointKind = true
			justExploredPoint.Traversed = true
			justExploredPoint.TraversedTime = time.Now().Unix()

			break
		}
	case Wall:
		{
			justExploredPoint.Point.X = r.CurLocation.X + r.CurPath.ListOfPCoordinates[0].Point.X
			justExploredPoint.Point.Y = r.CurLocation.Y + r.CurPath.ListOfPCoordinates[0].Point.Y
			justExploredPoint.PointKind = false
			justExploredPoint.Traversed = true
			justExploredPoint.TraversedTime = time.Now().Unix()

			break
		}
	case RightWall:
		{
			justExploredPoint.Point.X = r.CurLocation.X + r.CurPath.ListOfPCoordinates[0].Point.X + 1
			justExploredPoint.Point.Y = r.CurLocation.Y + r.CurPath.ListOfPCoordinates[0].Point.Y
			justExploredPoint.PointKind = false
			justExploredPoint.Traversed = true
			justExploredPoint.TraversedTime = time.Now().Unix()

			break
		}
	case LeftWall:
		{
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
	} else {
		r.RMap.ExploredPath[justExploredPoint.Point] = justExploredPoint
	}
	// Write to log to update RMap
	r.WriteToLog()
	return nil
}

func (r *RobotStruct) RespondToNeighoursAboutTask(taskToDo TaskPayload) {
	r.RobotNeighbours.Lock()
	for ids, neighbour := range r.RobotNeighbours.rNeighbour {
		client, err := rpc.Dial("tcp", neighbour.Addr)
		if err != nil {
			fmt.Println("1 RespondToNeighoursAboutTask() ", err)
			delete(r.RobotNeighbours.rNeighbour, ids)
			continue
			//fmt.Println("There is a problem respoing to neighbour about its task")
		}
		responsePayload := ResponseForNeighbourPayload{}

		if neighbour.NID == taskToDo.SenderID {
			messagepayload := 1
			finalsend := r.Logger.PrepareSend("Sending Message - "+"Accpeting task from my neighbour:"+neighbour.Addr, messagepayload)
			taskResponsePayloadYes := TaskDescisionPayload{
				SenderID:       r.RobotID,
				SenderAddr:     r.RobotIP,
				Descision:      true,
				SendlogMessage: finalsend,
			}
			fmt.Printf("RespondToNeighoursAboutTask() Will do NeighbourID [ %+v ] task \n", neighbour.NID)
			err = client.Call("RobotRPC.ReceiveTaskDecsionResponse", taskResponsePayloadYes, &responsePayload)
			if err != nil {
				fmt.Println("2 RespondToNeighoursAboutTask() RIP neighbour - Going to delete this")
				delete(r.RobotNeighbours.rNeighbour, ids)
			}
		} else {
			messagepayload := 1
			finalsend := r.Logger.PrepareSend("Sending Message - "+"Denying task from my neighbour:"+neighbour.Addr, messagepayload)
			taskResponsePayloadNo := TaskDescisionPayload{
				SenderID:       r.RobotID,
				SenderAddr:     r.RobotIP,
				Descision:      false,
				SendlogMessage: finalsend,
			}
			fmt.Printf("RespondToNeighoursAboutTask() Will not do NeighbourID [ %+v ] task \n", neighbour.NID)

			err = client.Call("RobotRPC.ReceiveTaskDecsionResponse", taskResponsePayloadNo, &responsePayload)
			if err != nil {
				fmt.Println("3 RespondToNeighoursAboutTask() RIP neighbour - Going to delete this")
				delete(r.RobotNeighbours.rNeighbour, ids)
			}
		}
		client.Close()
	}
	r.RobotNeighbours.Unlock()

}

// New version of merge maps, uses the Neighbour struct map feild
func (r *RobotStruct) MergeMaps() {
	refToOriginalMap := r.RMap
	r.RobotNeighbours.Lock()
	for _, neighbourRobot := range r.RobotNeighbours.rNeighbour {

		if len(refToOriginalMap.ExploredPath) == 0 {
			r.RMap.ExploredPath = neighbourRobot.NMap.ExploredPath
		} else {
			neighbourExploredPath := neighbourRobot.NMap.ExploredPath

			for neighbourCoordinate, neighbourPointInfo := range neighbourExploredPath {
				if currentPointInfo, ok := r.RMap.ExploredPath[neighbourCoordinate]; ok &&
					currentPointInfo.TraversedTime < neighbourPointInfo.TraversedTime {

					r.RMap.ExploredPath[neighbourCoordinate] = neighbourPointInfo
					continue
				}
				r.RMap.ExploredPath[neighbourCoordinate] = neighbourPointInfo
			}

		}
	}
	r.RobotNeighbours.Unlock()
	r.WriteToLog()
}

func (r *RobotStruct) GetMap() Map {
	return r.RMap
}

// TODO comment: update this when path type is updated
func (r *RobotStruct) UpdateCurLocation() {
	r.CurLocation.X = r.CurLocation.X + r.CurPath.ListOfPCoordinates[0].Point.X
	r.CurLocation.Y = r.CurLocation.Y + r.CurPath.ListOfPCoordinates[0].Point.Y
	r.RobotEnergy--

	r.WriteToLog()
}

func (robot *RobotStruct) CheckAliveNeighbour() {
	for idx, val := range robot.RobotNeighbours.rNeighbour {
		_, err := rpc.Dial("tcp", val.Addr)
		if err != nil {
			delete(robot.RobotNeighbours.rNeighbour, idx)
		}
	}
}

func (r *RobotStruct) WaitForEnoughTaskFromNeighbours() {
	r.RobotNeighbours.Lock()
WaitingForEnoughTask:
	for {
		fmt.Println("received Task", len(r.ReceivedTasks))
		fmt.Println("length neighbour", len(r.RobotNeighbours.rNeighbour))
		// Check how many neighbours are alive

		r.CheckAliveNeighbour() // TODO change to our new fn

		if len(r.ReceivedTasks) == len(r.RobotNeighbours.rNeighbour) {
			fmt.Println("waiting for my neighbours to send me tasks")
			// choose task
			// r.CurPath = something
			// should enter default Roaming state, aka don't need to do anything
			r.RobotNeighbours.Unlock()
			break WaitingForEnoughTask
		}
	}
}

func (r *RobotStruct) WaitForNeighbourTaskResponse() {
	r.RobotNeighbours.Lock()
WaitingForEnoughTaskResponse:
	for {
		//fmt.Println("received Task", len(r.ReceivedTasks))
		//fmt.Println("length neighbour", len(r.RobotNeighbours))
		fmt.Println("# task receive RESPONSE ", len(r.ReceivedTasksResponse), " num neighbour ", len(r.RobotNeighbours.rNeighbour))
		r.CheckAliveNeighbour()
		if len(r.ReceivedTasksResponse) >= len(r.RobotNeighbours.rNeighbour) {
			fmt.Println("waiting for my neighbours to send me tasks")
			// choose task
			// r.CurPath = something
			// should enter default Roaming state, aka don't need to do anything
			fmt.Println("(Inside) num of received task RESPONSE ", len(r.ReceivedTasksResponse), " num neighbour ", len(r.RobotNeighbours.rNeighbour))
			r.RobotNeighbours.Unlock()
			break WaitingForEnoughTaskResponse
		}
	}

}

func (r *RobotStruct) PickTaskWithLowestID(taskFromMe PointStruct) TaskPayload {
	localMin := 100000
	var taskToDo TaskPayload
	fmt.Println("IN PICK_TASKWITHLOWESTID ")
	fmt.Printf(" Robot ID %+v and its state %+v\n", r.RobotID, r.State.rState)
	for _, task := range r.ReceivedTasks {
		if task.SenderID < localMin {
			localMin = task.SenderID
			taskToDo = task
		}
		fmt.Println("PickTaskWithLowestID() received task ", "task sender ID ", task.SenderID, " => ", task.DestPoint)
	}
	// Check if the task assigned is larger than the one it assigned itself
	if r.RobotID < taskToDo.SenderID || len(r.ReceivedTasks) == 0{
		taskToDo.SenderID = r.RobotID
		taskToDo.DestPoint = taskFromMe
	}

	return taskToDo
}

func (r *RobotStruct) TaskAllocationToNeighbours(ldp []PointStruct) {
	//fmt.Printf( "The length of LDPN is  %v \n", len(ldp))
	ldpn := ldp
	rand.Seed(time.Now().UnixNano())
	fmt.Println("In TASK ALLOCATION TO NEIGHBOURS")
	r.RobotNeighbours.Lock()
	fmt.Printf("There are %+v robots \n", len(r.RobotNeighbours.rNeighbour))
	r.RobotNeighbours.Unlock()

	r.RobotNeighbours.Lock()
	for idx, robotNeighbour := range r.RobotNeighbours.rNeighbour {
		//fmt.Printf( "The length of LDPN is  %v \n", len(ldpn))
		fmt.Println("Examining Robot ", robotNeighbour.NID)
		dpn := ldpn[rand.Intn(len(ldpn))]
		removeElFromlist(dpn, &ldpn)
		//fmt.Printf("Current Neighour %s \n", robotNeighbour)
		// fmt.Println(neighbourRoboAddr)
		messagepayload := 1
		finalsend := r.Logger.PrepareSend("Sending Task number"+strconv.Itoa(idx)+"to Robot"+robotNeighbour.Addr, messagepayload)
		task := &TaskPayload{
			SenderID:       r.RobotID,
			SenderAddr:     r.RobotIP,
			DestPoint:      dpn,
			SendlogMessage: finalsend,
		}

		fmt.Println("finalSEND")
		//fmt.Printf("TaskAllocateToNeighbours(%s -------> %s) \n", task.SenderAddr, robotNeighbour.Addr)
		//data, _ := json.MarshalIndent(task, "", "")
		//fmt.Println(string(data)

		// Dial neighbour - if neighbour not there remove from its list of neighbours
		neighbourClient, err := rpc.Dial("tcp", robotNeighbour.Addr)
		if err != nil {

			fmt.Println("1 TaskAllocationToNeighbours() ", err)
			delete(r.RobotNeighbours.rNeighbour, idx)
			continue
		}

		//fmt.Printf("%+v", neighbourClient)
		alive := false
		// Here I send my robot the task
		fmt.Println("Going to send following task ", "Sender_ID", task.SenderID, "=>", task.DestPoint)

		err = neighbourClient.Call("RobotRPC.ReceiveTask", task, &alive)

		fmt.Println("Why are you hanging????????????")
		if err != nil {
			fmt.Println("2 TaskAllocationToNeighbours() ", err)
			delete(r.RobotNeighbours.rNeighbour, idx)
		}

		neighbourClient.Close()
	}

	r.RobotNeighbours.Unlock()

	return
}

// FN: payload to ask neighbour if I and my current hommies are within this new neighbours radius
func createFarNeighbourPayload(r RobotStruct, finalsend []byte) FarNeighbourPayload {

	farNeighbourPayload := FarNeighbourPayload{
		NeighbourID:         r.RobotID,
		NeighbourIPAddr:     r.RobotIP,
		NeighbourCoordinate: r.CurLocation,
		//NeighbourMap:        r.RMap,
		SendlogMessage: finalsend,
		//State: 				 r.State.rState,
		//ItsNeighbours:       r.RobotNeighbours,
	}
	r.RobotNeighbours.Lock()
	for _, robot := range r.RobotNeighbours.rNeighbour {
		farNeighbourPayload.ItsNeighbours = append(farNeighbourPayload.ItsNeighbours, robot)
	}
	r.RobotNeighbours.Unlock()
	return farNeighbourPayload
}

func SaveNeighbour(r *RobotStruct, robotsToAdd []Neighbour) {
	for idx, val := range robotsToAdd {
		if (robotsToAdd[idx].NID == r.RobotID) || CheckNeighbourExists(r, val) {
			continue
		}
		r.RobotNeighbours.Lock()
		r.RobotNeighbours.rNeighbour[robotsToAdd[idx].NID] = val
		r.RobotNeighbours.Unlock()
	}
}
func CheckNeighbourExists(r *RobotStruct, rn Neighbour) bool {
	r.RobotNeighbours.Lock()
	for i, _ := range r.RobotNeighbours.rNeighbour {
		if i == rn.NID {
			r.RobotNeighbours.Unlock()
			return true
		}
	}
	r.RobotNeighbours.Unlock()
	return false
}

func CheckNeighbourExistsByIPd(r *RobotStruct, rn string) bool {

	r.RobotNeighbours.Lock()
	for _, val := range r.RobotNeighbours.rNeighbour {

		//fmt.Println()
		//fmt.Println("Comparing neighbour with existing neighbour")
		//fmt.Println(val.Addr)
		//fmt.Println(rn)
		//fmt.Println()

		if val.Addr == rn {

			//fmt.Printf("\nI have neighbour[%s] already\n", val.Addr)
			r.RobotNeighbours.Unlock()
			return true
		}
	}
	r.RobotNeighbours.Unlock()
	return false
}

// Client -> R2
// Fn: From the list of possible neighbours (address' that were pinged befo)
func (r *RobotStruct) CallNeighbours() {
	for {
		for _, possibleNeighbour := range r.PossibleNeighbours.List() {
			client, err := rpc.Dial("tcp", possibleNeighbour.(string))
			if err != nil {
				//fmt.Println("HUGE ERROR on possible neighbour-- all nighter")
				fmt.Println("1 CallNeighbours() ", err.Error())
				r.PossibleNeighbours.Remove(possibleNeighbour)

				r.RobotNeighbours.Lock()
				for k, nei := range r.RobotNeighbours.rNeighbour {
					if nei.Addr == possibleNeighbour.(string) {
						fmt.Println("REMOVING THE NEIGHBOUR ", nei.Addr)
						delete(r.RobotNeighbours.rNeighbour, k)
						break
					}

				}
				r.RobotNeighbours.Unlock()
				continue
			}
			r.State.Lock()
			checkROAMState := r.State.rState == ROAM
			checkJOINState := r.State.rState == JOIN
			r.State.Unlock()

			r.exchangeFlag.Lock()
			exchangeState := r.exchangeFlag.flag
			r.exchangeFlag.Unlock()

			if (checkROAMState && exchangeState) || (checkJOINState) {
				//Test
				if CheckNeighbourExistsByIPd(r, possibleNeighbour.(string)) {
					continue
				}

				responsePayload := ResponseForNeighbourPayload{}

				messagepayload := 1
				finalsend := r.Logger.PrepareSend("Sending Message - "+"Trying to call my neighbour:"+possibleNeighbour.(string), &messagepayload)
				farNeighbourPayload := createFarNeighbourPayload(*r, finalsend)
				// This robot is calling its (potential) neighbour to see if its within the communication radius of itself and its current neighbours
				fmt.Println("CallNeighbours() my ID and state ", r.RobotID, " ", r.State.rState)
				err := client.Call("RobotRPC.ReceivePossibleNeighboursPayload", farNeighbourPayload, &responsePayload)
				fmt.Println("CallNeighbours() finished calling the following neighbours ")
				if err != nil {
					fmt.Println("2 CallNeighbours() ", err)
				}

				client.Close()

				//if other robot is in join/roam and within cr, current robot tries joining
				if !responsePayload.WithInComRadius {
					continue
				}

				SaveNeighbour(r, responsePayload.NeighboursNeighbourRobots) // Client robot saves the other robot and its neighbours which ARE in CR

				r.State.Lock()
				checkJOINState := r.State.rState == JOIN
				r.State.Unlock()

				if checkJOINState {
					continue
				}
				//Up to this point, Robot with JOINNING state should exit here

				r.State.Lock()
				r.State.rState = JOIN
				r.State.Unlock()

				StartClock(responsePayload.NeighbourState, r, responsePayload.RemainingTime)
			} else {
				client.Close()
			}

		}
		//time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("CallNeighbours() Current Neighbours ", r.RobotNeighbours)

}

func StartClock(state RobotState, r *RobotStruct, remainingTime time.Duration) {

	r.joinInfo.joiningTime = time.Now()
	ticker := time.NewTicker(1000 * time.Millisecond)
	// If neighbour is in JOIN state this Robot takes timer of neighbour
	if state == JOIN {
		//fmt.Println("The Pi has the remaining time of ", responsePayload.RemainingTime)
		//robot start its time using the remainingTime

		go func() {
			counter := 0
			for _ = range ticker.C {
				counter += 1
				fmt.Printf("Joining Neighbour timer--------> Counter is %s\n", counter)
				if time.Now().Sub(r.joinInfo.joiningTime) >= (TIMETOJOIN - remainingTime) {
					fmt.Println("WE ARE FINISHED. -- JOIN")
					r.RobotNeighbours.Lock()
					fmt.Println("TImer: # of Neighbour is --joining", len(r.RobotNeighbours.rNeighbour))
					r.RobotNeighbours.Unlock()
					ticker.Stop()
					r.State.Lock()
					r.State.rState = BUSY
					r.State.Unlock()
					r.BusySig <- true
				}
			}
		}()

	} else if state == ROAM {

		//robot starts its owner timer
		go func() {
			counter := 0
			for _ = range ticker.C {
				counter += 1
				fmt.Printf("Joining My timer--------> Counter is %s\n", counter)
				if counter >= TIMETOJOINSECONDUNIT {
					fmt.Println("WE ARE FINISHED -- ROAM")
					r.RobotNeighbours.Lock()
					fmt.Println("TImer: # of Neighbour is --ROAM", len(r.RobotNeighbours.rNeighbour))
					r.RobotNeighbours.Unlock()
					ticker.Stop()
					r.State.Lock()
					r.State.rState = BUSY
					r.State.Unlock()
					r.BusySig <- true
				}
			}
		}()

	} else {
		fmt.Println("StartClock() SHOULD NOT END UP HERE -- CODE WRONG ", state)
		//do nothing
		//neightbours are in hte busy state
	}
}

// TODO
// Decide the appropriate task that the neighbours assigned it and send response to neighbours
func (r *RobotStruct) decideTaskTodo() {
	// call ReceiveTaskDecsionResponse() here

}

func InitRobot(rID int, initMap Map, ic Coordinate, logger *govec.GoLog, robotIPAddr string, logname string) *RobotStruct {
	newRobot := RobotStruct{

		PossibleNeighbours: set.New(),
		RobotID:            rID,
		RobotIP:            robotIPAddr,
		RobotEnergy:        INITENERGY,
		CurLocation:        ic,
		RobotNeighbours:    RobotNeighboursMutex{rNeighbour: make(map[int]Neighbour)},
		RMap:               initMap,
		JoiningSig:         make(chan Neighbour),
		BusySig:            make(chan bool),
		WaitingSig:         make(chan bool),
		FreeSpaceSig:       make(chan bool),
		WallSig:            make(chan bool),
		RightWallSig:       make(chan bool),
		LeftWallSig:        make(chan bool),
		WalkSig:            make(chan bool),
		Logname:            logname,
		Logger:             logger,
		State:              RobotMutexState{rState: ROAM},
		joinInfo:           JoiningInfo{time.Now(), true},
		exchangeFlag:       CanExchangeInfoWithRobots{flag: true},
	}
	// newRobot.CurPath.ListOfPCoordinates = append(newRobot.CurPath.ListOfPCoordinates, shared.PointStruct{PointKind: true})

	//tempEXploredMap := make(map[Coordinate]PointStruct)
	//tempLocation := Coordinate{float64(newRobot.RobotID) + 10.0, float64(newRobot.RobotID) + 10.0}
	//tempEXploredMap[tempLocation] = PointStruct{Point: tempLocation}
	//newRobot.RMap = Map{tempEXploredMap, 0}

	return &newRobot
}

//func (r *RobotStruct) EncodeRobotLogInfo(robotLog RobotLog) string {
//	buf := new(bytes.Buffer)
//	encoder := gob.NewEncoder(buf)
//	err := encoder.Encode(robotLog)
//	if err != nil {
//		panic(err)
//	}
//	output := string(buf.Bytes())
//	return output
//	// fmt.Println(buf.Bytes())
//}
//
//func (r *RobotStruct) ReadFromLog() {
//	robotLogContent, _ := ioutil.ReadFile("./" + r.Logname)
//	buf := bytes.NewBuffer(robotLogContent)
//
//	var decodedRobotLog RobotLog
//
//	decoder := gob.NewDecoder(buf)
//	err := decoder.Decode(&decodedRobotLog)
//	if err != nil {
//		panic(err)
//	}
//
//	r.RMap = decodedRobotLog.RMap
//	r.CurLocation = decodedRobotLog.CurLocation
//	r.CurrTask = decodedRobotLog.CurrTask
//	fmt.Println(decodedRobotLog.RMap)
//	fmt.Println(decodedRobotLog.CurLocation)
//	fmt.Println(decodedRobotLog.CurrTask)
//	fmt.Println("finshed loading from log")
//}
//
//func (r *RobotStruct) CreateLog() (*os.File, error) {
//	file, err := os.Create("./" + r.Logname)
//	if err != nil {
//		fmt.Println("error creating robot log")
//	}
//	return file, err
//}
//
//func (r *RobotStruct) ProduceLogInfo() RobotLog {
//	robotLog := RobotLog{
//		CurrTask:    r.CurrTask,
//		RMap:        r.RMap,
//		CurLocation: r.CurLocation,
//	}
//	return robotLog
//}
//
//func (r *RobotStruct) LocateLog() (*os.File, error) {
//	file, err := os.Open(r.Logname)
//	return file, err
//}

func (r *RobotStruct) UpdateStateForNewJourney() {

	r.RobotNeighbours.Lock()
	r.RobotNeighbours.rNeighbour = make(map[int]Neighbour)
	r.RobotNeighbours.Unlock()
	r.ReceivedTasks = make([]TaskPayload, 0)
	r.ReceivedTasksResponse = make([]TaskDescisionPayload, 0)

	ticker := time.NewTicker(1000 * time.Millisecond)
	temp := time.Now()

	r.exchangeFlag.Lock()
	r.exchangeFlag.flag = false
	r.exchangeFlag.Unlock()

	r.State.Lock()
	r.State.rState = ROAM
	r.State.Unlock()

	go func() {
		counter := 0
		for _ = range ticker.C {
			counter += 1
			fmt.Printf("Flag timer. \n        Counter is %s\n", counter)
			if time.Now().Sub(temp) >= TIMETOJOIN {
				fmt.Println("WE ARE FINISHED -- CAN JOIN AGAIN")
				ticker.Stop()
				r.exchangeFlag.Lock()
				r.exchangeFlag.flag = true
				r.exchangeFlag.Unlock()
			}
		}
	}()

}

func (r *RobotStruct) SendMapToLocalServer() {
	for {
		// Encode Map info
		buf := new(bytes.Buffer)
		encoder := gob.NewEncoder(buf)
		err := encoder.Encode(RandomMapGenerator())
		if err != nil {
			continue
		}
		fmt.Println("Encoded map")
		// output := string(buf.Bytes())
		// Send it to local Server using TCP
		conn, err := net.Dial("tcp", ":8888")
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			conn.Write(buf.Bytes())
			fmt.Println("----------------------------------------------")
			fmt.Println("Encoded map send")
			fmt.Println("----------------------------------------------")
		}
		time.Sleep(5000 * time.Millisecond)
	}
}

// FN: Removes the dead nieghbours from this robots list
// TODO: MIGHT need to make this a go-routine in the main.go. Then need to remove checkAllNeighbours
func (r *RobotStruct) RemoveDeadNeighbours() {
	r.RobotNeighbours.Lock()
	maxCallTime := 1 // for network issue problem check later
	var err error
	for i, nei := range r.RobotNeighbours.rNeighbour {
		for i := 0; i < maxCallTime; i++ {
			client, error := rpc.Dial("tcp", nei.Addr)
			err = error
			if err == nil {
				break
			}
			client.Close()

		}
		if err != nil {
			delete(r.RobotNeighbours.rNeighbour, i)
			err = nil
		}

	}
	r.RobotNeighbours.Unlock()
}

func (r *RobotStruct) MonitorButtonsOnPins(leftButton bgpio.Pin) {
	for {
		if leftValue, err := leftButton.Read(); leftValue == 1 {
			fmt.Println("left wall button pressed")
			if err != nil {
				fmt.Println(err)
			}
			r.LeftWallSig <- true
			time.Sleep(300 * time.Millisecond)
		}
	}
}

func (r *RobotStruct) MonitorRightButtonOnPin(rightButton bgpio.Pin) {
	for {
		if rightValue, err := rightButton.Read(); rightValue == 1 {
			fmt.Println("right wall button pressed")
			if err != nil {
				fmt.Println(err)
			}
			r.RightWallSig <- true
			time.Sleep(300 * time.Millisecond)
		}
	}

}

func (r *RobotStruct) MonitorFrontButtonOnPin(frontEmptyButton bgpio.Pin) {
	for {
		if frontEmptyValue, err := frontEmptyButton.Read(); frontEmptyValue == 1 {
			fmt.Println("font empty button pressed")
			if err != nil {
				fmt.Println(err)
			}
			r.FreeSpaceSig <- true
			time.Sleep(300 * time.Millisecond)
		}
	}

}

func (r *RobotStruct) MonitorFrontObsButtonOnPin(frontObsButton bgpio.Pin) {
	for {
		if frontObsValue, err := frontObsButton.Read(); frontObsValue == 1 {
			fmt.Println("front wall button pressed")
			if err != nil {
				fmt.Println(err)
			}
			r.WallSig <- true
			time.Sleep(300 * time.Millisecond)
		}
	}

}

// func (r *RobotStruct) MonitorLeftObstacleButton() {
// 	for {
// 		select {
// 		case:
//
// 		}
// 	}
// }
