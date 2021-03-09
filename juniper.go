package main

import(
	"os"
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

type JuniperUtilizationReader struct {
	pid int
}

func NewJuniperUtilizationReader() (*JuniperUtilizationReader, error) {
	return &JuniperUtilizationReader{
		pid: os.Getpid(),
	}, nil
}

func getCPUPercentPid(pid int)(float64, error){
	sysInfo, err := stat(pid, "ps")
	return float64(sysInfo.CPU), err
}

func getMemoryStatsPid(pid int)(float64, error){
	sysInfo, err := stat(pid, "ps")
	return float64(sysInfo.Memory), err
}

func (jun *JuniperUtilizationReader) GetCPUPercent() (float64, error) {
	return getCPUPercentPid(jun.pid)
}

func (jun *JuniperUtilizationReader) GetMemoryPercent() (float64, error) {
	return getMemoryStatsPid(jun.pid)
}

func stat(pid int, statType string) (*SysInfo, error) {
	sysInfo := &SysInfo{}
	args := "-o pcpu,pmem -p"
	stdout, _ := exec.Command("ps", args, strconv.Itoa(pid)).Output()
	if len(stdout) == 0{
		return sysInfo, errors.New("Didn't get ps printout successfully with pid " + strconv.Itoa(pid))
	}
	ret := formatStdOut(stdout, 1)
	if len(ret) == 0{
		return sysInfo, errors.New("Can't find process with this PID: " + strconv.Itoa(pid))
	}
	sysInfo.CPU = parseFloat(ret[0])
	sysInfo.Memory = parseFloat(ret[1])
	return sysInfo, nil
}

func formatStdOut(stdout []byte, userfulIndex int) []string{
	infoArr := strings.Split(string(stdout), "\n")[userfulIndex]
	ret := strings.Fields(infoArr)
	return ret
}

func parseFloat(val string) float64 {
	floatVal, _ := strconv.ParseFloat(val, 64)
	return floatVal
}

type SysInfo struct {
	CPU	float64
	Memory  float64
}

type Stat struct {
	utime  float64
	stime  float64
	cutime float64
	cstime float64
	start  float64
	rss    float64
	uptime float64
}