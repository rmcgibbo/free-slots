// -*-  tab-width:4  -*-
package slurm

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
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
        for _, node := range expandBracket(part["Nodes"]) {
            _, ok := n2plist[node]
            if ok {
                n2plist[node] = append(n2plist[node], part["PartitionName"])
            } else {
                n2plist[node] = []string{part["PartitionName"]}
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
	data, _ := exec.Command("scontrol", "show", "hostnames",s).Output()
	values := make([]string, 0)
	for _, line := range strings.Split(string(data), "\n") {
        values = append(values, line);
	}
	return values
}
