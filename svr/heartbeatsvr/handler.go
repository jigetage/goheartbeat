package heartbeatsvr

import (
	"encoding/json"
	"fmt"
	"github.com/jigetage/goheartbeat/common"
	"log"
	"net"
	"sync"
)

type HeartBeatSvr struct {
	Port int64
}

var (
	HeartBeatServer *HeartBeatSvr
	once sync.Once
)

func NewHeartBeatSvr(port int64) *HeartBeatSvr {
	once.Do(func() {
		HeartBeatServer = &HeartBeatSvr{
			Port: port,
		}
	})
	return HeartBeatServer
}

func connHandler(c net.Conn) {
	defer c.Close()
	if c == nil {
		log.Println("invalid tcp connection")
	}

	buf := make([]byte, common.RecvBuf)
	for {
		log.Println("111111")
		cnt, err := c.Read(buf)
		log.Println("555555")

		if err != nil {
			log.Printf("read failed, error is: %v\n", err)
			break
		}
		log.Println("222222")

		if cnt == 0 {
			c.Write([]byte("oh, no data"))
			log.Println("333333")
			continue
		}

		log.Println("444444")
		datav := common.Msg{}
		err = json.Unmarshal(buf[:cnt], &datav)
		if err != nil {
			log.Printf("parse msg failed, error is: %v\n", err)
			c.Write([]byte("oh, msg illegal"))
			continue
		}

		c.Write([]byte("hello" + datav.CliName))
	}
}

// ServerSocket ServerSocket
func (hbs *HeartBeatSvr) ServerSocket() {
	addr := fmt.Sprintf(":%v", hbs.Port)
	server, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen socket failed, error is: %v\n", err)
		return
	}
	log.Println("listen success")

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Printf("accept failed, error is: %v\n", err)
			continue
		}

		go connHandler(conn)
	}
}