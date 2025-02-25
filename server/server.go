package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"

	pb "goldun.com/quotes/quotes_proto"
)

var (
	port       = flag.Int("port", 8080, "The server port")
	quotesFile = flag.String("quotes_file", "testdata/quotes.json", "The file with quotes")
)

type quotesServiceServer struct {
	pb.UnimplementedQuotesServiceServer
	quotes map[string][]*pb.Quote
}

func (s *quotesServiceServer) GetQuote(ctx context.Context, request *pb.QuoteRequest) (*pb.QuoteResponse, error) {
	q, exist := s.quotes[request.Source]
	if !exist {
		return nil, errors.New("quote for a given source isn't found")
	}
	return &pb.QuoteResponse{Quote: q}, nil
}

func (s *quotesServiceServer) loadQuotes(filePath string) {
	var data []byte
	var q []*pb.Quote
	var err error
	data, err = os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
	if err := json.Unmarshal(data, &q); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
	log.Printf("quotes: %v /n", q)
	for _, val := range q {
		key := val.Source
		s.quotes[key] = append(s.quotes[key], val)
	}
}

func newServer() *quotesServiceServer {
	s := &quotesServiceServer{quotes: make(map[string][]*pb.Quote)}
	s.loadQuotes(*quotesFile)
	return s
}

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterQuotesServiceServer(grpcServer, newServer())
	log.Printf("listening on port: %d \n", *port)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
		return
	}
}
