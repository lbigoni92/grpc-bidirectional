package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "grpc-bidirectional/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Errore creazione client: %v", err)
	}
	defer conn.Close()

	client := pb.NewCommunicatorClient(conn)
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalf("Errore creazione stream: %v", err)
	}

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Errore ricezione: %v", err)
			}
			fmt.Printf("Ricevuto: %s: %s\n", in.Sender, in.Content)
		}
	}()

	for i := 0; i < 5; i++ {
		msg := &pb.Message{
			Sender:    "client",
			Content:   fmt.Sprintf("Messaggio %d dal client", i),
			Timestamp: time.Now().Unix(),
		}
		if err := stream.Send(msg); err != nil {
			log.Fatalf("Errore invio: %v", err)
		}
		time.Sleep(3 * time.Second)
	}

	stream.CloseSend()
}
