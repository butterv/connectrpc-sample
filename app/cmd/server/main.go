package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	greetv1 "github.com/butterv/connectrpc-sample/app/gen/greet/v1"
	"github.com/butterv/connectrpc-sample/app/gen/greet/v1/greetv1connect"
)

type GreetServer struct{}

func (s *GreetServer) Greet(
	ctx context.Context,
	req *connect.Request[greetv1.GreetRequest],
) (*connect.Response[greetv1.GreetResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

// withCORS adds CORS support to a Connect HTTP handler.
func withCORS(h http.Handler) http.Handler {
	middleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})
	return middleware.Handler(h)
}

func main() {
	mux := http.NewServeMux()
	path, handler := greetv1connect.NewGreetServiceHandler(
		&GreetServer{},
		connect.WithInterceptors(NewAuthInterceptor()),
	)
	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:8080",
		withCORS(
			// Use h2c so we can serve HTTP/2 without TLS.
			h2c.NewHandler(mux, &http2.Server{}),
		),
	)

	//api := http.NewServeMux()
	//api.Handle(greetv1connect.NewGreetServiceHandler(greeter))
	//mux := http.NewServeMux()
	//mux.Handle("/grpc/", http.StripPrefix("/grpc", api))
	//http.ListenAndServe(":8080", mux)
}
