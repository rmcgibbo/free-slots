// -*-  tab-width:4  -*-
package main

import (
	"fmt"
	"github.com/BurntSushi/ty/fun"
	"os"
	"sge"
	"text/tabwriter"
	"torque"
)

func main() {
	if torque.IsTorque() {
		mainTORQUE()
	}
}

func mainTORQUE() {
	fmt.Printf("Summary of PBS nodes with free slots\n\n")
	counts := torque.CollectFreeSlots()
	var _ = counts
	rows := []string{
		"Number of Nodes  \tProperties\tUtilization\tFree slots",
		"---------------  \t----------\t-----------\t----------"}

	keys := make([]torque.NodeStatus, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	fun.Sort(func(a, b torque.NodeStatus) bool {
		return float64(a.NpAlloc)/float64(a.NpTotal) <
			float64(b.NpAlloc)/float64(b.NpTotal)
	}, keys)
	for _, k := range keys {
		rows = append(rows, fmt.Sprintf("%d\t%s\t%d/%d\t%d",
			counts[k], k.Properties, k.NpAlloc, k.NpTotal, k.NpTotal-k.NpAlloc))
	}
	printTable(rows)

}

func mainSGE() {
	fmt.Printf("Summary of SGE nodes with free slots\n\n")
	counts := sge.CollectFreeSlots()

	rows := []string{
		"Number of Nodes\tQueue\tUtilization\tFree slots",
		"---------------\t-----\t-----------\t----------"}

	keys := make([]sge.NodeStatus, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	fun.Sort(func(a, b sge.NodeStatus) bool {
		if a.Queue == b.Queue {
			return (a.NpTotal - a.NpAlloc) > (b.NpTotal - b.NpAlloc)
		}
		return a.Queue < b.Queue
	}, keys)
	for _, k := range keys {
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
