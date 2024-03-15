package terminal

import "fmt"

type coordinates struct {
	intitialX     int
	intitialY     int
	currentX      int
	currentY      int
	pendingDeltaX int
	pendingDeltaY int
}

func newCoords(x int, y int) coordinates {
	return coordinates{
		intitialX:     x,
		intitialY:     y,
		currentX:      x,
		currentY:      y,
		pendingDeltaX: 0,
		pendingDeltaY: 0,
	}
}

func (coord coordinates) String() string {
	return fmt.Sprintf(
		"\nInitial X: %v, Y: %v\nCurrent X: %v, Y: %v\nTarget X: %v, Y:%v\n",
		coord.intitialX,
		coord.intitialY,
		coord.currentX,
		coord.currentY,
		coord.pendingDeltaX,
		coord.pendingDeltaY,
	)
}

func (coord coordinates) onThisRow(i int) bool {
	return coord.currentY == i
}

func (coord coordinates) addRowDelta(i int) coordinates {
	return coord.setPendingRow(coord.pendingDeltaY + i)
}

func (coord coordinates) setPendingRow(i int) coordinates {
	coord.pendingDeltaY = i
	return coord
}

func (coord coordinates) addColumnDelta(i int) coordinates {
	return coord.setPendingColumn(coord.pendingDeltaX + i)
}

func (coord coordinates) setPendingColumn(i int) coordinates {
	coord.pendingDeltaX = i
	return coord
}

func (coord coordinates) outputDelataToTarget() (int, int) {
	return coord.pendingDeltaX - coord.currentX, coord.pendingDeltaY - coord.currentY
}

func (coord coordinates) applyPendingDeltas() coordinates {
	return coordinates{
		intitialX:     coord.intitialX,
		intitialY:     coord.intitialY,
		currentX:      coord.pendingDeltaX,
		currentY:      coord.pendingDeltaY,
		pendingDeltaX: coord.pendingDeltaX,
		pendingDeltaY: coord.pendingDeltaY,
	}
}
