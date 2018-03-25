package shared

import (
	"math"
)

const NUMOFPATHTOGENERATE = 10

func Round(n float64) float64 {
	i := int(n*1000)
	if (n*1000 - float64(i)) >= 0.5 {
		return math.Ceil(float64(i))/1000
	} else {
		return math.Floor(float64(i))/1000
	}
}


// xMin=yMin=0, xMax=yMax=9
func RandomMapGenerator() Map{
	var sampleMap = Map{}
	sampleMap.ExploredPath = make(map[Coordinate] PointStruct)
	for j := 0; j< NUMOFPATHTOGENERATE; j++{
		myPoint := Coordinate{float64(j), float64(j)}
		sampleMap.ExploredPath[myPoint] = PointStruct{myPoint, false, 0, false}
	}

	return sampleMap
}