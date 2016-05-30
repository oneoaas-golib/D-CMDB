package utils

import (
    "syscall"
    "fmt"
)

const (
        B  = 1
        KB = 1024 * B
        MB = 1024 * KB
        GB = 1024 * MB
)

func Memsize()(memsize string,err error){
    sysInfo := new(syscall.Sysinfo_t)
    err = syscall.Sysinfo(sysInfo)
    if err != nil {
                return
        }
    memsize=fmt.Sprintf("%.2f GB", float64(sysInfo.Totalram)/float64(GB))
    return memsize,nil
}