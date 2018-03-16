package shared

import (
	"math"
)

func Round(n float64) float64 {
	i := int(n)
	if (n - float64(i)) >= 0.5 {
		return math.Ceil(n)
	} else {
		return math.Floor(n)
	}
}

func CheckExist(coordinate PointStruct, cooArr []PointStruct) (bool, int) {
	for i, point := range cooArr {
		if point.Point.X == coordinate.Point.X && point.Point.Y == coordinate.Point.Y {
			return true, i
		}
	}
	return false, -1
}
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
