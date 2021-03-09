package main

import (
	"log"
	"fmt"
	"os"
)

func ShowCPU(jun *JuniperUtilizationReader) {
	cpuPercent, err := jun.GetCPUPercent()
	if err != nil {
		log.Printf("Failed to get CPU percentage: %s", err)
	}
	log.Printf("Baseline CPU usage: %d\%", cpuPercent)
}

func ShowMem(jun *JuniperUtilizationReader) {
	memPercent, err := jun.GetMemoryPercent()
	if err != nil {
		log.Printf("Failed to get CPU percentage: %s", err)
	}
	log.Printf("Baseline memory usage: %d\%", memPercent)
}

func main() {
	jun, err := NewJuniperUtilizationReader()
	if err != nil {
		log.Printf("Failed to create juniper utilization reader: %s", err)
	}

	for i := 0; i < 100; i++ {
		path := fmt.Sprintf("test-file-%v.txt", i)
		var file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Printf("Failed to open file %s: %v", path, err)
			break
		}
		defer file.Close()
		_, err = file.WriteString(fmt.Sprintf("hello this is file %v", i))
		if err != nil {
			log.Printf("Failed to write to file %s: %v", path, err)
		}
	}

	ShowCPU(jun)
	ShowMem(jun)

	// Delete files made
	for i := 0; i < 100; i++ {
		path := fmt.Sprintf("test-file-%v.txt", i)
		err := os.Remove(path)
		if err != nil {
			log.Printf("Failed to delete file %s: %v", path, err)
			break
		}
	}
	

	// Use a lot of CPU
	// Open large # of files
	// Write a lot to open files
	// Allocate a lot of memory
}

