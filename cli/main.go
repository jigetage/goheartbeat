package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/jigetage/goheartbeat/cli/heartbeatcli"
)

func main()  {
	cli := heartbeatcli.NewHeartBeatCli(":8999", 2 * time.Second, "cli", "")
	go cli.Run()

	// press ctrl + c to quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("quit (%v)\n", <- sig)
}