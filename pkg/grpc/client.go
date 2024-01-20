package ProtoBuf

/*
==================================================================================
Copyright (c) 2023 Capgemini .
File contains client handling and the exposed service calling.
==================================================================================
*/

import (
	"context"
	"flag"
	"fmt"
	"time"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = flag.String("addr", "service-ricplt-conflictmgr-grpc.ricplt:50051", "the address to connect to")

func SendE2GuidanceReq(e2GuidanceReq *E2GuidanceReq) (e2GuidanceResp *E2GuidanceResp, err error) {
	xapp.Logger.Debug("Inside SendE2GuidanceResponse: %v", e2GuidanceReq)
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		xapp.Logger.Error("GRPC Client did not connect: %v", err)
		return nil, err
	}
	fmt.Printf("\nGRPC Connection Established, conn: %v", conn)
	defer conn.Close()
	client := NewE2GuidanceServiceClient(conn)
	fmt.Printf("\nclient: %v", client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//Invokig Grpc FunctionCall for E2Guidance
	xapp.Logger.Debug("client : %v", client)
	e2GuidanceRes, err := client.E2GuidanceService(ctx, e2GuidanceReq)
	if err != nil {
		fmt.Printf("\nE2Guidance  fails: %v\n", err)
		return nil, err
	}
	return e2GuidanceRes, nil
}
