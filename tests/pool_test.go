package tests

import (
	"fmt"
	"testing"

	"github.com/spinvettle/gcp"
)

func BenchmarkPool(b *testing.B) {
	pool, _ := gcp.New(100)
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
func TestPool(b *testing.T) {
	pool, _ := gcp.New(2)
	defer pool.ShutDown()
	// var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		// wg.Add(1)
		for {
			err := pool.Submit(func() {
				for i := 0; i < 10000; i++ {
					_ = i * i
				}
				fmt.Println(i)
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
