package cluster

import (
	//"flag"
	//"fmt"
	"testing"
	"time"
	"fmt"
	//"github.com/sk4x0r/cluster"
)


func TestSimpleSendReceive(t *testing.T) {
	fmt.Println("Test: Simple Send Receive")
	sender := New(1001, PATH_TO_CONFIG)
	receiver := New(1002, PATH_TO_CONFIG)
	msgCount := 10
	
	success := make(chan bool, 1)
	go sendMessages(sender, msgCount, receiver.Pid())
	go receiveMessages(receiver, msgCount, success)
	select {
	case <-success:
		fmt.Println("SUCCESS!!")
		close(success)
		break
	case <-time.After(10 * time.Minute):
		t.Errorf("ERROR: TIMEOUT")
		break
	}
}

func TestSingleBroadcast(t *testing.T) {
	fmt.Println("Test: Single Broadcast")
	sender := New(1005, PATH_TO_CONFIG)
	receiver1 := New(1001, PATH_TO_CONFIG)
	receiver2 := New(1002, PATH_TO_CONFIG)
	receiver3 := New(1003, PATH_TO_CONFIG)
	//receiver4 := New(1004, PATH_TO_CONFIG)
	
	msgCount := 10
	
	success := make(chan bool, 1)
	go sendMessages(sender, msgCount, BROADCAST)
	go receiveMessages(receiver1, msgCount, success)
	go receiveMessages(receiver2, msgCount, success)
	go receiveMessages(receiver3, msgCount, success)
	//go receiveMessages(receiver4, msgCount, success)
	
	successCount:=0
	for successCount<3{
		select {
			case <-success:
				successCount=successCount+1
				break
			case <-time.After(10 * time.Minute):
				t.Errorf("Error:TIMEOUT")
				successCount=4
				break
		}
	}
	fmt.Println("SUCCESS!!")
}
