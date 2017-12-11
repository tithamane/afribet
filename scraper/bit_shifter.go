package main

import "errors"

const (
	ShiftedBitsExceedTotalBits = "Bits to be shifted exceed the total bits available"
)

// BitShifter makes it easier to generate possible combination  by shifting bits
type BitShifter struct {
	values []int
	bits   []bool
}

func NewBitShifter(values []int) (BitShifter, error) {

	bitShifter := BitShifter{
		values: values,
		bits:   make([]bool, len(values)),
	}
	return bitShifter, nil
}

func (b *BitShifter) resetBits(activeBits int) {
	for i := 0; i < len(b.bits); i++ {
		if i < activeBits {
			b.bits[i] = true
		} else {
			b.bits[i] = false
		}
	}
}

func (b *BitShifter) CombinationsN(n int) ([][]int, error) {
	if n > len(b.bits) {
		return [][]int{}, errors.New(ShiftedBitsExceedTotalBits)
	}
	b.resetBits(n)
	combinations := [][]int{}

	// Copy the initial state too
	combinations = append(combinations, b.CombinationState())

	for b.canShift(n) {
		if b.canShiftWithoutCarryOver() {
			b.shift()
		} else {
			b.shiftCarryOver()
		}
		combinations = append(combinations, b.CombinationState())
	}

	return combinations, nil
}

func (b *BitShifter) CombinationState() []int {
	values := []int{}
	for i := 0; i < len(b.bits); i++ {
		if b.bits[i] == true {
			values = append(values, b.values[i])
		}
	}

	return values
}

// Checks if the last nBits are all true, if so, we can't shift a bit anymore
func (b *BitShifter) canShift(shiftableBits int) bool {
	lastIndex := len(b.bits) - 1
	for i := lastIndex; i > lastIndex-shiftableBits; i-- {
		if b.bits[i] == false {
			return true
		}
	}

	return false
}

func (b *BitShifter) canShiftWithoutCarryOver() bool {
	lastIndex := len(b.bits) - 1
	if b.bits[lastIndex] == false {
		return true
	}

	return false
}

func (b *BitShifter) shift() {
	lastIndex := len(b.bits) - 1
	for i := lastIndex; i > -1; i-- {
		if b.bits[i] == true {
			b.bits[i] = false
			b.bits[i+1] = true
			break
		}
	}
}

// shiftCarryOver is only called if lastBit is checked
func (b *BitShifter) shiftCarryOver() {
	lastAvailableIndex := b.lastAvailableIndex()
	// shift active bit nearest and before lastAvailableIndex forward
	var newShiftedBitIndex int
	for i := lastAvailableIndex; i > -1; i-- {
		if b.bits[i] == true {
			b.bits[i] = false
			newShiftedBitIndex = i + 1
			b.bits[newShiftedBitIndex] = true
			break
		}
	}

	// Count active bits and reset them til the newShiftedIndex
	var activeBits int
	for i := newShiftedBitIndex + 1; i < len(b.bits); i++ {
		if b.bits[i] == true {
			activeBits++
		}
		b.bits[i] = false
	}

	// Set the active bits immediately after the newShiftedBitIndex
	lastActiveBitIndex := newShiftedBitIndex + 1 + activeBits
	for i := newShiftedBitIndex + 1; i < lastActiveBitIndex; i++ {
		b.bits[i] = true
	}
}

// Returns the last bitIndex where the value is false
func (b *BitShifter) lastAvailableIndex() int {
	lastIndex := len(b.bits) - 1
	for i := lastIndex; i > -1; i-- {
		if b.bits[i] == false {
			return i
		}
	}

	return 0
}
