package test

import (
	"testing"
	"../shared"
	//"fmt"
	//"math"
)
const NEIGHBOURS = 1
const NEIGHBOURPATH = 2

var WEST = shared.PointStruct{shared.Coordinate{-1.0, 0.0}, false, 0, false}
var EAST = shared.PointStruct{shared.Coordinate{1.0, 0.0}, false, 0, false}
var NORTH = shared.PointStruct{shared.Coordinate{0.0, 1.0}, false, 0, false}
var SOUTH = shared.PointStruct{shared.Coordinate{0.0, -1.0}, false, 0, false}



//func GetStartTestTitle(e string) string  {
//	return "<===================  Starting test case ["+ e + "]  ===================>\n"
//}
//
//func GetEndTestTitle(e string) string  {
//	return "<===================  Ending test case ["+ e + "] ===================>\n\n"
//}
//
//////User with non-empty map doing a merge
//func TestMapMerge_withNonEmptyMap(t *testing.T) {
//
//	fmt.Printf(GetStartTestTitle("MapMerge_withNonEmptyMap"))
//
//	var sampleMap = []shared.Map{}
//
//	exploredPath := make(map[shared.Coordinate]shared.PointStruct)
//	coordinate := shared.Coordinate{float64(1), float64(2)}
//	exploredPath[coordinate] = shared.PointStruct{
//		shared.Coordinate{float64(1), float64(2)},
//		false,
//		0,
//		false,
//	}
//
//	newMap := shared.Map{
//		ExploredPath:exploredPath,
//		FrameOfRef: 2,
//	}
//
//	var robot = shared.InitRobot(0, newMap, nil, "", "jjjh")
//
//	for i:= 0; i< NEIGHBOURS ; i++{
//		robId := i
//		myMap := new(shared.Map)
//		myMap.ExploredPath = make(map[shared.Coordinate] shared.PointStruct)
//		for j := 10; j< 10 + NEIGHBOURPATH; j++{
//			myPoint := shared.Coordinate{float64(i), float64(j)}
//			myMap.ExploredPath[myPoint] = shared.PointStruct{myPoint, false, 0, false}
//		}
//		myMap.FrameOfRef = robId
//		sampleMap = append(sampleMap, *myMap)
//	}
//
//	fmt.Printf("Map before merged ==> %v\n", robot.GetMap())
//	robot.MergeMaps(sampleMap)
//	fmt.Printf("Map after merged ==> %v\n", robot.GetMap())
//
//	fmt.Printf(GetEndTestTitle("MapMerge_withNonEmptyMap"))
//
//}
//
////User with empty doing a merge
//func TestMapMerge_withEmptyMap(t *testing.T) {
//	fmt.Printf(GetStartTestTitle("MapMerge_withEmptyMap"))
//
//	var sampleMap = []shared.Map{}
//
//	var robot = shared.InitRobot(0, shared.Map{}, nil, "", "343")
//
//	for i:= 0; i< NEIGHBOURS ; i++{
//		robId := i
//		myMap := new(shared.Map)
//		myMap.ExploredPath = make(map[shared.Coordinate]shared.PointStruct)
//		for j := 10; j< 10 + NEIGHBOURPATH; j++ {
//			myPoint := shared.Coordinate{float64(i), float64(j)}
//			myMap.ExploredPath[myPoint] = shared.PointStruct{myPoint, false, 0, false}
//		}
//		myMap.FrameOfRef = robId
//		sampleMap = append(sampleMap, *myMap)
//	}
//
//	fmt.Printf("Map before merged ==> %v\n", robot.GetMap())
//	robot.MergeMaps(sampleMap)
//	fmt.Printf("Map after merged ==> %v\n", robot.GetMap())
//
//	//shared.PrettyPrint_Map(robot.GetMap())
//
//	fmt.Println(GetEndTestTitle("MapMerge_withEmptyMap"))
//
//}
//
//func TestModifyPathForWall(t *testing.T){
//
//	fmt.Printf(GetStartTestTitle("ModifyPathForWall"))
//
//	robotStruct := shared.RobotStruct{}
//	sampleTask := []shared.PointStruct{}
//
//	//generate sample directions
//	for i:= 0; i< 4 ; i++{
//		sampleTask = append(sampleTask, EAST)
//	}
//
//	for i:= 0; i< 3 ; i++{
//		sampleTask = append(sampleTask, NORTH)
//	}
//
//	robotStruct.CurPath.ListOfPCoordinates = sampleTask
//
//	//finish generating sample data
//
//	fmt.Printf("Task before modified => %v\n", robotStruct.CurPath.ListOfPCoordinates)
//	//shared.PrettyPrint_Path(robotStruct.CurPath.ListOfPCoordinates)
//
//	robotStruct.ModifyPathForWall()
//
//	fmt.Printf("Task after modified => %v\n", robotStruct.CurPath.ListOfPCoordinates )
//	//shared.PrettyPrint_Path(robotStruct.CurPath.ListOfPCoordinates)
//
//	if len(robotStruct.CurPath.ListOfPCoordinates) != 3{
//		t.Errorf("The actual is %d but the expected value is %d",len(robotStruct.CurPath.ListOfPCoordinates), 3 )
//	}
//
//	fmt.Println(GetEndTestTitle("ModifyPathForWall"))
//
//}
//
//func TestTaskCreation(t *testing.T){
//	fmt.Printf(GetStartTestTitle("Task Creation"))
//
//	robotStruct := shared.RobotStruct{}
//	robotStruct.RMap = shared.RandomMapGenerator()
//	_ = shared.Neighbour{}
//	robotStruct.CurLocation = shared.Coordinate{float64(3.0), float64(4.0)}
//
//	task, _ :=robotStruct.TaskCreation()
//	fmt.Printf("The created task is %v\n", task)
//
//	path := shared.CreatePathBetweenTwoPoints(robotStruct.CurLocation, task[0].Point)
//	fmt.Println("The create path is ", path)
//	fmt.Println("The length of the path is ", len(path.ListOfPCoordinates))
//
//	fmt.Println(GetEndTestTitle("Task Creation"))
//}
//
//func TestPathCreation(t *testing.T){
//	fmt.Println(GetStartTestTitle("CreatePathBetweenTwoPoints"))
//
//	p1 := shared.Coordinate{0.0, 0.0}
//	p2 := shared.Coordinate{8.0, 8.0}
//	// robotStruct := shared.RobotStruct{}
//	path :=shared.CreatePathBetweenTwoPoints(p1, p2)
//
//	fmt.Println("The create path is ", path)
//	fmt.Println("The length of the path is ", len(path.ListOfPCoordinates))
//	fmt.Println(GetEndTestTitle("CreatePathBetweenTwoPoints"))
//}
////func TestTaskAllocation(t *testing.T){
////	fmt.Printf(GetStartTestTitle("Task Allocation"))
////
////	robotStruct := shared.RobotStruct{}
////	robotStruct.RMap = RandomMapGenerator()
////	rn := shared.Neighbour{}
////	robotStruct.RobotNeighbours = append(robotStruct.RobotNeighbours, rn, rn , rn)
////	robotStruct.CurLocation = shared.PointStruct{Point:shared.Coordinate{float64(3.0), float64(4.0)}}
////
////	tasks, _ :=robotStruct.TaskCreation()
////
////	robotStruct.AllocateTaskToNeighbours(tasks)
////
////}
//
//// TODO: TEST
//
//func TestFindDestPoints(t *testing.T) {
//
//	center := shared.Coordinate{0,0}
//	listOfPoints := shared.FindDestPoints(5, center)
//
//
//	for _, point := range listOfPoints {
//
//		calculatedR := math.Sqrt(point.Point.X*point.Point.X + point.Point.Y*point.Point.Y)
//		if math.Abs(calculatedR - shared.EXRADIUS) > 0.01 {
//			fmt.Println(math.Abs(calculatedR - shared.EXRADIUS))
//			t.FailNow()
//		}
//	}
//}

