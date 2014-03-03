package cluster

import (
	"fmt"
	//"os"
	"time"
	//"io/ioutil"
	//"encoding/json"
	zmq "github.com/pebbe/zmq4"
	"strconv"
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
	//fmt.Println("Inside getPeers()")
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
	//fmt.Println("Inside getPort()")
	for _, peer := range c.Peers {
		if peer.Pid == myId {
			return peer.Port
		}
	}
	panic("myId" + string(myId) + "doesn't exist")
}

func (c *Config) getPeerInfo(myId int) map[int]Peer {
	//fmt.Println("Inside getPeerInfo()")
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
	//ip          int
	peers       []int
	outbox      chan *Envelope
	inbox       chan *Envelope
	stopInbox   chan bool
	stopOutbox  chan bool
	port        int
	peerInfo    map[int]Peer
	connections map[int]*zmq.Socket
}

func (s *Server) StopServer(){
	s.stopInbox<-true
	s.stopOutbox<-true
	time.Sleep(5*time.Second)
	close(s.stopInbox)
	close(s.stopOutbox)
	close(s.inbox)
	close(s.outbox)
	//fmt.Println("Server stopped")
}

func (s *Server) Pid() int {
	//fmt.Println("Inside Pid()")
	return s.pid
}

func (s *Server) Peers() []int {
	//fmt.Println("Inside Peers()")
	return s.peers
}

func (s *Server) Outbox() chan *Envelope {
	//fmt.Println("Inside Outbox()")
	return s.outbox
}

func (s *Server) Inbox() chan *Envelope {
	//fmt.Println("Inside Inbox()")
	return s.inbox
}

func (s *Server) Port() int {
	return s.port
}

func (s *Server) handleInbox() {
	//fmt.Println("handleInbox()")
	responder, err := zmq.NewSocket(zmq.PULL)
	if err != nil {
		fmt.Println("Error creating socket")
		panic(err)
	}
	defer responder.Close()
	bindAddress := "tcp://*:" + strconv.Itoa(s.port)
	responder.Bind(bindAddress)
	responder.SetRcvtimeo(1000*time.Millisecond)
	//fmt.Println("socket created")
	for {
		select{
			//goroutine should return when something is received on this channel
			case <-s.stopInbox:
				return
			//otherwise.. keep receiving and processing requests
			default:
				msg, err := responder.RecvBytes(0)
				//fmt.Println("Received")
				if err != nil {
					//fmt.Println("Error receiving message", err.Error())
					break
				}
				envelope := gobToEnvelope(msg)
				s.inbox <- &envelope
			}
		}
	}

func (s *Server) initializeSockets(){
	s.connections = make(map[int]*zmq.Socket)
	for i := range s.peers {
		peerId := s.peers[i]

		sock, err := zmq.NewSocket(zmq.PUSH)
		if err != nil {
			fmt.Println("Error creating socket", err)
		}
		sockAddr := "tcp://" + s.peerInfo[peerId].Ip + ":" + strconv.Itoa(s.peerInfo[peerId].Port)
		sock.Connect(sockAddr)
		s.connections[peerId] = sock
		//s.peerInfo[peerId].soc= s.connections[peerId]
	}
}

func (s *Server) handleOutbox() {
	//fmt.Println("handleOutbox()")
	for {
		select {
		case message := <-s.outbox:
			envelope := *message
			if envelope.Pid == BROADCAST {
				//fmt.Println("Broadcasting")
				time.Sleep(50*time.Millisecond)
				for _, conn := range s.connections {
					//TODO: insert individual pids in each envelope
					msg := envelopeToGob(envelope)
					conn.SendBytes(msg, 0)
				}
			} else {
				peerId := envelope.Pid
				conn := s.connections[peerId]
				msg := envelopeToGob(envelope)
				//fmt.Println(msg)
				conn.SendBytes(msg, 0)
			}
		case <-s.stopOutbox:
			return
		}
	}
}
