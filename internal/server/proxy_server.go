package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"test_task/pkg/api"
)

type ProxyServer struct {
	server *Server
	*http.Server
}

func (ps *ProxyServer) Run(ctx context.Context) error {
	gw := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}))
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := api.RegisterRusprofileServiceHandlerFromEndpoint(ctx, gw, ps.server.config.GRPCAddr, opts)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", gw)
	mux.HandleFunc("/api.swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./swagger/api.swagger.json")
	})

	ps.Addr = ps.server.config.ProxyAddr
	ps.Handler = registerProxyServer(mux)

	fmt.Printf("The Proxy Server is running on port %s\n", ps.server.config.ProxyAddr)
	return ps.ListenAndServe()
}
