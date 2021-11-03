package main

import (
	"bufio"
	"context"
	pb "example.com/chittychat/chatservice"
	"flag"
	"fmt"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"sync"
)

var (
	client      pb.ChittychatClient
	wait        *sync.WaitGroup
	vectorclock = map[string]int64{}
	user        *pb.User
)

func init() {
	wait = &sync.WaitGroup{}
}

func connect() error {
	var streamErr error
	vectorclock[user.Id]++
	stream, err := client.Join(context.Background(), &pb.Connect{
		User:        user,
		Active:      true,
		Vectorclock: vectorclock,
	})

	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}

	wait.Add(1)
	go func(streamClient pb.Chittychat_JoinClient) {
		defer wait.Done()
		for {
			msg, err := streamClient.Recv()
			if err != nil {
				streamErr = fmt.Errorf("recieve error: %v", err)
				break
			}
			if msg.User.Name == "leave request" && msg.Msg == "" && msg.Vectorclock == nil {
				vectorclock[user.Id]++
				delete(vectorclock, msg.User.Id)
			} else {
				MergeVectorclocks(msg.Vectorclock)
				vectorclock[user.Id]++
				fmt.Printf("[%v]: %s | %v\n", msg.User.Name, msg.Msg, vectorclock)
			}
		}
		if streamErr != nil {
			if streamErr.Error() == "recieve error: rpc error: code = Unknown desc = disconnected" {
				fmt.Println("[Disconnected]")
				os.Exit(0)
			}
		}
	}(stream)

	return streamErr
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
	host := flag.String("h", "localhost", "Host server address")
	port := flag.String("p", "5000", "Host server port")
	name := flag.String("u", "Anonymous", "Username")
	id := uuid.NewV4().String()
	flag.Parse()

	done := make(chan int)
	vectorclock[id] = 0

	conn, err := grpc.Dial(*host+":"+*port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}

	client = pb.NewChittychatClient(conn)
	user = &pb.User{
		Id:   id,
		Name: *name,
	}

	connErr := connect()
	if connErr != nil {
		os.Exit(1)
	}

	wait.Add(1)
	go func() {
		defer wait.Done()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			vectorclock[user.Id]++
			var msg *pb.Message
			var err error
			if scanner.Text() == "/leave" {
				msg = &pb.Message{
					User:        user,
					Vectorclock: vectorclock,
				}
				_, err = client.Leave(context.Background(), msg)
			} else {
				msg = &pb.Message{
					User:        user,
					Msg:         scanner.Text(),
					Vectorclock: vectorclock,
				}
				_, err = client.Broadcast(context.Background(), msg)
			}
			if err != nil {
				fmt.Printf("Send Error: %v", err)
				break
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			msg := &pb.Message{
				User:        user,
				Vectorclock: vectorclock,
			}
			vectorclock[user.Id]++
			_, err := client.Leave(context.Background(), msg)
			if err != nil {
				fmt.Printf("Send Error: %v", err)
				break
			}
			os.Exit(0)
		}
	}()

	go func() {
		wait.Wait()
		close(done)
	}()
	<-done
}
