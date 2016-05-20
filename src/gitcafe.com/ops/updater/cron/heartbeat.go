package cron

import (
	"encoding/json"
	"fmt"
	"gitcafe.com/ops/common/model"
	"gitcafe.com/ops/common/utils"
	"gitcafe.com/ops/updater/g"
	"github.com/toolkits/net/httplib"
	"log"
	"time"
)

func Heartbeat() {
	SleepRandomDuration()
	for {
		//heartbeat()
		cmdbhb()
		d := time.Duration(g.Config().Interval) * time.Second
		time.Sleep(d)
	}
}

func cmdbhb() {
	hostname, err := utils.Hostname(g.Config().Hostname)
        if err != nil {
                return
        }	

	vminfo,err:=utils.Ifvmware()
	if err != nil {
                return
        }
        cpuinfo,err:=utils.Cpuinfo()
	if err != nil {
                return
        }
	
	cmdbhbRequest := BuildcmdbhbRequest(hostname, vminfo,cpuinfo)
        if g.Config().Debug {
                log.Println("====>>>>")
                log.Println(cmdbhbRequest)
        }
	
	bs, err := json.Marshal(cmdbhbRequest)
        if err != nil {
                log.Println("encode cmdbhb request fail", err)
                return
        }
	
	url := fmt.Sprintf("http://%s/cmdbhb", g.Config().Server)
        httpRequest := httplib.Post(url).SetTimeout(time.Second*10, time.Minute)
        httpRequest.Body(bs)
        httpResponse, err := httpRequest.Bytes()
        if err != nil {
                log.Printf("curl %s fail %v", url, err)
                return
        }
	var cmdbhbResponse model.CmdbhbResponse
        err = json.Unmarshal(httpResponse, &cmdbhbResponse)
        if err != nil {
                log.Println("decode cmdbhb response fail", err)
                return
        }

        if g.Config().Debug {
                log.Println("<<<<====")
                log.Println(cmdbhbResponse)
        }
	
}

func heartbeat() {
	agentDirs, err := ListAgentDirs()
	if err != nil {
		return
	}

	hostname, err := utils.Hostname(g.Config().Hostname)
	if err != nil {
		return
	}

	heartbeatRequest := BuildHeartbeatRequest(hostname, agentDirs)
	if g.Config().Debug {
		log.Println("====>>>>")
		log.Println(heartbeatRequest)
	}

	bs, err := json.Marshal(heartbeatRequest)
	if err != nil {
		log.Println("encode heartbeat request fail", err)
		return
	}

	url := fmt.Sprintf("http://%s/heartbeat", g.Config().Server)
	httpRequest := httplib.Post(url).SetTimeout(time.Second*10, time.Minute)
	httpRequest.Body(bs)
	httpResponse, err := httpRequest.Bytes()
	if err != nil {
		log.Printf("curl %s fail %v", url, err)
		return
	}

	var heartbeatResponse model.HeartbeatResponse
	err = json.Unmarshal(httpResponse, &heartbeatResponse)
	if err != nil {
		log.Println("decode heartbeat response fail", err)
		return
	}

	if g.Config().Debug {
		log.Println("<<<<====")
		log.Println(heartbeatResponse)
	}

	HandleHeartbeatResponse(&heartbeatResponse)

}
