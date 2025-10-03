package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// GirPayloadUDP: Payload-haye pishrafte baraye UDP
func GirPayloadUDP(noepayload string) []byte {
	switch strings.ToLower(noepayload) {
	case "dns":
		// DNS Query pishrafte ba ID tasadofi baraye A record example.com
		rand.Seed(time.Now().UnixNano())
		id := uint16(rand.Intn(65536))
		header := fmt.Sprintf("%04x01000001000000000000076578616d706c6503636f6d0000010001", id)
		data, _ := hex.DecodeString(header)
		return data
	case "version.bind":
		// DNS Query baraye version.bind (TXT record, CHAOS class)
		rand.Seed(time.Now().UnixNano())
		id := uint16(rand.Intn(65536))
		header := fmt.Sprintf("%04x010000010000000000000776657273696f6e0462696e640000100003", id)
		data, _ := hex.DecodeString(header)
		return data
	case "ntp":
		// NTP Request sadeh-tar baraye estekhraj noskhe
		data, _ := hex.DecodeString("1b00000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		return data
	case "snmp":
		// SNMP GetRequest sadeh-tar
		data, _ := hex.DecodeString("302902010004067075626c6963a01e0204000000010201000201003010300e060a2b060102010101000500")
		return data
	case "test":
		// Payload sadeh baraye test
		return []byte("TEST")
	default:
		return []byte{0x00} // Minimal probe
	}
}
