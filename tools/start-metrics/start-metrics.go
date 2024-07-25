package main

import (
	"github.com/voinetwork/voi-node/tools/utils"
	"log"
	"time"
)

const (
	nodeExporterCmd            = "/node/bin/node_exporter"
	getMetricsCmd              = "/node/bin/get-metrics"
	metricsDir                 = "/algod/metrics"
	nodeExporterListenAddr     = "0.0.0.0:8080"
	nodeExporterStartupTimeout = 30 * time.Second
)

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

	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
	case <-time.After(5 * time.Second):
		return nil
	}
	return nil
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

	log.Printf("Node exporter started successfully")
	log.Printf("Starting get-metrics")

	if err := executeGetMetrics(pu); err != nil {
		log.Fatalf("Error running get-metrics: %v", err)
	}
}
