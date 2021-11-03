package main

import (
	"context"
	"errors"
	pb "example.com/chittychat/chatservice"
	"flag"
	"fmt"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
	glog "google.golang.org/grpc/grpclog"
	"log"
	"net"
	"os"
	"sync"
)

var (
	vectorclock = map[string]int64{}
	logger      glog.LoggerV2
	id          string
)

func init() {
	logger = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

type Connection struct {
	stream pb.Chittychat_JoinServer
	user   *pb.User
	active bool
	error  chan error
}

type Server struct {
	Connection []*Connection
	pb.UnimplementedChittychatServer
}

func (s *Server) Join(connection *pb.Connect, stream pb.Chittychat_JoinServer) error {
	wait := sync.WaitGroup{}
	done := make(chan int)

	conn := &Connection{
		stream: stream,
		user:   connection.User,
		active: true,
		error:  make(chan error),
	}
	s.Connection = append(s.Connection, conn)
	MergeVectorclocks(connection.Vectorclock)
	vectorclock[id]++
	logger.Info("Join request from: ", connection.User.Name, " | ", vectorclock)

	vectorclock[id]++
	welcomeMsg := &pb.Message{
		User: &pb.User{
			Id:   id,
			Name: "Server",
		},
		Msg:         fmt.Sprintf("Welcome %s, there are %d users online", connection.User.Name, len(s.Connection)),
		Vectorclock: vectorclock,
	}
	err := conn.stream.Send(welcomeMsg)
	logger.Info("User Joined: ", conn.user.Name, " | ", vectorclock)
	if err != nil {
		logger.Errorf("Error: %v - %v", conn.stream, err)
		conn.active = false
		conn.error <- err
	}

	for _, conn := range s.Connection {
		wait.Add(1)
		go func(msg *pb.Message, conn *Connection) {
			defer wait.Done()
			if conn.active && conn.user.Id != connection.User.Id {
				vectorclock[id]++
				msg.Vectorclock = vectorclock
				err := conn.stream.Send(msg)
				logger.Info("Send join message to: ", conn.user.Name, " | ", msg.Vectorclock)
				if err != nil {
					logger.Errorf("Error: %v - %v", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(&pb.Message{
			User: &pb.User{
				Id:   id,
				Name: "Server",
			},
			Msg:         fmt.Sprintf("User \"%s\" has joined | %d users online", connection.User.Name, len(s.Connection)),
			Vectorclock: vectorclock,
		}, conn)
	}

	go func() {
		wait.Wait()
		close(done)
	}()
	<-done

	return <-conn.error
}

func (s *Server) Broadcast(_ context.Context, msg *pb.Message) (*pb.Empty, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	MergeVectorclocks(msg.Vectorclock)
	vectorclock[id]++
	logger.Info("Received from: ", msg.User.Name, ", with message: \"", msg.Msg, "\" | ", vectorclock)

	for _, conn := range s.Connection {
		wait.Add(1)
		go func(msg *pb.Message, conn *Connection) {
			defer wait.Done()
			if conn.active && conn.user.Id != msg.User.Id {
				vectorclock[id]++
				msg.Vectorclock = vectorclock
				err := conn.stream.Send(msg)
				logger.Info("Send to: ", conn.user.Name, " | ", msg.Vectorclock)
				if err != nil {
					logger.Errorf("Error: %v - %v", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(msg, conn)
	}
	go func() {
		wait.Wait()
		close(done)
	}()
	<-done
	return &pb.Empty{}, nil
}

func (s *Server) Leave(_ context.Context, msg *pb.Message) (*pb.Empty, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	MergeVectorclocks(msg.Vectorclock)
	vectorclock[id]++
	logger.Info("Leave request from: ", msg.User.Name, " | ", vectorclock)

	for idx, conn := range s.Connection {
		if conn.user.Id == msg.User.Id {
			conn.active = false
			conn.error <- errors.New("disconnected")
			s.Connection = append(s.Connection[:idx], s.Connection[idx+1:]...)
		}
	}

	for _, conn := range s.Connection {
		wait.Add(1)
		go func(leaveMsg *pb.Message, leavingUser *pb.Message, conn *Connection) {
			defer wait.Done()
			if conn.active && conn.user.Id != msg.User.Id {
				vectorclock[id]++
				leaveMsg.Vectorclock = vectorclock
				err := conn.stream.Send(leaveMsg)
				logger.Info("Send leave message to: ", conn.user.Name, " | ", leaveMsg.Vectorclock)
				if err != nil {
					logger.Errorf("Error: %v - %v", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
				vectorclock[id]++
				leaveMsg.Vectorclock = vectorclock
				err2 := conn.stream.Send(leavingUser)
				if err2 != nil {
					logger.Errorf("Error: %v - %v", conn.stream, err2)
					conn.active = false
					conn.error <- err2
				}
			}
		}(&pb.Message{
			User: &pb.User{
				Id:   id,
				Name: "Server",
			},
			Msg:         fmt.Sprintf("User \"%s\" has left | %d users online", msg.User.Name, len(s.Connection)),
			Vectorclock: vectorclock,
		}, &pb.Message{
			User: &pb.User{
				Id:   msg.User.Id,
				Name: "leave request",
			},
			Vectorclock: nil,
		}, conn)
	}

	delete(vectorclock, msg.User.Id)

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
	return &pb.Empty{}, nil
}

func MergeVectorclocks(vc map[string]int64) {
	for id, clock := range vc {
		vectorclock[id] = max(vectorclock[id], clock)
	}
}

func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func main() {
	port := flag.String("p", "5000", "Specify port. Default is 5000")
	id = uuid.NewV4().String()
	flag.Parse()

	vectorclock[id] = 0

	var connections []*Connection
	server := &Server{connections, pb.UnimplementedChittychatServer{}}
	chatServer := grpc.NewServer()

	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("Error when starting server %v", err)
	}
	logger.Info("Started server on port " + *port)
	pb.RegisterChittychatServer(chatServer, server)
	regerr := chatServer.Serve(listener)
	if regerr != nil {
		log.Fatalf("Error when starting server %v", err)
	}
}
