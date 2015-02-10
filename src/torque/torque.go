// -*-  tab-width:4  -*-
package torque

import (
	"encoding/xml"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

// Does this system use Torque?
func IsTorque() bool {
	_, err := exec.LookPath("pbsnodes")
	return err == nil
}

type NodeStatus struct {
	Properties string
	NpAlloc    int
	NpTotal    int
}

func CollectFreeSlots() map[NodeStatus]int {
	type Node struct {
		State      string `xml:"state"`
		Properties string `xml:"properties"`
		NpTotal    string `xml:"np"`
		Jobs       string `xml:"jobs"`
	}
	type Data struct {
		XMLName xml.Name `xml:"Data"`
		Nodes   []Node   `xml:"Node"`
	}

	xmlData, err := exec.Command("pbsnodes", "-x").Output()
	if err != nil {
		log.Fatal(err)
	}
	var root Data
	xml.Unmarshal(xmlData, &root)
	counts := make(map[NodeStatus]int)

	for _, n := range root.Nodes {
		if n.State != "free" {
			continue
		}

		var NpAlloc = 0
		if n.Jobs != "" {
			NpAlloc = len(strings.Split(n.Jobs, ", "))
		}
		NpTotal, _ := strconv.Atoi(n.NpTotal)

		ns := NodeStatus{n.Properties, NpAlloc, NpTotal}
		_, ok := counts[ns]
		if ok {
			counts[ns]++
		} else {
			counts[ns] = 1
		}
	}

	return counts
}
