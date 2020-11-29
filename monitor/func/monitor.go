package monitor

import (
	"Zetsu/common"
	pb "Zetsu/zetsu"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

type Monitor struct {
	Client      pb.ZetsuClient
	Ctx         context.Context
	Cancel      context.CancelFunc
	timeTicker  *time.Ticker
	maxSaveTime int32
	configItems []*pb.ConfigItem
	spec        string
}

func NewMonitor(address string) *Monitor {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Can't connect: %v", err)
	}

	cli := pb.NewZetsuClient(conn)
	ctx, cancel := context.WithCancel(context.Background())

	monitor := Monitor{Client: cli, Ctx: ctx, Cancel: cancel}
	hostInfo := NewHostInfo()
	if err := monitor.register(hostInfo); err != nil {
		log.Fatalf("Failed to register to remote: %s", err.Error())
	}
	if err := monitor.getConfig(hostInfo); err != nil {
		log.Fatalf("Failed to get config from remote: %s", err.Error())
	}

	uploadFunc := func() {
		hostInfo := NewHostInfo()
		err := monitor.uploadInfo(hostInfo)
		if err != nil {
			log.Printf("Failed to upload: %s", err.Error())
		}
	}
	uploadFunc()

	timeTicker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			select {
			case <-timeTicker.C:
				uploadFunc()
			}
		}
	}()
	monitor.timeTicker = timeTicker

	return &monitor
}

func (mon *Monitor) Stop() {
	mon.timeTicker.Stop()
}

func (mon *Monitor) register(host *HostInfo) error {
	r, err := mon.Client.RegisterMonitor(mon.Ctx, host.ToConnectionInfo())
	if err != nil {
		return err
	}

	if r.Status != int32(common.SUCCESS) {
		return errors.New(r.Info)
	}
	return nil
}

func (mon *Monitor) getConfig(host *HostInfo) error {
	mr, err := mon.Client.GetLatestConfig(mon.Ctx, host.ToMachineBasicInfo())
	if err != nil {
		return err
	}

	mon.spec = fmt.Sprintf("*/%d * * * * *", mr.Interval)
	mon.maxSaveTime = mr.MaxSaveTime
	mon.configItems = mr.Items
	return nil
}

func (mon *Monitor) uploadInfo(host *HostInfo) error {
	ur, err := mon.Client.UploadMonitorItem(mon.Ctx, host.getMonitorInfo(mon.configItems))
	if err != nil {
		return err
	}

	if ur.Status != int32(common.SUCCESS) {
		return errors.New(ur.Info)
	}

	return nil
}
