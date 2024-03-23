package vm

import (
	"math"
)

type vmMemory [math.MaxUint16 + 1]uint16

var memory vmMemory
