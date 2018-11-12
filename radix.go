package interval

import (
	"sort"
)

const (
	minSize      = 256
	radix   uint = 8
	bitSize uint = 64
)

// sortIntervalsStart sorts a slice of intervals.
func sortIntervalsStart(x []Interval) {
	if len(x) < 2 {
		return
	} else if len(x) < minSize {
		sort.Slice(x, func(i, j int) bool { return x[i].Start < x[j].Start })
	} else {
		doSortStart(x)
	}
}

func doSortStart(x []Interval) {
	from := x
	to := make([]Interval, len(x))
	var key uint8
	var offset [256]int

	for keyOffset := uint(0); keyOffset < bitSize; keyOffset += radix {
		keyMask := int64(0xFF << keyOffset)
		var counts [256]int
		sorted := true
		prev := int64(0)
		for _, elem := range from {
			key = uint8((elem.Start & keyMask) >> keyOffset)
			counts[key]++
			if sorted {
				sorted = elem.Start >= prev
				prev = elem.Start
			}
		}

		if sorted {
			if (keyOffset/radix)%2 == 1 {
				copy(to, from)
			}
			return
		}

		offset[0] = 0
		for i := 1; i < len(offset); i++ {
			offset[i] = offset[i-1] + counts[i-1]
		}
		for _, elem := range from {
			key = uint8((elem.Start & keyMask) >> keyOffset)
			to[offset[key]] = elem
			offset[key]++
		}
		to, from = from, to
	}
}

func sortIntervalsEnd(x []Interval) {
	if len(x) < 2 {
		return
	} else if len(x) < minSize {
		sort.Slice(x, func(i, j int) bool { return x[i].End < x[j].End })
	} else {
		doSortEnd(x)
	}
}

func doSortEnd(x []Interval) {
	from := x
	to := make([]Interval, len(x))
	var key uint8
	var offset [256]int

	for keyOffset := uint(0); keyOffset < bitSize; keyOffset += radix {
		keyMask := int64(0xFF << keyOffset)
		var counts [256]int
		sorted := true
		prev := int64(0)
		for _, elem := range from {
			key = uint8((elem.End & keyMask) >> keyOffset)
			counts[key]++
			if sorted {
				sorted = elem.End >= prev
				prev = elem.End
			}
		}

		if sorted {
			if (keyOffset/radix)%2 == 1 {
				copy(to, from)
			}
			return
		}

		offset[0] = 0
		for i := 1; i < len(offset); i++ {
			offset[i] = offset[i-1] + counts[i-1]
		}
		for _, elem := range from {
			key = uint8((elem.End & keyMask) >> keyOffset)
			to[offset[key]] = elem
			offset[key]++
		}
		to, from = from, to
	}
}
