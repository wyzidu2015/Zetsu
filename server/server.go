package main

import (
	common "Zetsu/common"
	pb "Zetsu/zetsu"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":5790"
)

type server struct {
	pb.UnimplementedZetsuServer
}

func (s *server) RegisterMonitor(ctx context.Context, in *pb.MachineConnectInfo) (*pb.StatusResponse, error) {
	log.Printf("Received connectinfo: %v\n", in)
	return &pb.StatusResponse{Status: int32(common.SUCCESS), Info: ""}, nil
}

func (s *server) GetLatestConfig(ctx context.Context, in *pb.MachineBasicInfo) (*pb.ConfigResponse, error) {
	log.Printf("Received machine basic info: %v\n", in)
	configItems := []*pb.ConfigItem{
		{Type: pb.ConfigItem_CPU},
		{Type: pb.ConfigItem_MEM},
		{Type: pb.ConfigItem_NET},
	}

	return &pb.ConfigResponse{Interval: 5, MaxSaveTime: 60, Items: configItems}, nil
}

func (s *server) UploadMonitorItem(ctx context.Context, in *pb.MonitorInfo) (*pb.StatusResponse, error) {
	log.Printf("Received monitor info: %v\n", in)
	return &pb.StatusResponse{Status: int32(common.SUCCESS), Info: ""}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen(%v): %v", port, err)
	}

	ser := grpc.NewServer()
	pb.RegisterZetsuServer(ser, &server{})
	if err := ser.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
