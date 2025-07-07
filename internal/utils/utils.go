package utils

import (
	"os"
	"os/exec"
	"sync"

	"github.com/spf13/viper"
)

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

func LogCommandIfVerbose(cmd *exec.Cmd) {
	verboseLogs := viper.GetBool("verbose_logs")

	if verboseLogs {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = nil
		cmd.Stderr = nil
	}
}
