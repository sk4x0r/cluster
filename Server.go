package cluster

import (
	"fmt"
	//"os"
	//"io/ioutil"
	"encoding/json"
	zmq "github.com/pebbe/zmq4"
	"strconv"
)

const (
	BROADCAST = -1
)

type Peer struct {
	Pid  int
	Ip   string
	Port int
	soc  *zmq.Socket
}

type Config struct {
	Peers []Peer
}

func (c *Config) getPeers(myId int) []int {
	//speak("Inside getPeers()")
	l := len(c.Peers)
	pids := make([]int, l-1)

	i := 0
	for _, peer := range c.Peers {
		if peer.Pid != myId {
			pids[i] = peer.Pid
			i = i + 1
		}
	}
	return pids
}

func (c *Config) getPort(myId int) int {
	//speak("Inside getPort()")
	for _, peer := range c.Peers {
		if peer.Pid == myId {
			return peer.Port
		}
	}
	panic("myId" + string(myId) + "doesn't exist")
}

func (c *Config) getPeerInfo(myId int) map[int]Peer {
	//speak("Inside getPeerInfo()")
	peerInfo := make(map[int]Peer)
	for _, peer := range c.Peers {
		if peer.Pid != myId {
			peerInfo[peer.Pid] = peer
		}
	}
	return peerInfo
}

type Server struct {
	pid         int
	peers       []int
	outbox      chan *Envelope
	inbox       chan *Envelope
	port        int
	peerInfo    map[int]Peer
	connections map[int]*zmq.Socket
}

func (s *Server) Pid() int {
	//speak("Inside Pid()")
	return s.pid
}

func (s *Server) Peers() []int {
	//speak("Inside Peers()")
	return s.peers
}

func (s *Server) Outbox() chan *Envelope {
	//speak("Inside Outbox()")
	//fmt.Println("Inside Outbox()")
	return s.outbox
}

func (s *Server) Inbox() chan *Envelope {
	//speak("Inside Inbox()")
	return s.inbox
}

func (s *Server) Port() int {
	return s.port
}

func msgToEnvelope(msg string) Envelope {
	//speak("Inside msgToEnvelope")
	var env Envelope
	json.Unmarshal([]byte(msg), &env)
	return env
}

func (s *Server) handleInbox() {
	//speak("handleInbox()")
	responder, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		fmt.Println("Error creating socket", err)
		return
	}
	defer responder.Close()
	bindAddress := "tcp://*:" + strconv.Itoa(s.port)
	responder.Bind(bindAddress)
	//fmt.Println("socket created")
	for {
		//fmt.Println("waiting for response")
		msg, err := responder.Recv(0)
		if err != nil {
			fmt.Println("Error receiving message", err.Error())
			break
		}
		envelope := msgToEnvelope(msg)
		s.inbox <- &envelope
		responder.Send("!", 0)
		//fmt.Println("sent response")
	}
}

func envelopeToMsg(e Envelope, peerId int) string {
	message := "{"
	pid := strconv.Itoa(peerId)
	message += "\"Pid\":" + pid + ","
	msgId := strconv.FormatInt(e.MsgId, 10)
	message += "\"MsgId\":" + msgId + ","
	msg := e.Msg.(string)
	message += "\"Msg\":\"" + msg
	message += "\"}"
	//fmt.Println("message=",message)
	return message
}

func (s *Server) handleOutbox() {
	//speak("handleOutbox()")
	//create sockets for each peer
	s.connections = make(map[int]*zmq.Socket)
	for i := range s.peers {
		peerId := s.peers[i]

		sock, err := zmq.NewSocket(zmq.REQ)
		if err != nil {
			fmt.Println("Error creating socket", err)
		}
		st := "tcp://" + s.peerInfo[peerId].Ip + ":" + strconv.Itoa(s.peerInfo[peerId].Port)
		sock.Connect(st)
		s.connections[peerId] = sock
		//s.peerInfo[peerId].soc= s.connections[peerId]
	}
	for {
		select {
		case message := <-s.outbox:
			envelope := *message
			if envelope.Pid == BROADCAST {
				for peerId, conn := range s.connections {
					msg := envelopeToMsg(envelope, peerId)
					conn.Send(msg, 0)
					conn.Recv(0)
				}
			} else {
				peerId := envelope.Pid
				conn := s.connections[peerId]
				msg := envelopeToMsg(envelope, peerId)
				//fmt.Println(msg)
				conn.Send(msg, 0)
				conn.Recv(0)
			}
		}
	}
}
