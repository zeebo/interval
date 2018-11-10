package interval

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/zeebo/mon"
)

func init() { rand.Seed(time.Now().UnixNano()) }

func TestCenter(t *testing.T) {
	var (
		counts  mon.Histogram
		centers mon.Histogram
	)

	for i := 0; i < 10000; i++ {
		points := make([]Point, 100)
		for i := range points {
			points[i].Start = int64(rand.Intn(100))
		}

		center := findCenter(points)
		count := int64(0)
		for _, pt := range points {
			if pt.Start < center {
				count++
			}
		}

		centers.Observe(center)
		counts.Observe(count)
	}

	fmt.Println(centers.Average(), centers.Quantile(.1), centers.Quantile(.9))
	fmt.Println(counts.Average(), counts.Quantile(.1), counts.Quantile(.9))
}
