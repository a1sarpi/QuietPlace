package main

import (
	"fmt"
	"github.com/a1sarpi/QuietPlace/currency/data"
	protos "github.com/a1sarpi/QuietPlace/currency/protos/currency"
	"github.com/a1sarpi/QuietPlace/currency/server"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

func main() {
	log := hclog.Default()

	rates, err := data.NewRates(log)
	if err != nil {
		log.Error("Unable to generate rates", "error", err)
	}

	// create a new gRPC server, use WithInsecure to allow http connections
	gs := grpc.NewServer()

	// register an instance of the Currency server
	c := server.NewCurrency(rates, log)

	// register the currency server
	protos.RegisterCurrencyServer(gs, c)

	// register the reflection service which allows clients to determine the methods
	// for this gRPC server
	reflection.Register(gs)

	// create a TCP socket for inbound server connections
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 9092))
	if err != nil {
		log.Error("Unable to listen", "err", err)
		os.Exit(1)
	}

	gs.Serve(l)
}
