package main

import (
	"log"
	"fmt"
	"os"
	"golang.org/x/sys/unix"
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

func ShowUsages(jun *JuniperUtilizationReader) {
	fmt.Println()
	ShowCPUPercent(jun)
	ShowMemPercent(jun)
	ShowCPU(jun)
	ShowMem(jun)
	fmt.Println()
}

func main() {
	jun, err := NewJuniperUtilizationReader()
	if err != nil {
		log.Printf("Failed to create juniper utilization reader: %s", err)
	}

	log.Printf("Starting up. Taking baseline stats")
	ShowUsages(jun)

	numFiles := 60
	log.Printf("Creating %v files and writing short string to them", numFiles)
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

	ShowUsages(jun)

	log.Printf("Cleaning up %v files", numFiles)
	for i := 0; i < numFiles; i++ {
		path := fmt.Sprintf("test-file-%v.txt", i)
		err := os.Remove(path)
		if err != nil {
			log.Printf("Failed to delete file %s: %v", path, err)
			break
		}
	}

	
	var rLimit unix.Rlimit
	err = unix.Getrlimit(unix.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Printf("Couldn't get rLimit: %v", err)
		return
	}
	log.Printf("Current limit on open files: %v", rLimit.Cur)
	log.Printf("Max limit on open files: %v", rLimit.Max)

	newFileLimit := int64(2 * numFiles - 5)
	log.Printf("Limiting number of open files to %v", newFileLimit)
	rLimit.Cur = newFileLimit
	err = unix.Setrlimit(unix.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Printf("Couldn't set rLimit: %v", err)
		return
	}

	log.Printf("Creating %v files and writing short string to them", numFiles)
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

	ShowUsages(jun)

	log.Printf("Cleaning up %v files", numFiles)
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

