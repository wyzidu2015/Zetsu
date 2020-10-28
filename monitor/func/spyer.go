package monitor

import (
	"fmt"
	// "log"
	pb "Zetsu/zetsu"
)

type Spyer struct {
	Types []pb.ConfigItem_MonitorItemType
	Interval int32
}

func GetSpyerByMonitorConfig(basicInfo *pb.MachineBasicInfo) *Spyer {
	types := []pb.ConfigItem_MonitorItemType{pb.ConfigItem_CPU, pb.ConfigItem_MEM, pb.ConfigItem_NET}
	return &Spyer{Types: types, Interval: 10}
}

func (s *Spyer) Start() {
	fmt.Println("In Spyder start func")
}

func (s *Spyer) End() {
	fmt.Println("In Spyder end func")
}

