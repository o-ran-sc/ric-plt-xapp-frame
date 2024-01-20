package ProtoBuf

/*
==================================================================================
Copyright (c) 2023 Capgemini .
File contains server handling of GRPC and of implemented module handling.
==================================================================================
*/

import (
	"context"
	"flag"
	"fmt"
	"net"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"google.golang.org/grpc"
)

type E2GuidanceReqHandle func(r *E2GuidanceReq) (*E2GuidanceResp, error)

var (
	port = flag.Int("port", 50051, "The server port")
)

type myServiceServer struct {
	UnimplementedE2GuidanceServiceServer
}

var E2GuidanceReqhandleInstance E2GuidanceReqHandle

/*
--------------------------------------------
Server calls this function to pass its handler for invoking on receive client request
--------------------------------------------
*/
func Listen(in E2GuidanceReqHandle) {
	xapp.Logger.Debug("Inside listen")
	E2GuidanceReqhandleInstance = in
	return
}
func (s *myServiceServer) E2GuidanceService(ctx context.Context, req *E2GuidanceReq) (resp *E2GuidanceResp, err error) {
	xapp.Logger.Debug("E2GuidanceService Received: %v", req)
	response, err := E2GuidanceReqhandleInstance(req)
	if err != nil {
		xapp.Logger.Error("Unable to process Guidance Request : %v", err)
		return nil, err
	}
	return response, nil

}

/*
--------------------------------------
This function starts GRPC server
-----------------------------------------
*/
func StartGrpcServer() {
	xapp.Logger.Info("Inside StartGrpcServer")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		xapp.Logger.Error("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterE2GuidanceServiceServer(s, &myServiceServer{})

	xapp.Logger.Debug("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		xapp.Logger.Error("failed to serve: %v", err)
	}
}
