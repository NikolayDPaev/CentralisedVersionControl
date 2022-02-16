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

// Starts the server.
// Does not block.
func (s *Server) Start() {
	s.running = true
	go s.runServer()
}

// Stops the server gracefully.
// Waits for all clients to be serviced.
func (s *Server) Stop() {
	s.running = false
	if err := s.sendEmptyRequest(); err != nil {
		log.Println(err)
	}
	s.wg.Wait()
}

// Used to unblock the server from waiting on accept()
func (s *Server) sendEmptyRequest() error {
	c, err := net.Dial("tcp", "localhost:"+s.port)
	if err != nil {
		return fmt.Errorf("error creating poison socket: %w", err)
	}
	defer c.Close()

	comm := netio.NewCommunication(CHUNK_SIZE, c, c)
	if err := comm.SendVarInt(clienthandler.EMPTY_REQUEST); err != nil {
		return fmt.Errorf("error sending empty request: %w", err)
	}
	return nil
}

// Goroutine that handles incomming client.
// When its ready it closes the socket and decrements the wait group.
// If empty request is received (sent by the sendEmptyRequest function) it returns immediately.
func handleClient(c net.Conn, wg *sync.WaitGroup) {
	defer c.Close()
	defer wg.Done()
	comm := netio.NewCommunication(CHUNK_SIZE, c, c)
	clientHandler, err := clienthandler.NewHandler(comm, &storage.FileStorage{})

	if clientHandler == nil {
		return
	}
	if err != nil {
		log.Println(err)
		return
	}

	err = clientHandler.Handle()
	if err != nil {
		log.Println(err)
		return
	}
}

// Server routine.
// Incomming clients are handled by a separate goroutine
// and the counter in the wait group is incremented
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
