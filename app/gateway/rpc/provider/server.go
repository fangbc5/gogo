package provider

import (
	context "context"

	"google.golang.org/grpc"
)

type ApiGrpcServer struct {
}

func (s ApiGrpcServer) Health(c context.Context, in *Req) (*Rsp, error) {
	return &Rsp{Code: 200, Msg: "ok", Data: nil}, nil
}


func GetGrpcServer() *grpc.Server {
	server := grpc.NewServer()
	RegisterApiServer(server, ApiGrpcServer{})
	return server
}
