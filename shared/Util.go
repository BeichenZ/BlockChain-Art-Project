package shared

import (
	"math"
)

const NUMOFPATHTOGENERATE = 10

func Round(n float64) float64 {
	i := int(n)
	if (n - float64(i)) >= 0.5 {
		return math.Ceil(n)
	} else {
		return math.Floor(n)
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