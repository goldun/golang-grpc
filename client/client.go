package main

import (
	"context"
	"flag"
	pb "goldun.com/quotes/quotes_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var (
	serverAddr = flag.String("addr", "localhost:8080", "The server address in the format of host:port")
)

func getQuote(client pb.QuotesServiceClient, request *pb.QuoteRequest) {
	log.Printf("Getting quotes for a source (%s)", request.Source)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	quote, err := client.GetQuote(ctx, request)
	if err != nil {
		log.Fatalf("client.GetQuote failed: %v", err)
	}
	log.Println(quote)
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}(conn)
	client := pb.NewQuotesServiceClient(conn)

	// Looking for existing quote
	getQuote(client, &pb.QuoteRequest{Source: "Lock, Stock and Two Smoking Barrels"})

	// Looking for missing quote
	getQuote(client, &pb.QuoteRequest{Source: ""})
}
