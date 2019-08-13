package heartbeatcli

import (
	"encoding/json"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"github.com/jigetage/goheartbeat/common"
)

type HeartBeatCli struct {
	TimeDur    time.Duration
	RemoteAddr string
	Msg        common.Msg
	atomic.Value
}

var (
	HeartBeatClient *HeartBeatCli
	once sync.Once
)

func NewHeartBeatCli(remoteAddr string, dur time.Duration, nameCli string, info string) *HeartBeatCli {
	once.Do(func() {
		HeartBeatClient = &HeartBeatCli{
			TimeDur:    dur,
			RemoteAddr: remoteAddr,
			Msg: common.Msg{
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
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Duration(common.Timeout) * time.Second))

	timer := time.NewTicker(hb.TimeDur)
	for  {
		select {
		case <- timer.C:
			hb.SendHeartBeat(conn)
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
	log.Println("msg send: " + string(jsonv))

	// send heartbeat msg
	_, err = conn.Write(jsonv)
	if nil != err {
		log.Printf("write data failed, error is: %v\n", err)
		return err
	}

	// receive msg
	buf := make([]byte, common.RecvBuf)
	cnt, err := conn.Read(buf)
	if nil != err || 0 == cnt {
		// set deadline duration again
		//conn.SetDeadline(time.Now().Add(time.Duration(common.Timeout) * time.Second))
		log.Printf("read data failed, error is: %v\n", err)
		return err
	}
	// process recv msg
	log.Println("msg recv: " + string(buf[:cnt]))

	return nil
}
