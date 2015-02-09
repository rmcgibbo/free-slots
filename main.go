// -*-  tab-width:4  -*-
package main

import (
	"fmt"
	"github.com/BurntSushi/ty/fun"
	"log"
	"os"
	"sge"
	"slurm"
	"text/tabwriter"
	"torque"
)

func main() {
	if slurm.IsSLURM() {
		mainSLURM()
	} else if torque.IsTorque() {
		mainTORQUE()
	} else if sge.IsSGE() {
		mainSGE()
	} else {
		log.Fatal("Didn't find any recognized scheduler")
	}
}

func mainSLURM() {
	fmt.Printf("Summary of SLURM nodes with free slots\n\n")
	counts := slurm.CollectFreeSlots()
	var _ = counts
	rows := []string{
		"Num Nodes\tPartition\tUtilization\tFree slots",
		"---------\t---------\t-----------\t----------"}

	keys := make([]slurm.NodeStatus, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}

	fun.Sort(func(a, b slurm.NodeStatus) bool {
		if a.Partition == b.Partition {
			return float64(a.NpAlloc)/float64(a.NpTotal) <
				float64(b.NpAlloc)/float64(b.NpTotal)
		}
		return a.Partition < b.Partition
	}, keys)

	for _, k := range keys {
		if k.NpTotal == k.NpAlloc {
			continue
		}
		rows = append(rows, fmt.Sprintf("%d\t%s\t%d/%d\t%d",
			counts[k], k.Partition, k.NpAlloc, k.NpTotal, k.NpTotal-k.NpAlloc))
	}

	printTable(rows)
}

func mainTORQUE() {
	fmt.Printf("Summary of PBS nodes with free slots\n\n")
	counts := torque.CollectFreeSlots()
	var _ = counts
	rows := []string{
		"Num Nodes\tProperties\tUtilization\tFree slots",
		"---------\t----------\t-----------\t----------"}

	keys := make([]torque.NodeStatus, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}

	fun.Sort(func(a, b torque.NodeStatus) bool {
		if a.Properties == b.Properties {
			return float64(a.NpAlloc)/float64(a.NpTotal) <
				float64(b.NpAlloc)/float64(b.NpTotal)
		}
		return a.Properties < b.Properties
	}, keys)

	for _, k := range keys {
		if k.NpTotal == k.NpAlloc {
			continue
		}
		rows = append(rows, fmt.Sprintf("%d\t%s\t%d/%d\t%d",
			counts[k], k.Properties, k.NpAlloc, k.NpTotal, k.NpTotal-k.NpAlloc))
	}
	printTable(rows)

}

func mainSGE() {
	fmt.Printf("Summary of SGE nodes with free slots\n\n")
	counts := sge.CollectFreeSlots()

	rows := []string{
		"Num Nodes\tQueue\tUtilization\tFree slots",
		"---------\t-----\t-----------\t----------"}

	keys := make([]sge.NodeStatus, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}

	fun.Sort(func(a, b sge.NodeStatus) bool {
		if a.Queue == b.Queue {
			return float64(a.NpAlloc)/float64(a.NpTotal) <
				float64(b.NpAlloc)/float64(b.NpTotal)
		}
		return a.Queue < b.Queue
	}, keys)

	for _, k := range keys {
		if k.NpTotal == k.NpAlloc {
			continue
		}
		rows = append(rows, fmt.Sprintf("%d\t%s\t%d/%d\t%d",
			counts[k], k.Queue, k.NpAlloc, k.NpTotal, k.NpTotal-k.NpAlloc))
	}

	printTable(rows)
}

func printTable(rows []string) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	for _, line := range rows {
		fmt.Fprintln(w, line)
	}
	w.Flush()
}
