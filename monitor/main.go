package main

// import (
// 	"context"
// 	"log"
// 	"time"
// 	pb "Zetsu/zetsu"
// 	"google.golang.org/grpc"
// )

import (
	mon "Zetsu/monitor/func"
)

const (
	address = "localhost:5790"
)

func main() {
	// get current spyer
	hostCollector := mon.NewHostCollector()
	spyder := mon.GetSpyerByMonitorConfig(hostCollector.ToMachineBasicInfo())
	spyder.Start()
	defer spyder.End()

	monitor := mon.NewMonitor(address)

	// register to server
	monitor.Register(hostCollector)

	// get config
	monitor.GetConfig(hostCollector)

	// upload info
	monitor.UploadInfo()


	// defer conn.Close()

	// cli := pb.NewZetsuClient(conn)

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()

}