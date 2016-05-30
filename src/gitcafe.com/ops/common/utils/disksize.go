package utils

import (
    "syscall"
    "fmt"
)



func Disksize()(disksize string,err error){
	var stat syscall.Statfs_t
	wd:="/"
	syscall.Statfs(wd, &stat)
 	disksize=fmt.Sprintf("%.2f GB", float64(stat.Blocks * uint64(stat.Bsize))/float64(GB))
	return disksize,nil
}
