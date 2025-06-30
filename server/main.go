package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "grpc-bidirectional/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCommunicatorServer
	mu      sync.Mutex
	streams map[pb.Communicator_ChatServer]struct{}
}

func (s *server) broadcast(msg *pb.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for st := range s.streams {
		if err := st.Send(msg); err != nil {
			log.Println("Errore broadcast, rimuovo stream:", err)
			delete(s.streams, st)
		}
	}
}
func (s *server) Chat(stream pb.Communicator_ChatServer) error {
	// Registra lo stream
	s.mu.Lock()
	s.streams[stream] = struct{}{}
	s.mu.Unlock()
	log.Println("Client connesso, stream registrato")

	// Assicura deregistrazione alla chiusura dello stream
	defer func() {
		s.mu.Lock()
		delete(s.streams, stream)
		s.mu.Unlock()
		log.Println("Client disconnesso, stream rimosso")
	}()

	// Ciclo di ricezione messaggi client
	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Printf("Messaggio ricevuto: %s: %s\n", in.Sender, in.Content)
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Errore apertura porta: %v", err)
	}

	grpcServer := grpc.NewServer()
	s := &server{streams: make(map[pb.Communicator_ChatServer]struct{})}
	pb.RegisterCommunicatorServer(grpcServer, s)

	go func() { // Esecuzione simulata della pipeline
		for i := 0; i < 10; i++ {
			log.Printf("Esecuzione pipeline step %d", i)
			time.Sleep(3 * time.Second)
			evt := &pb.Message{
				Sender:    "server",
				Content:   fmt.Sprintf("Pipeline step %d completato", i),
				Timestamp: time.Now().Unix(),
			}
			s.broadcast(evt)
		}
	}()
	log.Println("Server in ascolto su :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Errore avvio server: %v", err)
	}
}
