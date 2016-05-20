package main

import (
	"flag"
	"fmt"
	"gitcafe.com/ops/updater/cron"
	"gitcafe.com/ops/updater/g"
	//"gitcafe.com/ops/updater/http"
	//"github.com/toolkits/sys"
	//"gitcafe.com/ops/common/utils"
	"log"
	"os"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if err := g.ParseConfig(*cfg); err != nil {
		log.Fatalln(err)
	}

	g.InitGlobalVariables()

	//go http.Start()
	go cron.Heartbeat()
	

	select {}
}
