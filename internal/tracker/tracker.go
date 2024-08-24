package tracker

import (
	"fmt"
	"onBillion/config"
	"os"
	"runtime"
	"time"
	"sync"
)

type MetricsMap map[string]string

func timer(f func(), metrics chan MetricsMap) {
	// start the timer
	start := time.Now()
	f()

	// end the timer
	end := time.Now()
	diff := end.Sub(start)

	metrics <- MetricsMap{
		"Execution Time": diff.String(),
	}
}

func memory(f func(), metrics chan MetricsMap) {
	// Trigger garbage collection to minimize impact of stale allocations
	runtime.GC()

	// Memory statistics before executing the function
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	f()

	// Memory statistics after executing the function
	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	// Calculate memory used by the function
	allocMemory := float64(memAfter.Alloc-memBefore.Alloc) / 1024 / 1024
	totalAllocMemory := float64(memAfter.TotalAlloc-memBefore.TotalAlloc) / 1024 / 1024
	sysMemory := float64(memAfter.Sys-memBefore.Sys) / 1024 / 1024
	heapMemory := float64(memAfter.HeapAlloc-memBefore.HeapAlloc) / 1024 / 1024

	// Send memory metrics to the metrics channel
	metrics <- map[string]string{
		"Alloc Memory":       fmt.Sprintf("%.2f MB", allocMemory),
		"Total Alloc Memory": fmt.Sprintf("%.2f MB", totalAllocMemory),
		"Sys Memory":         fmt.Sprintf("%.2f MB", sysMemory),
		"Heap Memory":        fmt.Sprintf("%.2f MB", heapMemory),
	}
}

func Run(f func()) {
	done := make(chan bool)
	metrics := make(chan MetricsMap)
	aggregatedMetrics := make(MetricsMap)

	conf := config.GetInstance()

	version := conf.Version
	metricsOutputFilePath := conf.MetricsFilePath

	go func() {
		memory(func() {
			timer(f, metrics)
		}, metrics)

		done <- true
	}()

	for {
		select {
		case <-done:
			fmt.Println("Function execution completed.")
			saveMetrics(metricsOutputFilePath, version, aggregatedMetrics)
			return
		case m := <-metrics:
			for k, v := range m {
				fmt.Printf("--- Metrics %s: %s\n", k, v)
				aggregatedMetrics[k] = v
			}

		case <-time.After(10 * time.Minute):
			fmt.Println("Function execution timed out.")
			saveMetrics(metricsOutputFilePath, version, aggregatedMetrics)
			done <- true
			return
		}
	}
}

func saveMetrics(metricsOutputFilePath string, version string, metrics MetricsMap) {
	file, err := os.OpenFile(metricsOutputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Version: %s\n", version))
	if err != nil {
		fmt.Println(err)
		return
	}

	date := time.Now().Format("2006-01-02 15:04:05")

	for k, v := range metrics {
		_, err = file.WriteString(fmt.Sprintf("%s-%s-%s-%s\n", date, version, k, v))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// MeasureExecutionTime measures the average execution time of two functions and compares them
func MeasureExecutionTime(f1, f2 func(), iterations int) (time.Duration, time.Duration) {
	var wg sync.WaitGroup
	wg.Add(2)

	var totalTimeF1, totalTimeF2 time.Duration

	// Run the first function multiple times in a goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			start := time.Now()
			f1()
			totalTimeF1 += time.Since(start)
		}
	}()

	// Run the second function multiple times in a goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			start := time.Now()
			f2()
			totalTimeF2 += time.Since(start)
		}
	}()

	// Wait for both goroutines to finish
	wg.Wait()

	// Calculate average times
	avgTimeF1 := totalTimeF1 / time.Duration(iterations)
	avgTimeF2 := totalTimeF2 / time.Duration(iterations)

	// Print the results
	fmt.Printf("Function 1 average time: %v\n", avgTimeF1)
	fmt.Printf("Function 2 average time: %v\n", avgTimeF2)

	// Determine which function is slower
	if avgTimeF1 > avgTimeF2 {
		fmt.Println("Function 1 is slower.")
	} else if avgTimeF1 < avgTimeF2 {
		fmt.Println("Function 2 is slower.")
	} else {
		fmt.Println("Both functions have the same average execution time.")
	}

	return avgTimeF1, avgTimeF2
}