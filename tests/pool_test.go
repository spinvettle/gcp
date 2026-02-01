package tests

import (
	"testing"

	"github.com/spinvettle/gcp"
)

func BenchmarkPool(b *testing.B) {
	pool := gcp.New(100)
	defer pool.ShutDown()
	// var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		// wg.Add(1)
		for {
			err := pool.Submit(func() {
				_ = 1 + 1
			})
			if err == nil {
				break
				// b.Fatal(err.Error())

				// wg.Done()
			}
		}

	}
	pool.ShutDown()
	// wg.Wait()
}
