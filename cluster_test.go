package cluster

import (
	//"flag"
	"fmt"
	"testing"
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

func TestSendReceive(t *testing.T) {
	sender := New(1001, PATH_TO_CONFIG)
	receiver := New(1002, PATH_TO_CONFIG)
	msgCount := 100
	
	success := make(chan bool, 1)
	go sendMessages(sender, msgCount, receiver.Pid())
	go receiveMessages(receiver, msgCount, success)
	select {
	case <-success:
		fmt.Println("Send-Receive test passed successfully")
		break
	case <-time.After(5 * time.Minute):
		t.Errorf("Could not send ", msgCount, " messages in 5 minute")
		break
	}
}
