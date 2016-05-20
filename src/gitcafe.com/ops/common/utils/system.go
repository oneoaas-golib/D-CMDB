package utils

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"strconv"
	"path/filepath"
	"bufio"
)

func Hostname(configHostname string) (string, error) {
	if configHostname != "" {
		return configHostname, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Println("ERROR: os.Hostname() fail", err)
	}

	return hostname, err
}

func Ifvmware()(string,error) {
        lsCmd := exec.Command("bash", "-c", "lscpu | grep -i vmware | wc -l")
        lsOut, err := lsCmd.Output()
        if err != nil {
		log.Println("ERROR: ifvmware() fail", err)
       }
        //fmt.Println(string(lsOut))
        value:=strings.Trim(string(lsOut),"\n")
        //fmt.Println("value:",value)
        number, _ := strconv.ParseInt(value, 10, 0)
        if number==1 {
               return "virtual machine",nil 
        }
        return "",nil
}

func GetEnv(key string, dfault string, combineWith ...string) string {
	value := os.Getenv(key)
	if value == "" {
		value = dfault
	}

	switch len(combineWith) {
	case 0:
		return value
	case 1:
		return filepath.Join(value, combineWith[0])
	default:
		all := make([]string, len(combineWith)+1)
		all[0] = value
		copy(all[1:], combineWith)
		return filepath.Join(all...)
	}
	log.Println("invalid switch case")
	return ""
}

func HostProc(combineWith ...string) string {
	return GetEnv("HOST_PROC", "/proc", combineWith...)
}

func HostSys(combineWith ...string) string {
	return GetEnv("HOST_SYS", "/sys", combineWith...)
}

func HostEtc(combineWith ...string) string {
	return GetEnv("HOST_ETC", "/etc", combineWith...)
}

func Cpuinfo() (string, error) {
	filename := HostProc("cpuinfo")
	lines, _ := ReadLines(filename)
        var ModelName=""
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		switch key {
		    case "model name":
			ModelName = value
		}
	}
	return ModelName,nil	
}

func ReadLines(filename string) ([]string, error) {
	return ReadLinesOffsetN(filename, 0, -1)
}

// ReadLines reads contents from file and splits them by new line.
// The offset tells at which line number to start.
// The count determines the number of lines to read (starting from offset):
//   n >= 0: at most n lines
//   n < 0: whole file
func ReadLinesOffsetN(filename string, offset uint, n int) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []string{""}, err
	}
	defer f.Close()

	var ret []string

	r := bufio.NewReader(f)
	for i := 0; i < n+int(offset) || n < 0; i++ {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		if i < int(offset) {
			continue
		}
		ret = append(ret, strings.Trim(line, "\n"))
	}

	return ret, nil
}
