package shared

import "math"

type Path struct {
	ListOfPCoordinates []PointStruct
}

// ----------------------------------------- FUNCTIONS ---------------------------------------------------------- //
func CreatePathBetweenTwoPoints(sp PointStruct, dp PointStruct) Path {
	var myPath []PointStruct
	delX := Round(dp.Point.X - sp.Point.X)
	delY := Round(dp.Point.Y - sp.Point.Y)
	//iteration := int(math.Abs(delX) + math.Abs(delY))

	//create the path in X direction
	for i := 0; i < int(math.Abs(delX)); i++ {
		if delX > 0 {
			myPath = append(myPath, PointStruct{Point: Coordinate{1, 0}})
		} else if delX < 0 {
			myPath = append(myPath, PointStruct{Point: Coordinate{-1, 0}})
		} else {
			//do nonthing since the delX is 0
		}
	}

	//create path in Y direction
	for i := 0; i < int(math.Abs(delY)); i++ {
		if delY > 0 {
			myPath = append(myPath, PointStruct{Point: Coordinate{0, 1}})
		} else if delY < 0 {
			myPath = append(myPath, PointStruct{Point: Coordinate{0, -1}})
		} else {
			//do nonthing since the delY is 0
		}
	}

	return Path{myPath}
}
