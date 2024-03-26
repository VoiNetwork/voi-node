package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

func executeCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to execute %s: %v", command, err)
	}
}

func startProcess(command string, args ...string) chan bool {
	done := make(chan bool)

	go func() {
		executeCommand(command, args...)
		done <- true
	}()

	return done
}

func main() {
	done := startProcess("/node/bin/node_exporter",
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
		"--collector.textfile.directory=/algod/metrics",
		"--web.listen-address=:8080")

	// Wait for node_exporter to start
	time.Sleep(5 * time.Second)

	// Execute catch-catchpoint
	executeCommand("/node/bin/get-metrics", "-d", "/algod/metrics")

	<-done
}
