package shared

import (
	"fmt"
	"net/rpc"
	"time"
	"math"
)

const XMIN  = "xmin"
const XMAX  = "xmax"
const YMIN  = "ymix"
const YMAX  = "ymax"
const EXRADIUS = 6


type RobotStruct struct {
	RobotID         uint // hardcoded
	RobotIP         string
	RobotListenConn *rpc.Client
	RobotNeighbourNum int
	RMap            Map
	CurPath         Path
	CurLocation     PointStruct
	//CurrentStep        	Coordinate
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
	SendFreeSpaceSig()
}

var robotStruct RobotStruct

func (r *RobotStruct) SendMyMap(rId uint, rMap Map) {
	return
}

func (r *RobotStruct) SendFreeSpaceSig(){
	fmt.Println("got here")
	r.FreeSpaceSig <- true
}


//error is not nil when the task queue is empty
// FN: Return list of destination points for each node in the network (one point for each node)
func (r *RobotStruct) TaskCreation() ([]PointStruct, error) {
	//newTask := Path{}
	xmin := r.FindMapExtrema(XMIN)
	xmax := r.FindMapExtrema(XMAX)
	ymin := r.FindMapExtrema(YMIN)
	ymax := r.FindMapExtrema(YMAX)

	center := PointStruct{Point: Coordinate{float64((xmax-xmin)/2), float64((ymax-ymin)/2)}}

	DestNum := r.RobotNeighbourNum + 1;

	DestPoints := r.FindDestPoints(DestNum, center)

	DestPointForMe := r.FindClosestDest(DestPoints)
	// move DestpointForMe to beginning of list
	//assuming  the destPoint in DestPoints is unique

	tempEle := DestPoints[0]
	for idx, value := range DestPoints{
		if value == DestPointForMe{
			DestPoints[0] = value
			DestPoints[idx] = tempEle
			break
		}
	}

	return DestPoints, nil
}

func (r *RobotStruct) FindMapExtrema(e string) float64{

	if e == XMAX{
		var xMax float64 = math.MinInt64
		for _, point := range r.RMap.ExploredPath{
			if xMax < point.Point.X{
				xMax = point.Point.X
			}
		}
		return xMax
	}else if e == XMIN{
		var xMin float64 = math.MaxFloat64
		for _, point := range r.RMap.ExploredPath{
			if xMin > point.Point.X{
				xMin = point.Point.X
			}
		}
		return xMin
	}else if e == YMAX{
		var yMax float64 = math.MinInt64
		for _, point := range r.RMap.ExploredPath{
			if yMax < point.Point.Y{
				yMax = point.Point.Y
			}
		}
		return yMax
	}else{
		var yMin float64 = math.MaxFloat64
		for _, point := range r.RMap.ExploredPath{
			if yMin > point.Point.Y{
				yMin = point.Point.Y
			}
		}
		return yMin
	}
}

func (r *RobotStruct) FindClosestDest(lodp []PointStruct) PointStruct {
	dist := math.MaxFloat64
	var rdp PointStruct
	for _, dp := range lodp{
		if (DistBtwnTwoPoints(r.CurLocation, dp) < dist) {
			rdp = dp;
		}
	}
	return rdp;
}


func (r *RobotStruct) CreatePathBetweenTwoPoints(sp PointStruct, dp PointStruct) Path  {
	var myPath []PointStruct
	delX := dp.Point.X - sp.Point.X
	delY := dp.Point.Y - sp.Point.Y
	//iteration := int(math.Abs(delX) + math.Abs(delY))

	//create the path in X direction
	for i:= 0; i< int(math.Abs(delX)) ; i++{
		if delX > 0{
			myPath = append(myPath, PointStruct{Point: Coordinate{1, 0}})
		}else if delX < 0 {
			myPath = append(myPath, PointStruct{Point: Coordinate{-1, 0}})
		}else{
			//do nonthing since the delX is 0
		}
	}

	//create path in Y direction
	for i:= 0; i< int(math.Abs(delY)) ; i++{
		if delY > 0{
			myPath = append(myPath, PointStruct{Point: Coordinate{0, 1}})
		}else if delY < 0 {
			myPath = append(myPath, PointStruct{Point: Coordinate{0, -1}})
		}else{
			//do nonthing since the delY is 0
		}
	}

	return Path{myPath}
}

//return the list of dest points
func (r *RobotStruct) FindDestPoints(desNum int, center PointStruct) []PointStruct{

	destPointsToReturn := []PointStruct{}
	circumference := 2*math.Pi*EXRADIUS;
	arcLength := circumference / float64(desNum);

	for i := 0; i< desNum ; i++{
		theta := float64(i)*arcLength / EXRADIUS;
		delPoint := PointStruct{Point: Coordinate{float64(EXRADIUS*math.Cos(theta)) , float64(EXRADIUS*math.Sin(theta))}}
		destPoint := PointStruct{}
		destPoint.Point.X = center.Point.X + delPoint.Point.X
		destPoint.Point.Y = center.Point.Y + delPoint.Point.Y
		destPointsToReturn = append(destPointsToReturn, destPoint)
	}

	return destPointsToReturn
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
				dpts, err := r.TaskCreation()
				var newPath Path

				if (len(dpts) == 1) {

					//TODO
					newPath = r.CreatePathBetweenTwoPoints(r.CurLocation, dpts[0])

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

				r.UpdateMap(true)
				r.SetCurrentLocation()
				r.TookOneStep() //remove the first element from r.CurPath.ListOfPCoordinates

				// Display task with GPIO
			case <-r.WallSig:
				r.UpdateMap(false)
				// Change wall path
				r.ModifyPathForWall()
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
		if (wallCoor == c ){
			continue
		}
		r.CurPath.ListOfPCoordinates = r.CurPath.ListOfPCoordinates[i:]
		break
	}
}

func (r *RobotStruct) TookOneStep() {
	r.CurPath.ListOfPCoordinates = r.CurPath.ListOfPCoordinates[1:]
}

//update pointkind: true => freespace, false  => wall
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
	robotStruct.JoiningSig = make(chan bool)
	robotStruct.BusySig = make(chan bool)
	robotStruct.WaitingSig = make(chan bool)
	robotStruct.WallSig = make(chan bool)
	robotStruct.WalkSig = make(chan bool)
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
