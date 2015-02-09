// -*-  tab-width:4  -*-
package sge

import (
	"encoding/xml"
	"log"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

func IsSGE() bool {
	_, err := exec.LookPath("qhost")
	return err == nil
}

type NodeStatus struct {
	Queue   string
	NpAlloc int
	NpTotal int
}

func CollectFreeSlots() map[NodeStatus]int {
	type Item struct {
		Data  string `xml:",chardata"`
		Name  string `xml:"name,attr"`
		QName string `xml:"qname,attr"`
	}
	type Queue struct {
		Name  string `xml:"name,attr"`
		Items []Item `xml:"queuevalue"`
	}
	type Host struct {
		Name       string  `xml:"name,attr"`
		Hostvalues []Item  `xml:"hostvalue"`
		Queues     []Queue `xml:"queue"`
	}
	type QHost struct {
		XMLName xml.Name `xml:"qhost"`
		Hosts   []Host   `xml:"host"`
	}

	xmlData, _ := exec.Command("qhost", "-q", "-xml").Output()
	var root QHost
	xml.Unmarshal(xmlData, &root)

	counts := make(map[NodeStatus]int)
	allowed := queuePermissions()

	for _, h := range root.Hosts {
		for _, q := range h.Queues {
			var slots_used = -1
			var slots = -1

			if !allowed[q.Name] {
				continue
			}

			for _, v := range q.Items {
				if v.Name == "slots_used" {
					slots_used, _ = strconv.Atoi(v.Data)
				} else if v.Name == "slots" {
					slots, _ = strconv.Atoi(v.Data)
				}
			}
			if slots_used != -1 && slots != -1 {
				s := NodeStatus{q.Name, slots_used, slots}
				_, ok := counts[s]
				if ok {
					counts[s]++
				} else {
					counts[s] = 1
				}
			}
		}
	}
	return counts
}

func qconflines(cmd ...string) []string {
	out, err := exec.Command("qconf", cmd...).Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(out), "\n")
}

func userPermittedInQueue(queue, user string) bool {
	u, ok := parseQconf("-sq", queue)["user_lists"]
	if !ok || u == "NONE" {
		return true
	}
	panic("Not Implemented")
}

func queuePermissions() map[string]bool {
	user, _ := user.Current()
	allowed := make(map[string]bool)
	for _, q := range qconflines("-sql") {
		allowed[q] = userPermittedInQueue(q, user.Username)
	}
	return allowed
}

func parseQconf(cmd ...string) map[string]string {
	lines := qconflines(cmd...)
	props := make(map[string]string)
	var withinMultilineEntry = false
	var key = ""
	var value = ""

	for _, line := range lines {
		if withinMultilineEntry {
			value += strings.TrimSpace(strings.TrimSuffix(line, "\\"))
			withinMultilineEntry = strings.HasSuffix(line, "\\")

		} else {
			fields := strings.SplitN(strings.TrimSpace(line), " ", 2)
			if len(fields) == 2 {
				key = strings.TrimSpace(fields[0])
				withinMultilineEntry = strings.HasSuffix(line, "\\")
				value = strings.TrimSpace(strings.TrimSuffix(fields[1], "\\"))

			}
		}

		if !withinMultilineEntry {
			props[key] = value
		}

	}
	return props
}
