package monitor

import (
	"log"
	"fmt"
	"time"
	"context"
	"google.golang.org/grpc"
	pb "Zetsu/zetsu"
)


type Monitor struct {
	Client pb.ZetsuClient
	Ctx context.Context
	Cancel context.CancelFunc
}

func NewMonitor(address string) *Monitor {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Can't connect: %v", err)
	}

	cli := pb.NewZetsuClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return &Monitor{Client: cli, Ctx: ctx, Cancel: cancel}
}

func (mon *Monitor) Register(host *HostCollector) error {
	log.Println("In register func")

	r, err := mon.Client.RegisterMonitor(mon.Ctx, &pb.MachineConnectInfo{Uri: fmt.Sprintf("%s:%d", host.IP, 8765), GroupId: 123})
	if err != nil {
		log.Fatalf("Can't parse: %v", err)
	}
	log.Printf("Get message from register monitor: %v\n", r)
	return nil
}

func (mon *Monitor) GetConfig(host *HostCollector) error {
	log.Println("In get config func")
	
	mr, err := mon.Client.GetLatestConfig(mon.Ctx, &pb.MachineBasicInfo{CpuArch: host.CpuArch, CpuCores: host.CpuCores, MemorySize: host.MemorySize, OsType: host.OS})
	if err != nil {
		log.Fatalf("Can't parse: %v", err)
	}

	log.Printf("Get message from get latest config: %v\n", mr)
	return nil
}

func (mon *Monitor) UploadInfo() error {
	log.Println("In upload info func")

	ur, err := mon.Client.UploadMonitorItem(mon.Ctx, &pb.MonitorInfo{Items: []*pb.ConfigItem{}, EndTime: 100379842})
	if err != nil {
		log.Fatalf("Can't parse: %v", err)
	}

	log.Printf("Get message from get latest config: %v\n", ur)
	return nil
}


