package parser

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type Parser struct {
	Sessions map[uint64]Session
}

type Session struct {
	Packets []gopacket.Packet
}

func (s *Session) addPacket(p gopacket.Packet) {
	s.Packets = append(s.Packets, p)
}

func (s *Session) Summary() {
	for _, p := range s.Packets {
		fmt.Println(p.Dump())
	}
}

func (p *Parser) createSession(hash uint64) Session {
	p.Sessions[hash] = Session{}
	return p.Sessions[hash]
}

func (p *Parser) Parse(path string) {
	var (
		ok      bool
		current Session
	)
	if handle, err := pcap.OpenOffline(path); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			if networkLayer := packet.NetworkLayer(); networkLayer != nil {
				hash := networkLayer.NetworkFlow().FastHash()
				// New session
				if current, ok = p.Sessions[hash]; !ok {
					current = p.createSession(hash)
				}
				// Session already exists
				current.Packets = append(current.Packets, packet)
				p.Sessions[hash] = current
			}
		}
	}
}

func NewParser() *Parser {
	return &Parser{
		Sessions: make(map[uint64]Session),
	}
}
