package main

import (
	"log"
	"net"
	"time"

	pb "grpc-bidirectional/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCommunicatorServer
}

func (s *server) Chat(stream pb.Communicator_ChatServer) error {
	log.Println("Server: Chat stream opened")
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			msg := &pb.Message{
				Sender:    "server",
				Content:   "Evento lato server",
				Timestamp: time.Now().Unix(),
			}
			if err := stream.Send(msg); err != nil {
				log.Println("Errore invio server:", err)
				return
			}
		}
	}()

	for {
		in, err := stream.Recv()
		if err != nil {
			log.Println("Client disconnesso:", err)
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
	pb.RegisterCommunicatorServer(grpcServer, &server{})
	log.Println("Server in ascolto su :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Errore avvio server: %v", err)
	}
}
