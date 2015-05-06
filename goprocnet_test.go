package goprocnet

import (
	"testing"
)

func TestTypeToFileMap(t *testing.T) {
	result := getFilename("tcp")
	if result != "/proc/net/tcp" {
		t.Fail()
	}
	result = getFilename("udp")
	if result != "/proc/net/udp" {
		t.Fail()
	}
	result = getFilename("foobar")
	if result != "" {
		t.Fail()
	}
}

func TestParseIPPort(t *testing.T) {
	ip := getIP("00000000")
	port := getPort("006F")
	if ip != "0.0.0.0" {
		t.Fail()
	}
	if port != "111" {
		t.Fail()
	}
	ip = getIP("0B01010A")
	if ip != "10.1.1.11" {
		t.Fail()
	}
	port = getPort("CAEF")
	if port != "51951" {
		t.Fail()
	}
}
