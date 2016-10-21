package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	//"github.com/lightstep/lightstep-tracer-go"
	//stdopentracing "github.com/opentracing/opentracing-go"
	//zipkin "github.com/openzipkin/zipkin-go-opentracing"
	//appdashot "github.com/sourcegraph/appdash/opentracing"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	//"sourcegraph.com/sourcegraph/appdash"

	"github.com/pkg/errors"

	//"github.com/go-kit/kit/log"

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
	// provide the direct address of an addsvc. This presumption is reflected in
	// the addcli binary and the the client packages: the -transport.addr flags
	// and various client constructors both expect host:port strings. For an
	// example service with a client built on top of a service discovery system,
	// see profilesvc.

	var (
		httpAddr = flag.String("http.addr", "", "HTTP address of addsvc")
		grpcAddr = flag.String("grpc.addr", ":8082", "gRPC (HTTP) address of addsvc")
		//zipkinAddr     = flag.String("zipkin.addr", "", "Enable Zipkin tracing via a Kafka Collector host:port")
		//appdashAddr    = flag.String("appdash.addr", "", "Enable Appdash tracing via an Appdash server host:port")
		//lightstepToken = flag.String("lightstep.token", "", "Enable LightStep tracing via a LightStep access token")
		method = flag.String("method", "readcontexttestvalue", "readcontexttestvalue,readcontextmetadata")
	)

	var ()
	flag.Parse()

	var (
		service handler.Service
		err     error
	)
	if *httpAddr != "" {
		//service, err = httpclient.New(*httpAddr, tracer, log.NewNopLogger())
		service, err = httpclient.New(*httpAddr)
	} else if *grpcAddr != "" {
		conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while dialing grpc connection: %v", err)
			os.Exit(1)
		}
		defer conn.Close()
		service = grpcclient.New(conn /*, tracer, log.NewNopLogger()*/)
	} else {
		fmt.Fprintf(os.Stderr, "error: no remote address specified\n")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var ctx = context.Background()
	ctx = context.WithValue(ctx, "test", "truss")
	md := metadata.Pairs("test", "truss")
	ctx = metadata.NewContext(ctx, md)

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

	case "readcontextmetadata":

		var err error

		request, err := clientHandler.ReadContextMetadata()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calling clientHandler.ReadContextMetadata: %v\n", err)
			os.Exit(1)
		}

		v, err := service.ReadContextMetadata(ctx, request)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calling service.ReadContextMetadata: %v\n", err)
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
