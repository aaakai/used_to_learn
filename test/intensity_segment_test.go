package test

import (
	"testing"
)

func TestIntensitySegments(t *testing.T) {
	is1 := NewIntensitySegment()
	assertFailure(t, is1.toPrintData(), "[]")

	is1.add(10, 30, 1)
	assertFailure(t, is1.toPrintData(), "[[10 1] [30 0]]")

	is1.add(20, 40, 1)
	assertFailure(t, is1.toPrintData(), "[[10 1] [20 2] [30 1] [40 0]]")

	is1.add(10, 40, -2)
	assertFailure(t, is1.toPrintData(), "[[10 -1] [20 0] [30 -1] [40 0]]")

	//
	is2 := NewIntensitySegment()
	assertFailure(t, is2.toPrintData(), "[]")

	is2.add(10, 30, 1)
	assertFailure(t, is2.toPrintData(), "[[10 1] [30 0]]")

	is2.add(20, 40, 1)
	assertFailure(t, is2.toPrintData(), "[[10 1] [20 2] [30 1] [40 0]]")

	is2.add(10, 40, -1)
	assertFailure(t, is2.toPrintData(), "[[20 1] [30 0]]")

	is2.add(10, 40, -1)
	assertFailure(t, is2.toPrintData(), "[[10 -1] [20 0] [30 -1] [40 0]]")
}

func assertFailure(t *testing.T, a, b string) {
	if a != b {
		t.Errorf("\nActually: %s\nExpected: %s\n", a, b)
	}
}
