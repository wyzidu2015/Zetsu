package monitor

import (
	pb "Zetsu/zetsu"
	"context"
	"errors"
	"fmt"
	cron "github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"log"
	"time"
)

type Monitor struct {
	Client pb.ZetsuClient
	Ctx context.Context
	Cancel context.CancelFunc
	crontab *cron.Cron
	maxSaveTime int32
	configItems []*pb.ConfigItem
	spec string
}

func NewMonitor(address string) *Monitor {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Can't connect: %v", err)
	}

	cli := pb.NewZetsuClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	monitor := Monitor{Client: cli, Ctx: ctx, Cancel: cancel}
	hostInfo := NewHostInfo()
	if err := monitor.register(hostInfo); err != nil {
		log.Fatalf("Failed to register to remote: %s", err.Error())
	}
	if err := monitor.getConfig(hostInfo); err != nil {
		log.Fatalf("Failed to register to remote: %s", err.Error())
	}

	crontab := cron.New()
	crontab.AddFunc(monitor.spec, func() {
		hostInfo := NewHostInfo()
		err := monitor.uploadInfo(hostInfo)
		if err != nil {
			log.Printf("Failed to upload: %s", err.Error())
		}
	})
	monitor.crontab = crontab

	return &monitor
}

func (mon *Monitor) Start() {
	mon.crontab.Start()
}

func (mon *Monitor) Stop()  {
	mon.crontab.Stop()
}

func (mon *Monitor) register(host *HostInfo) error {
	r, err := mon.Client.RegisterMonitor(mon.Ctx, host.ToConnectionInfo())
	if err != nil {
		return err
	}

	if r.Status != 0 {
		return errors.New(r.Info)
	}
	return nil
}

func (mon *Monitor) getConfig(host *HostInfo) error {
	mr, err := mon.Client.GetLatestConfig(mon.Ctx, host.ToMachineBasicInfo())
	if err != nil {
		return err
	}

	mon.spec = fmt.Sprintf("*/%d * * * * ?", mr.Interval)
	mon.maxSaveTime = mr.MaxSaveTime
	mon.configItems = mr.Items
	return nil
}

func (mon *Monitor) uploadInfo(host *HostInfo) error {
	ur, err := mon.Client.UploadMonitorItem(mon.Ctx, host.getMonitorInfo(mon.configItems))
	if err != nil {
		return err
	}

	if ur.Status != 0 {
		return errors.New(ur.Info)
	}

	return nil
}


