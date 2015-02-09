// -*-  tab-width:4  -*-
package main

import (
	"fmt"
	"github.com/BurntSushi/ty/fun"
	"os"
	"sge"
	"text/tabwriter"
)

func main() {
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
	PrintTable(rows)
}

func PrintTable(rows []string) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	for _, line := range rows {
		fmt.Fprintln(w, line)
	}
	w.Flush()
}
