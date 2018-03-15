package shared

import (
	"math"
)

func Round(n float64) float64 {
	i := int(n)
	if (n-float64(i)) >= 0.5 {
		return math.Ceil(n)
	} else {
		return math.Floor(n)
	}
}