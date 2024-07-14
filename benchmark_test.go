package robin

import (
	"math/rand"
	"testing"
	"time"
)

type testS struct {
	a int
	b string
	c string
	d float64
	e bool
}

func BenchmarkWRRNext(b *testing.B) {
	b.N = 10000000
	b.ReportAllocs()
	rand.Seed(time.Now().UnixNano())
	w := &WRR[testS]{}
	for i := 0; i < 1000; i++ {
		w.Add(testS{}, rand.Intn(10)+100)
	}

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w.Next()
		}
	})
}
