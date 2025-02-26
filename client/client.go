package main

import (
	"context"
	"flag"
	"goldun.com/quotes/auth"
	pb "goldun.com/quotes/quotes_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var (
	tls                = flag.Bool("tls", true, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "localhost:1443", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
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
	if *tls {
		if *caFile == "" {
			*caFile = auth.Path("ca_cert.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

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

	// Looking for multiple quotes for a given source
	getQuote(client, &pb.QuoteRequest{Source: "Семесюк"})

	// Looking for missing quote
	getQuote(client, &pb.QuoteRequest{Source: "unknown"})
}
