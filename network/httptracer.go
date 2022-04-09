package network

import (
	"net/http/httptrace"
	"time"
)

type DNSInfo struct {
	DnsStartInfo    httptrace.DNSStartInfo
	DnsStartTime    time.Time
	DnsDoneInfo     httptrace.DNSDoneInfo
	DnsDoneDuration time.Duration
}

type ConnectionInfo struct {
	StartNetwork, StartAddr string
	EndNetwork, EndAddr     string
	ConnStartTime           time.Time
	ConnEndTime             time.Time
	Duration                time.Duration
	Err                     error
}

type FileMetadata struct {
	HostPort      string
	DNSInfo       DNSInfo
	ConnInfo      ConnectionInfo
	GotConnection httptrace.GotConnInfo
	RTT2FirstByte time.Duration
}

func GetTraceRequestFile(fileMetadata *FileMetadata, startTime *time.Time) *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			fileMetadata.HostPort = hostPort
		},
		GotFirstResponseByte: func() {
			fileMetadata.RTT2FirstByte = time.Since(*startTime)
		},
		GotConn: func(info httptrace.GotConnInfo) {
			fileMetadata.GotConnection = info
		},
		Got100Continue: nil,
		Got1xxResponse: nil,
		DNSStart: func(info httptrace.DNSStartInfo) {
			fileMetadata.DNSInfo.DnsStartInfo = info
			fileMetadata.DNSInfo.DnsStartTime = time.Now()
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			fileMetadata.DNSInfo.DnsDoneInfo = info
			fileMetadata.DNSInfo.DnsDoneDuration = time.Since(fileMetadata.DNSInfo.DnsStartTime)
		},
		ConnectStart: func(network, addr string) {
			fileMetadata.ConnInfo.ConnStartTime = time.Now()
			fileMetadata.ConnInfo.StartAddr = addr
			fileMetadata.ConnInfo.StartNetwork = network
		},
		ConnectDone: func(network, addr string, err error) {
			fileMetadata.ConnInfo.EndNetwork = network
			fileMetadata.ConnInfo.EndAddr = addr
			fileMetadata.ConnInfo.ConnEndTime = time.Now()
			fileMetadata.ConnInfo.Duration = time.Since(fileMetadata.ConnInfo.ConnStartTime)
		},
		TLSHandshakeStart: nil,
		TLSHandshakeDone:  nil,
		WroteHeaderField:  nil,
		WroteHeaders:      nil,
		Wait100Continue:   nil,
		WroteRequest:      nil,
	}
}
