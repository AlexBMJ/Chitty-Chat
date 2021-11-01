package main

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	pb "example.com/chittychat/chatservice"
	"google.golang.org/grpc"
	glog "google.golang.org/grpc/grpclog"
)

var logger glog.LoggerV2

func init() {
	logger = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

type Connection struct {
	stream pb.Chittychat_JoinServer
	id     string
	active bool
	error  chan error
}

type Server struct {
	Connection []*Connection
}

var (
	pCount      int64 = 0
	vectorclock       = make([]int64, 1)
)

func (s *Server) Join(connection *pb.Connect, stream pb.Chittychat_JoinServer) error {
	conn := &Connection{
		stream: stream,
		id:     connection.User.Id,
		active: true,
		error:  make(chan error),
	}
	s.Connection = append(s.Connection, conn)
	return <-conn.error
}

//vectorclock[0]++
//pCount++
//vectorclock = append(vectorclock, message.Vectorclock[0])
//var vc = append(make([]int64, pCount), message.Vectorclock...)
//MergeVectorclocks(vc)
//log.Printf("Participant %v Joined Chitty-Chat at Vector time %v", pCount, vectorclock)
//
//err := src.Send(&pb.TextMessage{Text: strconv.FormatInt(pCount, 10), Vectorclock: vectorclock})
//if err != nil {
//	return err
//}
//
//for {
//	err := src.Send(&pb.TextMessage{Text: "aaaa", Vectorclock: vectorclock})
//	if err != nil {
//		return err
//	}
//	vectorclock[0]++
//
//	time.Sleep(2 * time.Second)
//}

func MergeVectorclocks(vc []int64) {
	if len(vc) > len(vectorclock) {
		vectorclock = append(vectorclock, make([]int64, len(vc)-len(vectorclock))...)
	} else if len(vc) < len(vectorclock) {
		vc = append(vc, make([]int64, len(vectorclock)-len(vc))...)
	}
	for i := 0; i < len(vc); i++ {
		if vc[i] > vectorclock[i] {
			vectorclock[i] = vc[i]
		}
	}
	vectorclock[0]++
}

func (s *ChittychatServer) Publish(ctx context.Context, message *pb.TextMessage) (*pb.Empty, error) {
	vectorclock[0]++
	MergeVectorclocks(message.Vectorclock)
	vectorclock[0]++
	log.Printf("%v : %v", message.Text, vectorclock)
	return &pb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterChittychatServer(server, &ChittychatServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
