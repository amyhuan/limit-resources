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

func() TestOpenFileLimit() {
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
}

func TestMemoryLimit() {
	ShowUsages(jun)

	bytesToWrite := 1000000
	fileName := fmt.Sprintf("%vbyte-file.txt", bytesToWrite)

	// Write big file
	log.Printf("Writing %v bytes to file %v", bytesToWrite, fileName)
	var file, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Failed to open file %s: %v", fileName, err)
		return
	}
	defer file.Close()

	bytes := make([]byte, bytesToWrite)
    rand.Read(bytes)
	_, err = file.Write(bytes)
	if err != nil {
		log.Printf("Failed to write to file %s: %v", fileName, err)
	}

	ShowUsages(jun)

	log.Printf("Cleaning up file %v", fileName)
	err := os.Remove(fileName)
	if err != nil {
		log.Printf("Failed to delete file %s: %v", fileName, err)
		return
	}

	// Change memory limit
	var rLimit unix.Rlimit
	err = unix.Getrlimit(unix.RLIMIT_RSS, &rLimit)
	if err != nil {
		log.Printf("Couldn't get rLimit: %v", err)
		return
	}
	log.Printf("Current limit on memory: %v bytes", rLimit.Cur)
	log.Printf("Max limit on memory: %v bytes", rLimit.Max)

	newMemLimit := 500000
	log.Printf("Limiting memory to %v", newMemLimit)
	rLimit.Cur = newMemLimit
	err = unix.Setrlimit(unix.RLIMIT_RSS, &rLimit)
	if err != nil {
		log.Printf("Couldn't set rLimit: %v", err)
		return
	}

	// Write big file again, should fail
	log.Printf("Writing %v bytes to file %v", bytesToWrite, fileName)
	var file, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Failed to open file %s: %v", fileName, err)
		return
	}
	defer file.Close()

	bytes := make([]byte, bytesToWrite)
    rand.Read(bytes)
	_, err = file.Write(bytes)
	if err != nil {
		log.Printf("Failed to write to file %s: %v", fileName, err)
	}

	ShowUsages(jun)

	log.Printf("Cleaning up file %v", fileName)
	err := os.Remove(fileName)
	if err != nil {
		log.Printf("Failed to delete file %s: %v", fileName, err)
		return
	}
}

func TestMemoryLimit() {
	ShowUsages(jun)

	bytesToWrite := 1000000
	fileName := fmt.Sprintf("%vbyte-file.txt", bytesToWrite)

	// Write big file
	log.Printf("Writing %v bytes to file %v", bytesToWrite, fileName)
	var file, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Failed to open file %s: %v", fileName, err)
		return
	}
	defer file.Close()

	bytes := make([]byte, bytesToWrite)
    rand.Read(bytes)
	_, err = file.Write(bytes)
	if err != nil {
		log.Printf("Failed to write to file %s: %v", fileName, err)
	}

	ShowUsages(jun)

	log.Printf("Cleaning up file %v", fileName)
	err := os.Remove(fileName)
	if err != nil {
		log.Printf("Failed to delete file %s: %v", fileName, err)
		return
	}

	// Change memory limit
	var rLimit unix.Rlimit
	err = unix.Getrlimit(unix.RLIMIT_RSS, &rLimit)
	if err != nil {
		log.Printf("Couldn't get rLimit: %v", err)
		return
	}
	log.Printf("Current limit on memory: %v bytes", rLimit.Cur)
	log.Printf("Max limit on memory: %v bytes", rLimit.Max)

	newMemLimit := 500000
	log.Printf("Limiting memory to %v", newMemLimit)
	rLimit.Cur = newMemLimit
	err = unix.Setrlimit(unix.RLIMIT_RSS, &rLimit)
	if err != nil {
		log.Printf("Couldn't set rLimit: %v", err)
		return
	}

	// Write big file again, should fail
	log.Printf("Writing %v bytes to file %v", bytesToWrite, fileName)
	var file, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Failed to open file %s: %v", fileName, err)
		return
	}
	defer file.Close()

	bytes := make([]byte, bytesToWrite)
    rand.Read(bytes)
	_, err = file.Write(bytes)
	if err != nil {
		log.Printf("Failed to write to file %s: %v", fileName, err)
	}

	ShowUsages(jun)

	log.Printf("Cleaning up file %v", fileName)
	err := os.Remove(fileName)
	if err != nil {
		log.Printf("Failed to delete file %s: %v", fileName, err)
		return
	}
}

func TestCPULimit() {
	ShowUsages(jun)

	// Do CPU intensive task
	log.Printf("Doing CPU intensive task")
	for i := 0; i < 2000000000; i++ {
		x = x + 1
		x = x - 1
		if i == 1000000000 {
			ShowUsages(jun)
		}
	}

	// Change CPU limit
	var rLimit unix.Rlimit
	err = unix.Getrlimit(unix.RLIMIT_CPU, &rLimit)
	if err != nil {
		log.Printf("Couldn't get rLimit: %v", err)
		return
	}
	log.Printf("Current limit on CPU time: %v bytes", rLimit.Cur)
	log.Printf("Max limit on CPU time: %v bytes", rLimit.Max)

	newCPULimit := 1
	log.Printf("Limiting CPU time to %v", newCPULimit)
	rLimit.Cur = newCPULimit
	err = unix.Setrlimit(unix.RLIMIT_CPU, &rLimit)
	if err != nil {
		log.Printf("Couldn't set rLimit: %v", err)
		return
	}

	// Do CPU intensive task again
	log.Printf("Doing CPU intensive task")
	for i := 0; i < 2000000000; i++ {
		x = x + 1
		x = x - 1
		if i == 1000000000 {
			ShowUsages(jun)
		}
	}
}


func main() {
	jun, err := NewJuniperUtilizationReader()
	if err != nil {
		log.Printf("Failed to create juniper utilization reader: %s", err)
	}

	//TestOpenFileLimit()
	TestMemoryLimit()
	TestCPULimit()
}

