package shared

import "math"

type Map struct {
	ExploredPath []PointStruct
	FrameOfRef   uint // Robot id
}
type PointStruct struct {
	Point         Coordinate
	PointKind     bool // true - free space, false - wall
	TraversedTime int64
	Traversed     bool
}
type Coordinate struct {
	X float64
	Y float64
}

type Path struct {
	ListOfPCoordinates []PointStruct
}

func DistBtwnTwoPoints(dp PointStruct, cp PointStruct) float64 {
	d := math.Sqrt(math.Pow(dp.Point.X-cp.Point.X, 2) + math.Pow(dp.Point.Y-cp.Point.Y, 2))
	return d
}
