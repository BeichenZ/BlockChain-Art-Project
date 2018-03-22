package sharedsis
import "math"

type PointStruct struct {
	Point         Coordinate
	PointKind     bool
	TraversedTime int64
	Traversed     bool
}
type Coordinate struct {
	X float64
	Y float64
}
// Button
// 1 - free space at current point,
// 2 - wall at current point,
// 3 - wall RIGHT to current point
// 4 - wall LEFT to current point
type Button int
const(
	FreeSpace Button = 1
	Wall Button = 2
	RightWall Button = 3
	LeftWall Button = 4
)
var WEST = PointStruct{Coordinate{-1.0, 0.0}, false, 0, false}
var EAST = PointStruct{Coordinate{1.0, 0.0}, false, 0, false}
var NORTH = PointStruct{Coordinate{0.0, 1.0}, false, 0, false}
var SOUTH = PointStruct{Coordinate{0.0, -1.0}, false, 0, false}

// ----------------------------------------- FUNCTIONS ---------------------------------------------------------- //
// FN: Finds magnitiude of the distance btwn two points
func DistBtwnTwoPoints(dp PointStruct, cp PointStruct) float64 {
	d := math.Sqrt(math.Pow(dp.Point.X-cp.Point.X, 2) + math.Pow(dp.Point.Y-cp.Point.Y, 2))
	return d
}

// FN: Check if coordinate of point exists in point Array
func CheckExist(coordinate PointStruct, cooArr []PointStruct) (bool, int) {
	for i, point := range cooArr {
		if point.Point.X == coordinate.Point.X && point.Point.Y == coordinate.Point.Y {
			return true, i
		}
	}
	return false, -1
}
// FN: Remove point from list of points
func removeElFromlist(p PointStruct, listp *[]PointStruct) {

	copyInput := *listp
	for idx, val := range copyInput{

		if val == p{
			temp1:= copyInput[:idx]
			temp2:= copyInput[idx:]
			if len(temp2) == 1{
				*listp = temp1
				return
			}else{
				temp2 = temp2[1:]
				*listp = append(temp1, temp2...)
				return
			}
		}
	}
}
//FN: Return list of destNum destination points EXRADIUS away from the given center
func FindDestPoints(desNum int, center PointStruct) []PointStruct {

	destPointsToReturn := []PointStruct{}

	for i := 0; i < desNum; i++ {
		theta := float64(i) * 2 * math.Pi / float64(desNum)
		delPoint := PointStruct{Point: Coordinate{float64(EXRADIUS * math.Cos(theta)), float64(EXRADIUS * math.Sin(theta))}}
		destPoint := PointStruct{}
		destPoint.Point.X = center.Point.X + delPoint.Point.X
		destPoint.Point.Y = center.Point.Y + delPoint.Point.Y
		destPointsToReturn = append(destPointsToReturn, destPoint)
	}

	return destPointsToReturn
}
// TODO
func CompareCoordinateTimeStamp(t1 int, t2 int ) bool {
	return false;
}
