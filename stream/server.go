package main

import (
	"fmt"
	"io"
	"log"
	pb "stream/proto"
	"strings"
	"time"
)

type server struct {
	pb.UnsafeStreamServiceServer
}

func NewServer() pb.StreamServiceServer {
	return &server{}
}

func (s *server) StreamNumber(req *pb.NumberRequest, stream pb.StreamService_StreamNumberServer) error {
	from := req.GetFrom()
	count := req.GetCount()
	log.Printf("number start from %d and count %d\n", from, count)

	now := from
	for i := 0; i < int(count); i++ {
		resp := &pb.Response{Result: fmt.Sprintf("here is your number: %d", now)}
		if err := stream.Send(resp); err != nil {
			log.Println(err)
			return err
		}
		log.Printf("finishing request number :%d", now)
		now++
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func (s *server) StreamPerson(stream pb.StreamService_StreamPersonServer) error {
	var people []string
	for {
		person, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return err
		}
		name := person.GetName()
		log.Printf("person(%s) from client\n", name)
		people = append(people, name)
	}

	joined := strings.Join(people, ", ")
	result := fmt.Sprintf("you sent %d people, %s", len(people), joined)
	if len(people) < 2 {
		result = fmt.Sprintf("you sent %d person: %s", len(people), joined)
	}

	if err := stream.SendAndClose(&pb.Response{Result: result}); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *server) StreamHello(stream pb.StreamService_StreamHelloServer) error {
	var people []string
	for {
		person, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return err
		}

		name := person.GetName()
		log.Printf("greeting from %s\n", name)
		people = append(people, name)
		if err = stream.Send(&pb.Response{Result: fmt.Sprintf("Hello, %s", name)}); err != nil {
			log.Println(err)
			return err
		}
	}

	if len(people) > 0 {
		log.Printf("bye %s", strings.Join(people, ", "))
	}
	return nil
}
