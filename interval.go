package interval

import "fmt"

// Tree is an interval tree that allows efficient querying of which
// intervals contain a given point.
type Tree struct {
	left  *Tree // all intervals to the left
	right *Tree // all intervals to the right

	center int64
	starts []Interval
	ends   []Interval
}

// Interval represents an interval with an attached key.
type Interval struct {
	Start int64
	End   int64
	Key   int64
}

// New constructs an interval tree from the set of intervals.
func New(ints []Interval) *Tree {
	if len(ints) == 0 {
		return new(Tree)
	}

	center := findCenter(ints)
	if len(ints) <= 40 {
		return &Tree{
			center: center,
			starts: sortByStart(ints),
			ends:   sortByEnd(ints),
		}
	}

	left, right := split(ints, center)
	return &Tree{
		left:  New(ints[:left]),
		right: New(ints[right:]),

		center: center,
		starts: sortByStart(ints[left:right]),
		ends:   sortByEnd(ints[left:right]),
	}
}

// Size returns the length of the nodes that match the center
func (t *Tree) Size() int {
	if t == nil {
		return 0
	}
	return len(t.starts)
}

// String returns a string summary of the tree.
func (t *Tree) String() string {
	return fmt.Sprintf("Tree(center:%d, size:%d, left:%d, right:%d)",
		t.center, t.Size(), t.left.Size(), t.right.Size())
}

// findCenter attempts to find a midinterval that cuts the intervals
// into half before the interval, and half after the interval.
func findCenter(ints []Interval) int64 {
	// Tukey's ``Ninther,'' median of three medians of three.
	h, m := len(ints), len(ints)/2
	if h > 40 {
		s := h / 8
		sortThree(ints, 0, s, 2*s)
		sortThree(ints, m, m-s, m+s)
		sortThree(ints, h-1, h-1-s, h-1-2*s)
	}

	sortThree(ints, 0, m, h-1)
	return ints[0].Start
}

// sortThree sorts the three provided indicies.
func sortThree(ints []Interval, m1, m0, m2 int) {
	_, _, _ = ints[m1], ints[m0], ints[m2]

	if ints[m1].Start < ints[m0].Start {
		ints[m0], ints[m1] = ints[m1], ints[m0]
	}
	if ints[m2].Start < ints[m1].Start {
		ints[m2], ints[m1] = ints[m1], ints[m2]
		if ints[m1].Start < ints[m0].Start {
			ints[m1], ints[m0] = ints[m0], ints[m1]
		}
	}
}

// leftOf returns true if the interval is entirely to the left of center.
func leftOf(in Interval, center int64) bool {
	return in.End <= center
}

// rightOf returns true if the interval is entirely to the right of center.
func rightOf(in Interval, center int64) bool {
	return in.Start > center
}

// centeredOn returns true if the interval contains the center.
func centeredOn(in Interval, center int64) bool {
	return !leftOf(in, center) && !rightOf(in, center)
}

// split sorts and partitions the intervals into a group that's left of, centered on, and
// right of the given center.
func split(ints []Interval, center int64) (left, right int) {
	right = len(ints) - 1

	for i := 0; i < len(ints) && i < right; i++ {
		in := ints[i]

		if leftOf(in, center) {
			for left < i && leftOf(ints[left], center) {
				left++
			}
			if left != i {
				ints[left], ints[i] = ints[i], ints[left]
				i--
			}
		} else if rightOf(in, center) {
			for right > i && rightOf(ints[right], center) {
				right--
			}
			if right != i {
				ints[right], ints[i] = ints[i], ints[right]
				i--
			}
		}
	}

	for left < len(ints) && leftOf(ints[left], center) {
		left++
	}
	for right >= 0 && rightOf(ints[right], center) {
		right--
	}

	return left, right + 1
}

// copyIntervals does what it says on the tin.
func copyIntervals(ints []Interval) []Interval {
	return append([]Interval(nil), ints...)
}

// sortByStart does what it says on the tin on a copy of the intervals.
func sortByStart(ints []Interval) []Interval {
	sorted := copyIntervals(ints)
	sortIntervalsStart(sorted)
	return sorted
}

// sortByEnd does what it says on the tin on a copy of the intervals.
func sortByEnd(ints []Interval) []Interval {
	sorted := copyIntervals(ints)
	sortIntervalsEnd(sorted)
	return sorted
}

// Find calls the callback with every interval that contains the center.
func (t *Tree) Find(center int64, cb func(Interval)) {
	if t.left == nil || t.right == nil { // "leaf" node
		for _, in := range t.starts {
			if centeredOn(in, center) {
				cb(in)
			}
		}
		return
	}

	switch {
	case center == t.center:
		for _, in := range t.starts {
			cb(in)
		}
		return

	case center < t.center:
		for _, in := range t.starts {
			if in.Start > center {
				break
			}
			cb(in)
		}
		t.left.Find(center, cb)

	case center > t.center:
		for _, in := range t.ends {
			if in.End <= center {
				break
			}
			cb(in)
		}
		t.right.Find(center, cb)
	}
}
