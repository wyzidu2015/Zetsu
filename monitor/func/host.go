package monitor

import (
	"fmt"
	"net"
	"runtime"
	"syscall"
	pb "Zetsu/zetsu"
)

type HostInfo struct {
	IP string
	OS string
	CpuArch pb.MachineBasicInfo_CPUArch
	CpuCores int32
	MemorySize int32
}

func NewHostInfo() *HostInfo {
	host := HostInfo {
		IP: getIpAddress(),
		OS: getOSString(), 
		CpuArch: getCPUArch(), 
		CpuCores: getCPUCores(), 
		MemorySize: getMemorySize(),
	}
	
	return &host
}

func (h *HostInfo) ToMachineBasicInfo() *pb.MachineBasicInfo {
	return &pb.MachineBasicInfo{
		CpuArch: h.CpuArch,
		CpuCores: h.CpuCores,
		MemorySize: h.MemorySize,
		OsType: h.OS,
	}
}

func (h *HostInfo) ToConnectionInfo() *pb.MachineConnectInfo {
	return &pb.MachineConnectInfo{
		Uri:     fmt.Sprintf("%s:%d", h.IP, 8765),
		GroupId: 123,
	}
}

func (h *HostInfo) getMonitorInfo(configItems []*pb.ConfigItem) *pb.MonitorInfo {
	return &pb.MonitorInfo{
		Items:   configItems,
		EndTime: 0,
	}
}

func getIpAddress() string {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet
		isIpNet bool
		err     error
	)

	if addrs, err = net.InterfaceAddrs(); err != nil {
		return ""
	}

	ipv4 := ""
	for _, addr = range addrs {
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()
				return ipv4
			}
		}
	}
	return ""
}

func getOSString() string {
	return runtime.GOOS
}

func getCPUArch() pb.MachineBasicInfo_CPUArch {
	switch runtime.GOARCH {
	case "amd64":
		return pb.MachineBasicInfo_X86_64
	case "arm64":
		return pb.MachineBasicInfo_AARCH64
	default:
		return pb.MachineBasicInfo_X86
	}
}

func getCPUCores() int32 {
	return int32(runtime.NumCPU())
}

func getMemorySize() int32 {
	memStat := new(runtime.MemStats)
	runtime.ReadMemStats(memStat)

	memSize := memStat.Sys * uint64(syscall.Getpagesize())
	memSize = memSize / 1024 / 1024 / 1024 / 8
	
	return int32(memSize)
}