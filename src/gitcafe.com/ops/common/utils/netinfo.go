package utils

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"net"
	"io/ioutil"
	"os/exec"
	"bufio"
	"bytes"
	"io"
)

type Dev struct {
	Dev_name string
	Dev_type string
	Dev_ip []string
	Dev_gw []string
}


const (
    PROC_TCP = "/proc/net/tcp"
    PROC_UDP = "/proc/net/udp"
    PROC_TCP6 = "/proc/net/tcp6"
    PROC_UDP6 = "/proc/net/udp6"
    PROC_ROUTE = "/proc/net/route"

)


type Rout struct {
    iface string
    dst string
    gateway string
}


func check(err error) {
        if err != nil {
                panic(err)
        }

}


func getData(t string) []string {
    // Get data from tcp or udp file.

    var proc_t string

    if t == "tcp" {
        proc_t = PROC_TCP
    } else if t == "udp" {
        proc_t = PROC_UDP
    } else if t == "tcp6" {
        proc_t = PROC_TCP6
    } else if t == "udp6" {
        proc_t = PROC_UDP6
    } else if t == "route" {
        proc_t = PROC_ROUTE
    } else {
        fmt.Printf("%s is a invalid type, tcp and udp only!\n", t)
        os.Exit(1)
    }


    data, err := ioutil.ReadFile(proc_t)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    lines := strings.Split(string(data), "\n")

    // Return lines without Header line and blank line on the end
    return lines[1:len(lines) - 1]

}

func removeEmpty(array []string) []string {
    // remove empty data from line
    var new_array [] string
    for _, i := range(array) {
        if i != "" {
           new_array = append(new_array, i)
        }
    }
    return new_array
}


func convertIp(ip string) string {
    // Convert the ipv4 to decimal. Have to rearrange the ip because the
    // default value is in little Endian order.

    var out string

    // Check ip size if greater than 8 is a ipv6 type
    if len(ip) > 8 {
        i := []string{ ip[30:32],
                        ip[28:30],
                        ip[26:28],
                        ip[24:26],
                        ip[22:24],
                        ip[20:22],
                        ip[18:20],
                        ip[16:18],
                        ip[14:16],
                        ip[12:14],
                        ip[10:12],
                        ip[8:10],
                        ip[6:8],
                        ip[4:6],
                        ip[2:4],
                        ip[0:2]}
        out = fmt.Sprintf("%v%v:%v%v:%v%v:%v%v:%v%v:%v%v:%v%v:%v%v",
                            i[14], i[15], i[13], i[12],
                            i[10], i[11], i[8], i[9],
                            i[6],  i[7], i[4], i[5],
                            i[2], i[3], i[0], i[1])

    } else {
        i := []int64{ hexToDec(ip[6:8]),
                       hexToDec(ip[4:6]),
                       hexToDec(ip[2:4]),
                       hexToDec(ip[0:2]) }

       out = fmt.Sprintf("%v.%v.%v.%v", i[0], i[1], i[2], i[3])
    }
   return out
}

func hexToDec(h string) int64 {
    // convert hexadecimal to decimal.
    d, err := strconv.ParseInt(h, 16, 32)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    return d
}



func gw2(t string) []Rout{
    var Routs []Rout
    data:=getData(t)
    for _,line := range(data) {
        line_array := removeEmpty(strings.Split(strings.TrimSpace(line), "\t"))
        iface:=line_array[0]
        dst:=line_array[1]
        gateway:=convertIp(line_array[2])
	if dst=="00000000" {
		dst="default"
        	p:=Rout{iface,dst,gateway}
        	Routs=append(Routs,p)
	}
    }
    return Routs
}



func Route2() []string {
    data:=gw2("route")
    //return data  //data is []Rout
    var s []string
    for _,data_struct:=range data {
	x:=data_struct.iface+" "+data_struct.gateway
	s=append(s,x)
    }
    return s
}



func Getnetinfo() []Dev{
	dev_array:=[]Dev{}
        interfaces, err := net.Interfaces()
        check(err)
	gw_a:=Route2()
        for _, i := range interfaces {
                if strings.Contains(i.Flags.String(), "up") {
                        if strings.Contains(i.Name,"lo") {
                                continue
                        }
                        if strings.Contains(i.Name,"virbr") {
                                continue
                        }
                        if strings.Contains(i.Name,"docker") {
                                continue
                        }
                        fmt.Printf("Name: %v up \n", i.Name)
                        addrs,err:=i.Addrs()
                        check(err)
			var ipstr []string
                        for _,a:=range addrs {
                                switch v:=a.(type){
                                case *net.IPNet:
					ss:=fmt.Sprintf("%s",v)
					if strings.Contains(ss,":") {
						continue
					}
					fmt.Printf("%v : %s  \n",i.Name,v)
					ipstr=append(ipstr,ss)
                                }
                        }
		devx:=Dev{}
		devx.Dev_name=i.Name
		devx.Dev_type=intftype(i.Name)
		devx.Dev_ip=ipstr
		devx.Dev_gw=gw_a
		dev_array=append(dev_array,devx)	
                fmt.Printf("\n")
                }
        }
	return dev_array
}



func intftype(dev string)string {
        ifvm,_:=Ifvmware()
        var line string
        lsCmd := exec.Command("bash", "-c", "ethtool "+dev)
        lsOut, err := lsCmd.Output()
        if err != nil {
                //panic(err)
		return "unknown"
        }

        r := bytes.NewReader(lsOut)
        r1 := bufio.NewReader(r)
        for {
                data, err := r1.ReadSlice('\n')
                if err == io.EOF {
                        return line+" "+ifvm
                        break
                }
                if strings.Contains(string(data), "Speed:") || strings.Contains(string(data), "Port:") {
                        inputstring := strings.Trim(string(data), "\n")
                        inputstring1 := strings.Replace(inputstring,"Speed:","",-1)
                        inputstring2 := strings.Replace(inputstring1,"Port:","",-1)
                        inputstring3 := strings.TrimSpace(string(inputstring2))
                        line = line + " "+inputstring3
                }
        }
        return line+" "+ifvm
}
