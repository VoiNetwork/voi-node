package main

import (
	"fmt"
	"github.com/voinetwork/voi-node/tools/utils"
	"log"
	"net"
	"time"
)

const (
	nodeExporterCmd            = "/node/bin/node_exporter"
	getMetricsCmd              = "/node/bin/get-metrics"
	metricsDir                 = "/algod/metrics"
	nodeExporterListenAddr     = "0.0.0.0:8080"
	nodeExporterStartTimeout   = 5 * time.Second
	nodeExporterStartupTimeout = 30 * time.Second
)

func isPortOpen(address string) bool {
	conn, err := net.DialTimeout("tcp", address, nodeExporterStartTimeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func startNodeExporter(pu utils.ProcessUtils) error {
	errChan := pu.StartProcess(nodeExporterCmd,
		"--no-collector.arp",
		"--no-collector.bcache",
		"--no-collector.bonding",
		"--no-collector.buddyinfo",
		"--no-collector.conntrack",
		"--no-collector.drbd",
		"--no-collector.edac",
		"--no-collector.entropy",
		"--no-collector.hwmon",
		"--no-collector.infiniband",
		"--no-collector.interrupts",
		"--no-collector.ipvs",
		"--no-collector.ksmd",
		"--no-collector.logind",
		"--no-collector.mdadm",
		"--no-collector.meminfo_numa",
		"--no-collector.mountstats",
		"--no-collector.nfs",
		"--no-collector.nfsd",
		"--no-collector.qdisc",
		"--no-collector.runit",
		"--no-collector.supervisord",
		"--no-collector.systemd",
		"--no-collector.tcpstat",
		"--no-collector.timex",
		"--no-collector.wifi",
		"--no-collector.xfs",
		"--no-collector.zfs",
		"--collector.textfile.directory="+metricsDir,
		"--web.listen-address="+nodeExporterListenAddr)

	// Wait for process to start or fail
	err := <-errChan
	if err != nil {
		return err
	}

	timeout := time.After(nodeExporterStartupTimeout)
	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for node_exporter to start")
		default:
			if isPortOpen(nodeExporterListenAddr) {
				return nil
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func executeGetMetrics(pu utils.ProcessUtils) error {
	_, err := pu.ExecuteCommand(getMetricsCmd, "-d", metricsDir)
	return err
}

func main() {
	pu := utils.ProcessUtils{}

	if err := startNodeExporter(pu); err != nil {
		log.Fatalf("Error starting node_exporter: %v", err)
	}

	if err := executeGetMetrics(pu); err != nil {
		log.Fatalf("Error running get-metrics: %v", err)
	}
}
