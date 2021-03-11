package main

import (
	"log"
	"fmt"
	"os"
)

func ShowCPUPercent(jun *JuniperUtilizationReader) {
	val, err := jun.GetCPUPercent()
	if err != nil {
		log.Printf("Failed to get CPU percentage: %s", err)
	}
	log.Printf("CPU usage: %.2f percent", val)
}

func ShowMemPercent(jun *JuniperUtilizationReader) {
	val, err := jun.GetMemoryPercent()
	if err != nil {
		log.Printf("Failed to get mem percentage: %s", err)
	}
	log.Printf("Memory usage: %.2f percent", val)
}

func ShowCPU(jun *JuniperUtilizationReader) {
	val, err := jun.GetCPUTime()
	if err != nil {
		log.Printf("Failed to get CPU time: %s", err)
	}
	log.Printf("CPU time usage: %.0f ", val)
}

func ShowMem(jun *JuniperUtilizationReader) {
	val, err := jun.GetMemoryMB()
	if err != nil {
		log.Printf("Failed to get memory: %s", err)
	}
	log.Printf("Memory usage: %.0f MB", val)
}

func main() {
	jun, err := NewJuniperUtilizationReader()
	if err != nil {
		log.Printf("Failed to create juniper utilization reader: %s", err)
	}

	log.Printf("Creating 10000 files and writing short string to them")
	numFiles := 10000
	for i := 0; i < numFiles; i++ {
		path := fmt.Sprintf("test-file-%v.txt", i)
		var file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Printf("Failed to open file %s: %v", path, err)
			break
		}
		defer file.Close()
		// About 20 bytes, so 10^5 of these is about .2MB
		_, err = file.WriteString(fmt.Sprintf("hello this is file %v", i))
		if err != nil {
			log.Printf("Failed to write to file %s: %v", path, err)
		}
	}

	ShowCPUPercent(jun)
	ShowMemPercent(jun)
	ShowCPU(jun)
	ShowMem(jun)

	log.Printf("Cleaning up 100 files")
	for i := 0; i < numFiles; i++ {
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

