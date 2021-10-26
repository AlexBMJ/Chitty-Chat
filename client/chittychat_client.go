package main

import (
    "context"
    pb "example.com/chittychat/chatservice"
    "google.golang.org/grpc"
    "log"
    "strconv"
)

const (
	address = "localhost:50051"
)

var (
    id int64
    vectorclock = make([]int64, 1)
)

func main() {
    var conn, err = grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    var client = pb.NewChittychatClient(conn)

    var ctx = context.Background()

    var res, joinErr = client.Join(ctx, &pb.Message{ Vectorclock: vectorclock })
    if joinErr != nil {
        log.Fatalf("could not join chittychat: %v", joinErr)
    }

    recv, err := res.Recv()
    if err != nil {
        log.Fatal(err)
        return
    }
    id, _ = strconv.ParseInt(recv.Text, 10, 64)
    MergeVectorclocks(recv.Vectorclock)

    go func() {
        for {
            recv, err := res.Recv()
            MergeVectorclocks(recv.Vectorclock)
            vectorclock[id]++
            if err != nil {
                log.Fatal(err)
                return
            }
            log.Printf("%v:    %v", recv.Text, vectorclock)
        }
    }()


    for {
    }


    //id = res.Id
    //MergeVectorclocks(res.Vectorclock)
    //vectorclock[id]++
    //log.Printf("Join Response: %v", vectorclock)

    //for i := 0; i < 5; i++ {
    //    var _, pubErr = client.Publish(ctx, &pb.TextMessage{ UId: id, Text: "testMessage", Vectorclock: vectorclock })
    //    if pubErr != nil {
    //        log.Fatal(pubErr)
    //    }
    //    vectorclock[id]++
    //    time.Sleep(2 * time.Second)
    //}
}

func MergeVectorclocks(vc []int64) {
    if len(vc) > len(vectorclock) {
        vectorclock = append(vectorclock, make([]int64, len(vc) - len(vectorclock))...)
    } else if len(vc) < len(vectorclock) {
        vc = append(vc, make([]int64, len(vectorclock) - len(vc))...)
    }
    for i := 0; i < len(vc); i++ {
        if vc[i] > vectorclock[i] {
            vectorclock[i] = vc[i]
        }
    }
    vectorclock[id]++
}
