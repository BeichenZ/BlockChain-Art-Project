package shared

import "math"
// TODO comment : i dont think we need this
type Map struct {
	ExploredPath map[Coordinate]PointStruct
	FrameOfRef   int // Robot id
}