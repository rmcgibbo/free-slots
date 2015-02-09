// -*-  tab-width:4  -*-
package slurm

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	//	"reflect"
)

var _ = fmt.Println

func IsSLURM() bool {
	_, err := exec.LookPath("scontrol")
	return err == nil
}

type NodeStatus struct {
	Partition string
	NpAlloc   int
	NpTotal   int
}

func CollectFreeSlots() map[NodeStatus]int {
	n2p := nodeToPartition()
	nodes := scontrolShow("nodes")
	counts := make(map[NodeStatus]int)

	for _, n := range nodes {
		CPUAlloc, _ := strconv.Atoi(n["CPUAlloc"])
		CPUTot, _ := strconv.Atoi(n["CPUTot"])
		name := n2p[n["NodeHostName"]]
		if name == "" {
			continue
		}

		s := NodeStatus{n2p[n["NodeHostName"]], CPUAlloc, CPUTot}
		_, ok := counts[s]
		if ok {
			counts[s] += 1
		} else {
			counts[s] = 1
		}
	}
	return counts
}

// Mapping from NodeHostName to PartitionName for each
// node
func nodeToPartition() map[string]string {
	partitions := scontrolShow("partition")
	n2plist := make(map[string][]string)

	for _, part := range partitions {
		for _, nodeGroup := range strings.Split(part["Nodes"], ",") {
			for _, node := range expandBracket(nodeGroup) {
				_, ok := n2plist[node]
				if ok {
					n2plist[node] = append(n2plist[node], part["PartitionName"])
				} else {
					n2plist[node] = []string{part["PartitionName"]}
				}
			}
		}
	}

	out := make(map[string]string)
	for k, v := range n2plist {
		out[k] = strings.Join(v, ",")
	}

	return out
}

func scontrolShow(cmd string) []map[string]string {
	data, _ := exec.Command("scontrol", "show", "-o", cmd).Output()
	values := make([]map[string]string, 0)

	for _, line := range strings.Split(string(data), "\n") {
		if len(line) == 0 {
			continue
		}

		fields := strings.Fields(line)
		row := make(map[string]string)
		for _, field := range fields {
			items := strings.Split(field, "=")
			if len(items) == 2 {
				row[items[0]] = items[1]
			}
		}
		values = append(values, row)
	}
	return values
}

func expandBracket(s string) []string {
	m := regexp.MustCompile(`(.*)\[(\d+)\-(\d+)(?:,(\d+)\-(\d+))*\]`)
	groups := m.FindStringSubmatch(s)
	out := make([]string, 0)
	if len(groups) == 0 {
		return out
	}

	// for i, g := range groups {
	// 	fmt.Println(i, g, len(g))
	// }

	prefix := groups[1]
	for i := 2; i < len(groups); i += 2 {
		//	 	leading0 := (groups[i][0:1] == "0")
		if len(groups[i]) > 0 && len(groups[i+1]) > 0 {
			first, _ := strconv.Atoi(groups[i])
			last, _ := strconv.Atoi(groups[i+1])
			for j := first; j < last+1; j++ {
				suffix := strconv.Itoa(j)
				out = append(out, prefix+suffix)
			}
		}
	}

	//fmt.Println(out)
	return out
}
