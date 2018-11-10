package interval

import "sort"

type Tree struct {
	left  *Tree // all intervals to the left
	right *Tree // all intervals to the right

	center int64
	starts []Point
	ends   []Point
}

type Point struct {
	Start int64
	End   int64
	Key   int64
}

func New(points []Point) *Tree {
	// We expect ~5 points to fit in a cache line because they are
	// 24 bytes, assuming a cache line is 128 bytes. So, we try to
	// bisect until we get to 5 points

	center := findCenter(points)
	if len(points) <= 5 {
		return &Tree{
			center: center,
			starts: sortByStart(points),
			ends:   sortByEnd(points),
		}
	}

	left, overlap, right := split(points, center)
	return &Tree{
		left:  New(left),
		right: New(right),

		center: center,
		starts: sortByStart(overlap),
		ends:   sortByEnd(overlap),
	}
}

// findCenter attempts to find a midpoint that cuts the points
// into half before the point, and half after the point.
func findCenter(points []Point) int64 {
	h, m := len(points), len(points)/2
	if h > 40 { // Tukey's ``Ninther,'' median of three medians of three.
		s := h / 8
		medianOfThree(points, 0, s, 2*s)
		medianOfThree(points, m, m-s, m+s)
		medianOfThree(points, h-1, h-1-s, h-1-2*s)
	}
	medianOfThree(points, 0, m, h-1)
	return points[0].Start
}

func medianOfThree(points []Point, m1, m0, m2 int) {
	if points[m1].Start < points[m0].Start {
		points[m0], points[m1] = points[m1], points[m0]
	}
	if points[m2].Start < points[m1].Start {
		points[m2], points[m1] = points[m1], points[m2]
		if points[m1].Start < points[m0].Start {
			points[m1], points[m0] = points[m0], points[m1]
		}
	}
}

func split(points []Point, center int64) (left, overlap, right []Point) {
	return nil, nil, nil
}

func copyPoints(points []Point) []Point {
	return append([]Point(nil), points...)
}

func sortByStart(points []Point) []Point {
	sorted := copyPoints(points)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Start < sorted[j].Start })
	return sorted
}

func sortByEnd(points []Point) []Point {
	sorted := copyPoints(points)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].End < sorted[j].End })
	return sorted
}

func (t *Tree) Find(cb func(key int64), start, end int64) {

}
