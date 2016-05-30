package utils

import (
        "fmt"
        "strings"
        "bufio"
        "os"
        "regexp"
        "io/ioutil"
)

var (
        pat_ver    = regexp.MustCompile(`\d+.*`)
)

func Exist(filename string) bool {
    _, err := os.Stat(filename)
    return err == nil || os.IsExist(err)
}

func ReadRelease() (string, string, error) {
        var (
                ID         string
                VERSION    string
        )
        var err error

        fpath := "/etc/redhat-release"
        if Exist(fpath) {
                err = readLine(fpath, func(line string) error {
                lines := strings.SplitN(line, " ", 2)
                if len(lines) == 2 {
                        ID = strings.ToLower(lines[0])
                        VERSION = pat_ver.FindString(lines[1])
                        return nil
                } else {
                        return fmt.Errorf("invalid file format: %v, %d", lines, len(lines))
                }
                })
        }
        fpath = "/etc/SuSE-release"
        if Exist(fpath) {
                ID,VERSION:=getOs(fpath)
                return ID,VERSION,nil
        }

        return ID, VERSION, err
}

func readLine(fname string, line func(string) error) error {
        f, err := os.Open(fname[:])
        if err != nil {
                return err
        }
        defer f.Close()
        scanner := bufio.NewScanner(f)
        for scanner.Scan() {
                if err := line(scanner.Text()); err != nil {
                        return err
                }
        }
        return scanner.Err()
}

func getOs(fpath string) (string,string) {
    data, err := ioutil.ReadFile(fpath)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    lines := strings.Split(string(data), "\n")
    ID:=strings.Split(lines[0]," ")[0]
    VERSION:=strings.Replace(strings.Split(lines[1],"=")[1]," ","",-1)+"SP"+strings.Replace(strings.Split(lines[2],"=")[1]," ","",-1)
    return ID,VERSION
}
