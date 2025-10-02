package main

import (
	"encoding/hex"
	"fmt"
	"math/rand/v2"
	"strings"
)

// GirPayloadUDP: Payload-haye pishrafte baraye UDP
func GirPayloadUDP(noepayload string) []byte {
	switch strings.ToLower(noepayload) {
	case "dns":
		// DNS Query pishrafte ba ID tasadofi baraye A record example.com
		id := rand.Uint16()
		header := fmt.Sprintf("%04x01000001000000000000076578616d706c6503636f6d0000010001", id)
		data, _ := hex.DecodeString(header)
		return data
	case "ntp":
		// NTP Monlist Request pishrafte baraye estekhraj noskhe va info
		data, _ := hex.DecodeString("17 00 03 2a 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00")
		return data
	case "snmp":
		// SNMP GetBulkRequest pishrafte baraye sysDescr va sysObjectID
		data, _ := hex.DecodeString("30 2e 02 01 00 04 06 70 75 62 6c 69 63 a5 21 02 04 00 00 00 01 02 01 00 02 01 00 30 13 30 11 06 0d 2b 06 01 02 01 01 03 00 05 00 30 0f 06 0b 2b 06 01 02 01 01 01 00 05 00")
		return data
	default:
		return []byte{0x00} // Minimal probe for no payload
	}
}
