package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	pb "stream/proto"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcPort = 50001
	httpPort = 50002
	wsPort   = 50003
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	go func() {
		if err := runGRPC(); err != nil {
			log.Panicln(err)
		}
	}()
	log.Printf("start grpc server, port: %d\n", grpcPort)

	go func() {
		if err := runHTTP(); err != nil {
			log.Panicln(err)
		}
	}()
	log.Printf("start http server, port: %d\n", httpPort)

	go func() {
		if err := runWebsocket(); err != nil {
			log.Panicln(err)
		}
	}()
	log.Printf("start websocket server, port: %d\n", wsPort)

	waitSignal()
}

func runGRPC() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return err
	}
	defer func() {
		if err = lis.Close(); err != nil {
			log.Println(err)
		}
	}()

	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, NewServer())

	return s.Serve(lis)
}

func runHTTP() error {
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("0.0.0.0:%d", grpcPort),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	gwmux := runtime.NewServeMux()
	if err = pb.RegisterStreamServiceHandler(context.Background(), gwmux, conn); err != nil {
		return err
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: cors.AllowAll().Handler(gwmux),
	}

	return gwServer.ListenAndServe()
}

func runWebsocket() error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterStreamServiceHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf(":%d", grpcPort), opts); err != nil {
		return err
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", wsPort), wsproxy.WebsocketProxy(mux))
}

func waitSignal() {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	log.Println("exit: ^C")
	sig := <-sigterm
	log.Printf("terminating: via signal %v\n", sig)
}
