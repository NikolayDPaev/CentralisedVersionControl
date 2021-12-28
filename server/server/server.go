package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/clienthandler"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
)

type Server struct {
	port    string
	wg      sync.WaitGroup
	running bool
}

func NewServer(port string) *Server {
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
		return fmt.Errorf("error creating poison socket: %w", err)
	}
	defer c.Close()
	if err := netIO.SendVarInt(clienthandler.EMPTY_REQUEST, c); err != nil {
		return fmt.Errorf("error sending empty request: %w", err)
	}
	return nil
}

func handleClient(c net.Conn, wg *sync.WaitGroup) {
	if err := clienthandler.Communication(c, c); err != nil {
		log.Println(err)
	}
	c.Close()
	wg.Done()
}

// func logUnwrappedError(err error) {
// 	currentErr := err
// 	for errors.Unwrap(currentErr) != nil {
// 		currentErr = errors.Unwrap(currentErr)
// 		log.Println(err)
// 	}
// }

func (s *Server) runServer() {
	l, err := net.Listen("tcp4", ":"+s.port)
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
