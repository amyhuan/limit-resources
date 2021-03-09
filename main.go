package main

import (
	"log"
)

func main() {
	jun, err := NewJuniperUtilizationReader()
	if err != nil {
		log.Printf("Failed to create juniper utilization reader: %s", err)
	}

	cpuPercent, err := jun.GetCPUPercent()
	if err != nil {
		log.Printf("Failed to get CPU percentage: %s", err)
	}
	log.Printf("Baseline CPU usage: %f%", cpuPercent)

	memPercent, err := jun.GetMemoryPercent()
	if err != nil {
		log.Printf("Failed to get CPU percentage: %s", err)
	}
	log.Printf("Baseline memory usage: %f%", memPercent)
	// Maybe do like...2 MB
	// Perhaps get measuring stuff first

	// Use a lot of CPU
	// Open large # of files
	// Write a lot to open files
	// Allocate a lot of memory
}

