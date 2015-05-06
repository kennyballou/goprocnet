package goprocnet

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	PROC_NET_TCP = "/proc/net/tcp"
	PROC_NET_UDP = "/proc/net/udp"
)

var STATE = map[string]string{
	"01": "ESTABLISHED",
	"02": "SYN_SENT",
	"03": "SYN_RECV",
	"04": "FIN_WAIT1",
	"05": "FIN_WAIT2",
	"06": "TIME_WAIT",
	"07": "CLOSE",
	"08": "CLOSE_WAIT",
	"09": "LAST_ACK",
	"0A": "LISTEN",
	"0B": "CLOSING",
}

var Version string

type Socket struct {
	PID        string
	UID        string
	Bin        string
	Name       string
	State      string
	LocalIP    string
	LocalPort  string
	RemoteIP   string
	RemotePort string
}

func (s Socket) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s:%s\t%s:%s",
		s.PID,
		s.UID,
		s.Bin,
		s.State,
		s.LocalIP,
		s.LocalPort,
		s.RemoteIP,
		s.RemotePort)
}

func GetTCPSockets() []Socket {
	sockets, err := netstat("tcp")
	if err != nil {
		log.Panic(err)
	}
	return sockets
}

func GetUDPSockets() []Socket {
	sockets, err := netstat("udp")
	if err != nil {
		log.Panic(err)
	}
	return sockets
}

func getFilename(t string) string {
	switch {
	case t == "tcp":
		return PROC_NET_TCP
	case t == "udp":
		return PROC_NET_UDP
	}
	return ""
}

func readFile(t string) ([]string, error) {
	filename := getFilename(t)
	if filename == "" {
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Unable to read file", err)
		return nil, err
	}
	lines := strings.Split(string(data), "\n")

	// return lines without header or tailing blank line
	return lines[1 : len(lines)-1], nil
}

func convertHexToDec(hex string) int64 {
	d, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		log.Panic("Unable to convert hex to dec", err)
	}
	return d
}

func getIP(hexip string) string {
	// we are going to assume IPv4 addresses for now...
	// We are also assuming address is in little endian
	return fmt.Sprintf("%v.%v.%v.%v",
		convertHexToDec(hexip[6:8]),
		convertHexToDec(hexip[4:6]),
		convertHexToDec(hexip[2:4]),
		convertHexToDec(hexip[0:2]))
}

func getPort(hex string) string {
	return fmt.Sprintf("%v", convertHexToDec(hex))
}

func removeInnerSpace(values []string) []string {
	var newValues []string
	for _, value := range values {
		if value != "" {
			newValues = append(newValues, value)
		}
	}
	return newValues
}

func netstat(t string) ([]Socket, error) {
	lines, err := readFile(t)
	if err != nil {
		log.Println("Failed to retrieve proc data", err)
		return nil, err
	}

	var sockets []Socket

	for _, line := range lines {
		values := removeInnerSpace(strings.Split(strings.TrimSpace(line), " "))

		ipPort := strings.Split(values[1], ":")
		localIP := getIP(ipPort[0])
		localPort := getPort(ipPort[1])

		remIpPort := strings.Split(values[2], ":")
		remoteIP := getIP(remIpPort[0])
		remotePort := getPort(remIpPort[1])

		state := STATE[values[3]]
		uid := values[7]
		pid := getProcessID(values[9])
		exe := getProcessExecutable(pid)
		name := getProcessName(exe)
		s := Socket{pid, uid, exe, name, state, localIP, localPort, remoteIP, remotePort}
		sockets = append(sockets, s)
	}
	return sockets, nil
}

func getProcessID(inode string) string {
	// Loop over the fd directories under /proc/pid for the right inode
	pid := "-"

	procDirs, err := filepath.Glob("/proc/[0-9]*/fd/[0-9]*")
	if err != nil {
		log.Panic(err)
	}

	re := regexp.MustCompile(inode)
	for _, item := range procDirs {
		path, _ := os.Readlink(item)
		out := re.FindString(path)
		if len(out) != 0 {
			pid = strings.Split(item, "/")[2]
		}
	}
	return pid
}

func getProcessExecutable(pid string) string {
	path, _ := os.Readlink(fmt.Sprintf("/proc/%s/exe", pid))
	return path
}

func getProcessName(exe string) string {
	n := strings.Split(exe, "/")
	name := n[len(n)-1]
	return strings.Title(name)
}
