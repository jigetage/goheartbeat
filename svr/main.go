package main

import (
	"fmt"
	"github.com/jigetage/goheartbeat/svr/heartbeatsvr"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
	svr := heartbeatsvr.NewHeartBeatSvr(8999)
	go svr.ServerSocket()

	// press ctrl + c to quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("quit (%v)\n", <- sig)
}
