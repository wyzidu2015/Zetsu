package main

import (
	"log"
	"net"
	"context"
	"google.golang.org/grpc"
	pb "Zetsu/zetsu"
	common "Zetsu/common"
)

const (
	port = ":5790"
)

type server struct {
	pb.UnimplementedZetsuServer
}

func (s *server) RegisterMonitor(ctx context.Context, in *pb.MachineConnectInfo) (*pb.StatusResponse, error) {
	log.Printf("Received connectinfo: %v\n", in)
	return &pb.StatusResponse{Status: int32(common.SUCCESS), Info: "aaa"}, nil
}

func (s *server) GetLatestConfig(ctx context.Context, in *pb.MachineBasicInfo) (*pb.ConfigResponse, error) {
	log.Printf("Received machine basic info: %v\n", in)
	return &pb.ConfigResponse{Interval: 100, MaxSaveTime: 60, Items: []*pb.ConfigItem{}}, nil
}

func (s *server) UploadMonitorItem(ctx context.Context, in *pb.MonitorInfo) (*pb.StatusResponse, error) {
	log.Printf("Received monitor info: %v\n", in)
	return &pb.StatusResponse{Status: int32(common.SUCCESS), Info: "bbbb"}, nil
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


