package cluster

import (
	//"flag"
	//"fmt"
	//"testing"
	"time"
	"log"
	"bytes"
	"encoding/gob"
	"strconv"
	"encoding/json"
	//"github.com/sk4x0r/cluster"
)

const (
	BROADCAST = -1
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
		time.Sleep(50 * time.Millisecond)
	}
	s.StopServer()
	//fmt.Println("Sent ", count, " messages to ", pid)
}

func receiveMessages(s Server, count int, success chan bool) {
	inbox := s.Inbox()
	for i := 0; i < count; i++ {
		<-inbox
	}
	s.StopServer()
	//fmt.Println("Received ", count, "messages")
	success <- true
}

//cite:http://blog.golang.org/gobs-of-data
//cite: Pushkar Khadilkar
func gobToEnvelope(gobbed []byte) Envelope {
	buf := bytes.NewBuffer(gobbed)
	dec := gob.NewDecoder(buf)
	var envelope Envelope
	dec.Decode(&envelope)
	return envelope
}

//cite: http://blog.golang.org/gobs-of-data
//cite: Pushkar Khadilkar
func envelopeToGob(envelope Envelope) []byte {
	var buf bytes.Buffer
    enc := gob.NewEncoder(&buf)
    err:=enc.Encode(envelope)
    if err != nil {
        log.Fatal("encode error:", err)
    }
	return buf.Bytes()
}

//depreciated, use envelopeToGob instead
func envelopeToMsg(e Envelope, peerId int) string {
	message := "{"
	pid := strconv.Itoa(peerId)
	message += "\"Pid\":" + pid + ","
	msgId := strconv.FormatInt(e.MsgId, 10)
	message += "\"MsgId\":" + msgId + ","
	msg := e.Msg.(string)
	message += "\"Msg\":\"" + msg
	message += "\"}"
	return message
}

//depreciated, use gobToEnvelope instead
func msgToEnvelope(msg string) Envelope {
	//fmt.Println("Inside msgToEnvelope")
	var env Envelope
	json.Unmarshal([]byte(msg), &env)
	return env
}
