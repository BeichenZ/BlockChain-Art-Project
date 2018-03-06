package shared


type Map struct {
	ExploredPath []PointStruct
	FrameOfRef   uint // Robot id
}
type PointStruct struct {
	Point     Coordinate
	PointKind bool // true - free space, false - wall
	TraversedTime int64
	Traversed bool
}
type Coordinate struct {
	X float64
	Y float64
}

type Path struct {
	ListOfPCoordinates []PointStruct;
}
