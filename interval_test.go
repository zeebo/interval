package interval

import (
	"math/rand"
	"testing"
	"time"

	"github.com/zeebo/mon"
)

func init() { rand.Seed(time.Now().UnixNano()) }

func sortInts(start, end int64) (int64, int64) {
	if start <= end {
		return start, end
	}
	return end, start
}

func randomInterval(max, offset int) Interval {
	start, end := sortInts(int64(rand.Intn(max)+offset), int64(rand.Intn(max)+offset))
	return Interval{Start: start, End: end}
}

func TestCenter(t *testing.T) {
	var (
		counts  mon.Histogram
		centers mon.Histogram
	)

	for i := 0; i < 10000; i++ {
		ints := make([]Interval, 100)
		for i := range ints {
			ints[i].Start = int64(rand.Intn(100))
		}

		center := findCenter(ints)
		count := int64(0)
		for _, pt := range ints {
			if pt.Start < center {
				count++
			}
		}

		centers.Observe(center)
		counts.Observe(count)
	}

	if avg := centers.Average(); avg < 48 || avg > 52 {
		t.Fatalf("bad average: %v", avg)
	}
	if avg := counts.Average(); avg < 48 && avg > 52 {
		t.Fatalf("bad average: %v", avg)
	}
}

func TestFind(t *testing.T) {
	ints := make([]Interval, 100000)
	for i := range ints {
		ints[i] = randomInterval(1000, 0)
	}

	tree := New(ints)
	tree.Find(500, func(i Interval) { t.Logf("%+v\n", i) })
}

func TestSplit(t *testing.T) {
	const center = 500

	for i := 0; i < 1000; i++ {
		max, offset := 2*center, 0
		if i == 0 { // try all to the left
			max, offset = 1, 0
		}
		if i == 1 { // try all to the right
			max, offset = 1, 1000
		}

		ints := make([]Interval, 100)
		class := make(map[Interval]int)

		for i := range ints {
			ints[i] = randomInterval(max, offset)
			if leftOf(ints[i], center) {
				class[ints[i]] = -1
			} else if rightOf(ints[i], center) {
				class[ints[i]] = 1
			} else {
				class[ints[i]] = 0
			}
		}

		rand.Shuffle(len(ints), func(i, j int) { ints[i], ints[j] = ints[j], ints[i] })
		left, right := split(ints, center)

		for _, in := range ints[:left] {
			if val, ok := class[in]; !ok || val != -1 {
				t.Fatal("bad split:", ints, left, right)
			}
		}

		for _, in := range ints[left:right] {
			if val, ok := class[in]; !ok || val != 0 {
				t.Fatal("bad split:", ints, left, right)
			}
		}

		for _, in := range ints[right:] {
			if val, ok := class[in]; !ok || val != 1 {
				t.Fatal("bad split:", ints, left, right)
			}
		}
	}
}

func BenchmarkTree(b *testing.B) {
	run := func(b *testing.B, size int64) {
		ints := make([]Interval, size)
		for n := range ints {
			start := int64(rand.Intn(int(10 * size)))
			ints[n] = Interval{
				Start: start,
				End:   start + int64(rand.Intn(int(size/10+1))),
			}
		}

		starts := sortByStart(ints)
		tree := New(ints)
		cb := func(Interval) {}
		b.Log(tree)

		run := func(b *testing.B, center int64) {
			b.Run("New", func(b *testing.B) {
				b.ReportAllocs()

				for i := 0; i < b.N; i++ {
					_ = New(ints)
				}
			})

			b.Run("Find", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					tree.Find(center, cb)
				}
			})

			b.Run("Linear", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					for _, in := range ints {
						if centeredOn(in, center) {
							cb(in)
						}
					}
				}
			})

			b.Run("Linear+Break", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					for _, in := range starts {
						if centeredOn(in, center) {
							cb(in)
						}
						if rightOf(in, center) {
							break
						}
					}
				}
			})
		}

		b.Run("Left", func(b *testing.B) { run(b, 0) })
		b.Run("Center", func(b *testing.B) { run(b, 10*size/2) })
		b.Run("Right", func(b *testing.B) { run(b, 10*size) })
	}

	b.Run("10", func(b *testing.B) { run(b, 10) })
	b.Run("100", func(b *testing.B) { run(b, 100) })
	b.Run("1000", func(b *testing.B) { run(b, 1000) })
	b.Run("10000", func(b *testing.B) { run(b, 10000) })
	b.Run("100000", func(b *testing.B) { run(b, 100000) })
	b.Run("1000000", func(b *testing.B) { run(b, 1000000) })
}
