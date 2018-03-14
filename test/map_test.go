package test

import (
	"testing"
	"../shared"
	"log"
)
const NEIGHBOURS = 1
const NEIGHBOURPATH = 2
const NUMOFPATHTOGENERATE = 10

var WEST = shared.PointStruct{shared.Coordinate{-1.0, 0.0}, false, 0, false}
var EAST = shared.PointStruct{shared.Coordinate{1.0, 0.0}, false, 0, false}
var NORTH = shared.PointStruct{shared.Coordinate{0.0, 1.0}, false, 0, false}
var SOUTH = shared.PointStruct{shared.Coordinate{0.0, -1.0}, false, 0, false}

func RandomMapGenerator() shared.Map{
	var sampleMap = shared.Map{}
	for j := 0; j< NUMOFPATHTOGENERATE; j++{
		myPoint := shared.Coordinate{float64(j%4), float64(j)}
		sampleMap.ExploredPath = append(sampleMap.ExploredPath, shared.PointStruct{myPoint, false, 0, false})
	}

	return sampleMap
}

//User with non-empty map doing a merge
func TestMapMerge_withNonEmptyMap(t *testing.T) {
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
		FrameOfRef: uint(2),
	}

	var robot = shared.InitRobot(0, newMap)

	for i:= 0; i< NEIGHBOURS ; i++{
		robId := i
		myMap := new(shared.Map)
		for j := 10; j< 10 + NEIGHBOURPATH; j++{
			myPoint := shared.Coordinate{float64(i), float64(j)}
			myMap.ExploredPath = append(myMap.ExploredPath, shared.PointStruct{myPoint, false, 0, false})
		}
		myMap.FrameOfRef = uint(robId)
		sampleMap = append(sampleMap, *myMap)
	}

	_ = robot.MergeMaps(sampleMap)
	shared.PrettyPrint_Map(robot.GetMap())

}

//User with empty doing a merge
func TestMapMerge_withEmptyMap(t *testing.T) {
	var sampleMap = []shared.Map{}

	var robot = shared.InitRobot(0, shared.Map{})

	for i:= 0; i< NEIGHBOURS ; i++{
		robId := i
		myMap := new(shared.Map)
		for j := 10; j< 10 + NEIGHBOURPATH; j++{
			myPoint := shared.Coordinate{float64(i), float64(j)}
			myMap.ExploredPath = append(myMap.ExploredPath, shared.PointStruct{myPoint, false, 0, false})
		}
		myMap.FrameOfRef = uint(robId)
		sampleMap = append(sampleMap, *myMap)
	}

	_ = robot.MergeMaps(sampleMap)
	shared.PrettyPrint_Map(robot.GetMap())
	return
}

func TestModifyPathForWall(t *testing.T){

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

	log.Printf("Task before modified =>")
	shared.PrettyPrint_Path(robotStruct.CurPath.ListOfPCoordinates)

	robotStruct.ModifyPathForWall()

	log.Printf("Task after modified =>")
	shared.PrettyPrint_Path(robotStruct.CurPath.ListOfPCoordinates)

	if len(robotStruct.CurPath.ListOfPCoordinates) != 3{
		t.Errorf("The actual is %d but the expected value is %d",len(robotStruct.CurPath.ListOfPCoordinates), 3 )
	}
}

func TestTaskCreation(t *testing.T){
	robotStruct := shared.RobotStruct{}
	robotStruct.RMap = RandomMapGenerator()
	robotStruct.CurLocation = shared.PointStruct{Point:shared.Coordinate{float64(3.0), float64(4.0)}}

	task, error :=robotStruct.TaskCreation()

}
