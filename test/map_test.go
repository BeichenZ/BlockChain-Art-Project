package test

import (
	"testing"
	"../shared"
)
const NEIGHBOURS = 1
const NEIGHBOURPATH = 2

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