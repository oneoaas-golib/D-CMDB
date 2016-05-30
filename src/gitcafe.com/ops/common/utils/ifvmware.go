package utils

import (
	"log"	
	"os/exec"
	"strings"
	"strconv"
)


//Vminfo
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
               return "vmware",nil 
        }
        return "",nil
}

