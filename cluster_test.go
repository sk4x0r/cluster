package cluster

import (
        //"flag"
        "fmt"
        "time"
        "testing"
        //"github.com/sk4x0r/cluster"
)
var(
	PATH_TO_CONFIG="config.json"
)

func createDummyMessage(msgId int, pid int) Envelope{
	e:=Envelope{Pid:pid, MsgId:int64(msgId), Msg:"Dummy Message"}
	return e
}

func sendMessages(s Server, count int, pid int){
	outbox:=s.Outbox()
	for i:=0;i<count;i++ {
		msg:=createDummyMessage(i,pid)
		outbox <- &msg
		time.Sleep(20*time.Millisecond)
	}
	fmt.Println("Sent ", count, " messages to ", pid)
}

func receiveMessages(s Server, count int){
	inbox:=s.Inbox()
	for i:=0;i<count;i++{
		//fmt.Println("Received",i)
		<- inbox
	}
	fmt.Println("Received ", count, "messages")
}

func TestSendReceive(t *testing.T) {
	sender:=New(1001, PATH_TO_CONFIG)
	receiver:=New(1002, PATH_TO_CONFIG)
	
	msgCount:=10
	go sendMessages(sender, msgCount, receiver.Pid())
	go receiveMessages(receiver, msgCount)
	time.Sleep(15*time.Second)
}
