package test

import (
	"testing"
	"../shared"
	"fmt"
)
const NEIGHBOURS = 1
const NEIGHBOURPATH = 2
const NUMOFPATHTOGENERATE = 10

var WEST = shared.PointStruct{shared.Coordinate{-1.0, 0.0}, false, 0, false}
var EAST = shared.PointStruct{shared.Coordinate{1.0, 0.0}, false, 0, false}
var NORTH = shared.PointStruct{shared.Coordinate{0.0, 1.0}, false, 0, false}
var SOUTH = shared.PointStruct{shared.Coordinate{0.0, -1.0}, false, 0, false}

// xMin=yMin=0, xMax=yMax=9
func RandomMapGenerator() shared.Map{
	var sampleMap = shared.Map{}
	for j := 0; j< NUMOFPATHTOGENERATE; j++{
		myPoint := shared.Coordinate{float64(j), float64(j)}
		sampleMap.ExploredPath = append(sampleMap.ExploredPath, shared.PointStruct{myPoint, false, 0, false})
	}

	return sampleMap
}

func GetStartTestTitle(e string) string  {
	return "<===================  Starting test case ["+ e + "]  ===================>\n"
}

func GetEndTestTitle(e string) string  {
	return "<===================  Ending test case ["+ e + "] ===================>\n\n"
}

////User with non-empty map doing a merge
func TestMapMerge_withNonEmptyMap(t *testing.T) {

	fmt.Printf(GetStartTestTitle("MapMerge_withNonEmptyMap"))


	var sampleMap = []shared.Map{}
	newMap := shared.Map{
		ExploredPath: []shared.PointStruct{
			shared.PointStruct{
				shared.Coordinate{float64(1), float64(2)},
				false,
				0,
				false,
			},
		},
		FrameOfRef: 2,
	}

	var robot = shared.InitRobot(0, newMap, nil)

	for i:= 0; i< NEIGHBOURS ; i++{
		robId := i
		myMap := new(shared.Map)
		for j := 10; j< 10 + NEIGHBOURPATH; j++{
			myPoint := shared.Coordinate{float64(i), float64(j)}
			myMap.ExploredPath = append(myMap.ExploredPath, shared.PointStruct{myPoint, false, 0, false})
		}
		myMap.FrameOfRef = robId
		sampleMap = append(sampleMap, *myMap)
	}

	fmt.Printf("Map before merged ==> %v\n", robot.GetMap())
	_ = robot.MergeMaps(sampleMap)
	fmt.Printf("Map after merged ==> %v\n", robot.GetMap())

	fmt.Printf(GetEndTestTitle("MapMerge_withNonEmptyMap"))

}

//User with empty doing a merge
func TestMapMerge_withEmptyMap(t *testing.T) {
	fmt.Printf(GetStartTestTitle("MapMerge_withEmptyMap"))

	var sampleMap = []shared.Map{}

	var robot = shared.InitRobot(0, shared.Map{}, nil)

	for i:= 0; i< NEIGHBOURS ; i++{
		robId := i
		myMap := new(shared.Map)
		for j := 10; j< 10 + NEIGHBOURPATH; j++{
			myPoint := shared.Coordinate{float64(i), float64(j)}
			myMap.ExploredPath = append(myMap.ExploredPath, shared.PointStruct{myPoint, false, 0, false})
		}
		myMap.FrameOfRef = robId
		sampleMap = append(sampleMap, *myMap)
	}

	fmt.Printf("Map before merged ==> %v\n", robot.GetMap())
	_ = robot.MergeMaps(sampleMap)
	fmt.Printf("Map after merged ==> %v\n", robot.GetMap())

	//shared.PrettyPrint_Map(robot.GetMap())

	fmt.Println(GetEndTestTitle("MapMerge_withEmptyMap"))

}

func TestModifyPathForWall(t *testing.T){

	fmt.Printf(GetStartTestTitle("ModifyPathForWall"))

	robotStruct := shared.RobotStruct{}
	sampleTask := []shared.PointStruct{}

	//generate sample directions
	for i:= 0; i< 4 ; i++{
		sampleTask = append(sampleTask, EAST)
	}

	for i:= 0; i< 3 ; i++{
		sampleTask = append(sampleTask, NORTH)
	}

	robotStruct.CurPath.ListOfPCoordinates = sampleTask

	//finish generating sample data

	fmt.Printf("Task before modified => %v\n", robotStruct.CurPath.ListOfPCoordinates)
	//shared.PrettyPrint_Path(robotStruct.CurPath.ListOfPCoordinates)

	robotStruct.ModifyPathForWall()

	fmt.Printf("Task after modified => %v\n", robotStruct.CurPath.ListOfPCoordinates )
	//shared.PrettyPrint_Path(robotStruct.CurPath.ListOfPCoordinates)

	if len(robotStruct.CurPath.ListOfPCoordinates) != 3{
		t.Errorf("The actual is %d but the expected value is %d",len(robotStruct.CurPath.ListOfPCoordinates), 3 )
	}

	fmt.Println(GetEndTestTitle("ModifyPathForWall"))

}

func TestTaskCreation(t *testing.T){
	fmt.Printf(GetStartTestTitle("Task Creation"))

	robotStruct := shared.RobotStruct{}
	robotStruct.RMap = RandomMapGenerator()
	rn := shared.Neighbour{}
	robotStruct.RobotNeighbours = append(robotStruct.RobotNeighbours, rn, rn , rn)
	robotStruct.CurLocation = shared.PointStruct{Point:shared.Coordinate{float64(3.0), float64(4.0)}}

	task, _ :=robotStruct.TaskCreation()
	fmt.Printf("The created task is %v\n", task)

	path := shared.CreatePathBetweenTwoPoints(robotStruct.CurLocation, task[0])
	fmt.Println("The create path is ", path)
	fmt.Println("The length of the path is ", len(path.ListOfPCoordinates))

	fmt.Println(GetEndTestTitle("Task Creation"))
}

func TestPathCreation(t *testing.T){
	fmt.Println(GetStartTestTitle("CreatePathBetweenTwoPoints"))

	p1 :=shared.PointStruct{Point:shared.Coordinate{0.0, 0.0}}
	p2 :=shared.PointStruct{Point:shared.Coordinate{8.0, 8.0}}
	// robotStruct := shared.RobotStruct{}
	path :=shared.CreatePathBetweenTwoPoints(p1, p2)

	fmt.Println("The create path is ", path)
	fmt.Println("The length of the path is ", len(path.ListOfPCoordinates))
	fmt.Println(GetEndTestTitle("CreatePathBetweenTwoPoints"))
}
//func TestTaskAllocation(t *testing.T){
//	fmt.Printf(GetStartTestTitle("Task Allocation"))
//
//	robotStruct := shared.RobotStruct{}
//	robotStruct.RMap = RandomMapGenerator()
//	rn := shared.Neighbour{}
//	robotStruct.RobotNeighbours = append(robotStruct.RobotNeighbours, rn, rn , rn)
//	robotStruct.CurLocation = shared.PointStruct{Point:shared.Coordinate{float64(3.0), float64(4.0)}}
//
//	tasks, _ :=robotStruct.TaskCreation()
//
//	robotStruct.AllocateTaskToNeighbours(tasks)
//
//}

// TODO: TEST
// FindDestPoints() - All the destination points are unqie
// UpdateMap() - error free case gets the right signal
// UpdateMap() - error case
// UpdateMap() - point is in the hash map already
// UpdateMap() - point isn't in the hashmap