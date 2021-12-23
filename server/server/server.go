package server

import (
	"log"
	"net"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/clienthandler"
)

type Server struct {
	port          string
	notifyStopped chan struct{}
	running       bool
}

func NewServer(port string) *Server {
	return &Server{port: port, notifyStopped: make(chan struct{}), running: false}
}

func (s *Server) Start() {
	s.running = true
	go s.runServer()
}

func (s *Server) Stop() {
	s.running = false
	<-s.notifyStopped
}

func handleClient(c net.Conn) {
	if err := clienthandler.Communication(c, c); err != nil {
		log.Println(err)
	}
}

func (s *Server) runServer() {
	l, err := net.Listen("tcp4", s.port)
	if err != nil {
		log.Println(err)
		return
	}
	defer l.Close()

	for s.running {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go handleClient(c)
	}
	s.notifyStopped <- struct{}{}
}
