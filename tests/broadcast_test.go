package cluster

import (
	//"flag"
	"fmt"
	"testing"
	"time"
	//"github.com/sk4x0r/cluster"
)

func TestBroadcast(t *testing.T) {
	sender := New(1005, PATH_TO_CONFIG)
	receiver1 := New(1001, PATH_TO_CONFIG)
	receiver2 := New(1002, PATH_TO_CONFIG)
	receiver3 := New(1003, PATH_TO_CONFIG)
	receiver4 := New(1004, PATH_TO_CONFIG)
	
	msgCount := 100
	
	success := make(chan bool, 1)
	go sendMessages(sender, msgCount, BROADCAST)
	go receiveMessages(receiver1, msgCount, success)
	go receiveMessages(receiver2, msgCount, success)
	go receiveMessages(receiver3, msgCount, success)
	go receiveMessages(receiver4, msgCount, success)
	
	successCount:=0
	for successCount<4{
		select {
			case <-success:
				successCount=successCount+1
				if successCount==4{
					fmt.Println("Broadcast Test passed successfully")
				}
				break
			case <-time.After(10 * time.Minute):
				t.Errorf("Could not send ", msgCount, " messages in 10 minute")
				successCount=4
				break
		}
	}
}
