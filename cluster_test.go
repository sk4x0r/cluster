package cluster

import (
	//"flag"
	"fmt"
	"testing"
	"time"
	"strconv"
	//"github.com/sk4x0r/cluster"
)


func TestSendReceive(t *testing.T) {
	sender := New(1001, PATH_TO_CONFIG)
	receiver := New(1002, PATH_TO_CONFIG)
	msgCount := 10
	
	success := make(chan bool, 1)
	go sendMessages(sender, msgCount, receiver.Pid())
	go receiveMessages(receiver, msgCount, success)
	select {
	case <-success:
		fmt.Println("Send-Receive test passed successfully")
		time.Sleep(5*time.Second)
		break
	case <-time.After(10 * time.Minute):
		t.Errorf("Could not send ", strconv.Itoa(msgCount), " messages in 10 minute")
		break
	}
}