func createNeighbour(num1 float64, num2 float64) shared.Neighbour{
	return shared.Neighbour{NeighbourCoordinate:shared.Coordinate{num1, num2}}
}

func TestJoinState(t *testing.T){

	robot1 := shared.InitRobot(0, shared.Map{}, nil, "", "")
	robot2 := shared.InitRobot(1, shared.Map{}, nil, "", "")

	robot1.CurLocation = shared.Coordinate{0.0,0.0}
	robot2.CurLocation = shared.Coordinate{5.0,5.0}

	payload:= shared.FarNeighbourPayload{NeighbourCoordinate: robot2.CurLocation, ItsNeighbours: []shared.Neighbour{}}

	result :=robot1.WithinRadiusOfNetwork(&payload)

	if result{
		t.FailNow()
	}

	robot2.CurLocation = shared.Coordinate{5.0,0.0}
	payload = shared.FarNeighbourPayload{NeighbourCoordinate: robot2.CurLocation, ItsNeighbours: []shared.Neighbour{}}

	result =robot1.WithinRadiusOfNetwork(&payload)

	if !result{
		t.FailNow()
	}

}
// case2: caller, callee within radius, but caller neighbour not within radius
func TestJoinState_case2(t *testing.T){

	robot1 := shared.InitRobot(0, shared.Map{}, nil, "", "")
	robot2 := shared.InitRobot(1, shared.Map{}, nil, "", "")
	robot1.CurLocation = shared.Coordinate{0.0,0.0}
	robot2.CurLocation = shared.Coordinate{5.0,0.0}

	robot1.RobotNeighbours[0] = createNeighbour(-5.0,3.0)

	//robot2.RobotNeighbours[0] = createNeighbour(10.0,10.0)

	payload:= shared.FarNeighbourPayload{NeighbourCoordinate: robot2.CurLocation, ItsNeighbours: []shared.Neighbour{}}
	result :=robot1.WithinRadiusOfNetwork(&payload)

	if result {
		t.FailNow()
	}


}
// case3: caller, callee within radius, but both caller and callees neighbour not within radius
func TestJoinState_case3(t *testing.T){

	robot1 := shared.InitRobot(0, shared.Map{}, nil, "", "")
	robot2 := shared.InitRobot(1, shared.Map{}, nil, "", "")
	robot1.CurLocation = shared.Coordinate{0.0,0.0}
	robot2.CurLocation = shared.Coordinate{5.0,0.0}

	robot1.RobotNeighbours[0] = createNeighbour(-5.0,3.0)

	robot2.RobotNeighbours[0] = createNeighbour(10.0,3.0)

	payload:= shared.FarNeighbourPayload{NeighbourCoordinate: robot2.CurLocation, ItsNeighbours: []shared.Neighbour{}}
	result :=robot1.WithinRadiusOfNetwork(&payload)

	if result {
		t.FailNow()
	}


}
// case4: everyone is within radius
func TestJoinState_case4(t *testing.T){

	robot1 := shared.InitRobot(0, shared.Map{}, nil, "", "")
	robot2 := shared.InitRobot(1, shared.Map{}, nil, "", "")
	robot1.CurLocation = shared.Coordinate{0.0,0.0}
	robot2.CurLocation = shared.Coordinate{5.0,0.0}

	robot1.RobotNeighbours[0] = createNeighbour(1.0,3.0)

	robot2.RobotNeighbours[0] = createNeighbour(2.0,3.0)

	payload:= shared.FarNeighbourPayload{NeighbourCoordinate: robot2.CurLocation, ItsNeighbours: []shared.Neighbour{}}
	result :=robot1.WithinRadiusOfNetwork(&payload)

	if !result {
		t.FailNow()
	}
// if time permits do stress test lols
}

// FindDestPoints() - All the destination points are unique
// UpdateMap() - error free case gets the right signal
// UpdateMap() - error case
// UpdateMap() - point is in the hash map already
// UpdateMap() - point isn't in the hashmap
// UpdateCurentLocation()