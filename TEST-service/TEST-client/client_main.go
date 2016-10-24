package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/pkg/errors"

	// This Service
	pb "github.com/hasAdamr/test-service/TEST-service"
	grpcclient "github.com/hasAdamr/test-service/TEST-service/generated/client/grpc"
	httpclient "github.com/hasAdamr/test-service/TEST-service/generated/client/http"
	clientHandler "github.com/hasAdamr/test-service/TEST-service/handlers/client"
	handler "github.com/hasAdamr/test-service/TEST-service/handlers/server"
)

var (
	_ = strconv.ParseInt
	_ = strings.Split
	_ = json.Compact
	_ = errors.Wrapf
	_ = pb.RegisterTestServiceServer
)

func main() {
	// The addcli presumes no service discovery system, and expects users to
	// provide the direct address of an service. This presumption is reflected in
	// the cli binary and the the client packages: the -transport.addr flags
	// and various client constructors both expect host:port strings.

	var (
		httpAddr = flag.String("http.addr", "", "HTTP address of addsvc")
		grpcAddr = flag.String("grpc.addr", ":5040", "gRPC (HTTP) address of addsvc")
		method   = flag.String("method", "readcontexttestvalue", "readcontexttestvalue")
	)

	var ()
	flag.Parse()

	var (
		service handler.Service
		err     error
	)
	if *httpAddr != "" {
		service, err = httpclient.New(*httpAddr, httpclient.CtxValuesToSend([]string{"Auth-truss"}))
	} else if *grpcAddr != "" {
		conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while dialing grpc connection: %v", err)
			os.Exit(1)
		}
		defer conn.Close()
		service = grpcclient.New(conn)
	} else {
		fmt.Fprintf(os.Stderr, "error: no remote address specified\n")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	ctx := context.WithValue(context.Background(), "Auth-truss", "kinda secret")

	switch *method {

	case "readcontexttestvalue":
		var err error

		request, err := clientHandler.ReadContextTestValue()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calling clientHandler.ReadContextTestValue: %v\n", err)
			os.Exit(1)
		}

		v, err := service.ReadContextTestValue(ctx, request)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calling service.ReadContextTestValue: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Client Requested with:")
		fmt.Println()
		fmt.Println("Server Responded with:")
		fmt.Println(v)
	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", method)
		os.Exit(1)
	}
}
