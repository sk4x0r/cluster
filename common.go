package cluster

import (
	//"flag"
	//"fmt"
	//"testing"
	"time"
	//"github.com/sk4x0r/cluster"
)


var (
	PATH_TO_CONFIG = "config.json"
)

func createDummyMessage(msgId int, pid int) Envelope {
	e := Envelope{Pid: pid, MsgId: int64(msgId), Msg: "Dummy Message"}
	return e
}

func sendMessages(s Server, count int, pid int) {
	outbox := s.Outbox()
	for i := 0; i < count; i++ {
		msg := createDummyMessage(i, pid)
		outbox <- &msg
		time.Sleep(5 * time.Millisecond)
	}
	//fmt.Println("Sent ", count, " messages to ", pid)
}

func receiveMessages(s Server, count int, success chan bool) {
	inbox := s.Inbox()
	for i := 0; i < count; i++ {
		<-inbox
	}
	//fmt.Println("Received ", count, "messages")
	success <- true
}
