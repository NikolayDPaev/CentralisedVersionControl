package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/clienthandler"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/storage"
)

const CHUNK_SIZE = 4096

type Server struct {
	netInterface string
	port         string
	wg           sync.WaitGroup
	running      bool
}

func NewServer(netInterface, port string) *Server {
	return &Server{port: port, running: false}
}

func (s *Server) Start() {
	s.running = true
	go s.runServer()
}

func (s *Server) Stop() {
	s.running = false
	if err := s.sendEmptyRequest(); err != nil {
		log.Println(err)
	}
	s.wg.Wait()
}

func (s *Server) sendEmptyRequest() error {
	c, err := net.Dial("tcp", "localhost:"+s.port)
	if err != nil {
		return fmt.Errorf("error creating poison socket:\n%w", err)
	}
	defer c.Close()

	comm := netio.NewCommunicator(CHUNK_SIZE, c, c)
	if err := comm.SendVarInt(clienthandler.EMPTY_REQUEST); err != nil {
		return fmt.Errorf("error sending empty request:\n%w", err)
	}
	return nil
}

func handleClient(c net.Conn, wg *sync.WaitGroup) {
	comm := netio.NewCommunicator(CHUNK_SIZE, c, c)
	clientHandler, err := clienthandler.NewHandler(comm, &storage.FileStorage{})
	if err != nil {
		log.Println(err)
	}
	err = clientHandler.Handle()
	if err != nil {
		log.Println(err)
	}
	c.Close()
	wg.Done()
}

func (s *Server) runServer() {
	l, err := net.Listen("tcp4", s.netInterface+":"+s.port)
	if err != nil {
		log.Println(err)
		return
	}
	defer l.Close()

	for s.running {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		s.wg.Add(1)
		go handleClient(c, &s.wg)
	}
}
