package utils

import "sync"

func ParallelForEach[T any](items []T, concurrencyLimit int, fn func(itemIndex int, item T)) {
	sem := make(chan struct{}, concurrencyLimit)
	var wg sync.WaitGroup

	for i, item := range items {
		wg.Add(1)
		sem <- struct{}{}

		go func(itemIndex int, val T) {
			defer wg.Done()
			defer func() { <-sem }()
			fn(itemIndex, val)
		}(i, item)
	}

	wg.Wait()
}
