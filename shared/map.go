package shared

import "math"

type Map struct {
	ExploredPath []PointStruct
	FrameOfRef   int // Robot id
}