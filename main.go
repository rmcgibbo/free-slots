package main

import (
	// "free_slots"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	// "os/exec"
)

type Status struct {
	Queue    string
	np_alloc int
	np_total int
}

func main() {
	fmt.Printf("hello, world\n")

	counts := Collect()
	fmt.Println(counts)
	ParseQconfMap()
}
func ParseQconfList() []string {
}

func ParseQconfMap() map[string]string {
	file, err := os.Open("qconf.sq.batch.q.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	inBytes, _ := ioutil.ReadAll(file)
	lines := strings.Split(string(inBytes), "\n")

	var withinMultilineEntry = false
	var key = ""
	var value = ""

	d := make(map[string]string)

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
			d[key] = value
		}

	}
	fmt.Println(d)
	return d
}

func QueuePermissions() {
    queues = parse_qconf('-sql')

}

func Collect() map[Status]int {
	type Item struct {
		Data  string `xml:",chardata"`
		Name  string `xml:"name,attr"`
		QName string `xml:"qname,attr"`
	}
	type Host struct {
		Name       string `xml:"name,attr"`
		Hostvalues []Item `xml:"hostvalue"`
		Queuevalue []Item `xml:"queue>queuevalue"`
	}
	type QHost struct {
		XMLName xml.Name `xml:"qhost"`
		Hosts   []Host   `xml:"host"`
	}

	file, err := os.Open("qhost_q_xml.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	xmlData, _ := ioutil.ReadAll(file)

	var q QHost
	xml.Unmarshal(xmlData, &q)

	counts := make(map[Status]int)

	for _, h := range q.Hosts {
		var slots_used = -1
		var slots = -1
		for _, v := range h.Queuevalue {
			if v.Name == "slots_used" {
				slots_used, _ = strconv.Atoi(v.Data)
			} else if v.Name == "slots" {
				slots, _ = strconv.Atoi(v.Data)
			}
		}
		if slots_used != -1 && slots != -1 {
			s := Status{h.Queuevalue[0].QName, slots_used, slots}
			_, ok := counts[s]
			if ok {
				counts[s]++
			} else {
				counts[s] = 1
			}
		}
	}
	return counts
}
