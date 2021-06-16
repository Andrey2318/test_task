package server

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"net/url"
	"test_task/internal/service"
	"test_task/pkg/api"
)

type GRPCServer struct {
	server *Server
	*grpc.Server
}

func (gr *GRPCServer) Run() error {
	lis, err := net.Listen("tcp", gr.server.config.GRPCAddr)
	if err != nil {
		return err
	}

	p := make([]*url.URL, 0)

	for _, v := range gr.server.config.ProxyPool {
		u, err := url.Parse(v)
		if err != nil {
			return err
		}
		p = append(p, u)
	}

	se := service.New(p)

	api.RegisterRusprofileServiceServer(gr, se.RusprofileController())

	fmt.Printf("The GRPC Server is running on port %s\n", gr.server.config.GRPCAddr)
	return gr.Serve(lis)
}
