package heartbeatcli

import (
	"encoding/json"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type HeartBeatCli struct {
	TimeDur    time.Duration
	RemoteAddr string
	Msg        Msg
	atomic.Value
}

var (
	HeartBeatClient *HeartBeatCli
	once sync.Once
)


type Msg struct {
	CliName string `json:"cli_name"`
	Info    string `json:"info"`
}

func NewHeartBeatCli(remoteAddr string, dur time.Duration, nameCli string, info string) *HeartBeatCli {
	once.Do(func() {
		HeartBeatClient = &HeartBeatCli{
			TimeDur:    dur,
			RemoteAddr: remoteAddr,
			Msg: Msg{
				CliName: nameCli,
				Info:    info,
			},
		}
	})
	return HeartBeatClient
}

func (hb *HeartBeatCli) Run() error {
	conn, err := net.Dial("tcp", hb.RemoteAddr)
	if nil != err {
		log.Fatalf("dial tcp failed, error is: %v\n", err)
		return err
	}

	timer := time.NewTicker(hb.TimeDur)
	for  {
		select {
		case <- timer.C:
			err := hb.SendHeartBeat(conn)
			if nil != err {
				return err
			}
		case <- time.After(5 * hb.TimeDur):
			log.Fatal("oops, no data")
		}
	}

	return nil
}

func (hb *HeartBeatCli) SendHeartBeat (conn net.Conn) error {
	jsonv, err := json.Marshal(hb.Msg)
	if nil != err {
		log.Printf("format json failed, error is: %v\n", err)
		return err
	}
	log.Println(string(jsonv))

	_, err = conn.Write(jsonv)
	if nil != err {
		log.Printf("write data failed, error is: %v\n", err)
		return err
	}

	return nil
}
