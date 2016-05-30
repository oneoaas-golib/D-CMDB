package utils

import (
	"log"	
	"os/exec"
	"strings"
	//"strconv"
)


//Cpunum
func Cpunum()(string,error) {
        lsCmd := exec.Command("bash", "-c", "lscpu | grep ^CPU\\(s\\):")
        lsOut, err := lsCmd.Output()
        if err != nil {
                log.Println("ERROR: Cpunum() fail", err)
       }
        value:=strings.Trim(string(lsOut),"\n")
        cpunum_str:=strings.Split(value,":")[1]
        cpunum:=strings.TrimSpace(cpunum_str)
        return cpunum,nil
}

