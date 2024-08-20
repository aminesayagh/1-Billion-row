package tracker

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
	"onBillion/config"
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
		"Alloc Memory":      fmt.Sprintf("%.2f MB", allocMemory),
		"Total Alloc Memory": fmt.Sprintf("%.2f MB", totalAllocMemory),
		"Sys Memory":        fmt.Sprintf("%.2f MB", sysMemory),
		"Heap Memory":       fmt.Sprintf("%.2f MB", heapMemory),
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
	tempFilePath := metricsOutputFilePath + ".tmp"

	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		fmt.Println("Error creating temporary metrics file: ", err)
		return
	}
	defer tempFile.Close()

	// save the metrics to the file in the format: date-version-key=value
	date := time.Now().Format("2006-01-02")
	for key, value := range metrics {
		if hasMetric, _ := searchMetric(metricsOutputFilePath, version, key); hasMetric != 0 {
			// remove the line from original file
			removeMetric(metricsOutputFilePath, hasMetric)
		}
		_, _ = tempFile.WriteString(fmt.Sprintf("%s-%s-%s=%s\n", date, version, key, value))
	}

	// Rename temp file to actual file
	if err := os.Rename(tempFilePath, metricsOutputFilePath); err != nil {
		fmt.Println("Error renaming temporary metrics file: ", err)
	} else {
		fmt.Println("Metrics saved to file: ", metricsOutputFilePath)
	}
}

func searchMetric(metricsOutputFilePath, version, key string) (int, string) {
	metricsOutputFile, err := os.Open(metricsOutputFilePath)
	if err != nil {
		fmt.Println("Error opening metrics file: ", err)
		return 0, ""
	}
	defer metricsOutputFile.Close()

	scanner := bufio.NewScanner(metricsOutputFile)
	currentLine := 0

	for scanner.Scan() {
		currentLine++
		line := scanner.Text()
		parts := strings.Split(line, "-")
		if len(parts) < 3 {
			continue
		}
		if parts[1] == version && strings.Contains(parts[2], key) {
			return currentLine, line
		}
	}
	return 0, ""
}

func removeMetric(metricsOutputFilePath string, line int) {
	metricsOutputFile, err := os.Open(metricsOutputFilePath)
	if err != nil {
		fmt.Println("Error opening metrics file: ", err)
		return
	}
	defer metricsOutputFile.Close()

	tempFilePath := metricsOutputFilePath + ".tmp"
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		fmt.Println("Error creating temporary file: ", err)
		return
	}
	defer tempFile.Close()

	scanner := bufio.NewScanner(metricsOutputFile)
	currentLine := 0

	for scanner.Scan() {
		currentLine++
		if currentLine == line {
			continue
		}
		_, _ = tempFile.WriteString(scanner.Text() + "\n")
	}

	if err := os.Rename(tempFilePath, metricsOutputFilePath); err != nil {
		fmt.Println("Error renaming temporary file: ", err)
	}
}
