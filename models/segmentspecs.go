package models

import "github.com/massimo-gollo/DASHpher/network"

// https://streaminglearningcenter.com/blogs/itu-t-p1203-p1204.html info about P1203

//SegmentInfo store information about Segment Obtained
type SegmentInfo struct {
	NetDetails  network.FileMetadata
	ArrivalTime int
	//deliveriTime of segmentRequested
	DeliveryTime int
	StallTime    int
	Bandwidth    int
	DelRate      int
	ActRate      int
	SegmentSize  int
	//P1203HeaderSize float64
	// buffer = difference in arr_times for adjacent segments + segment duration of this segment
	BufferLevel       int
	Adapt             string
	SegmentDuration   int
	ReprCodec         string
	ReprWidth         int
	ReprHeight        int
	ReprFps           int
	PlayStartPosition int
	PlaybackTime      int
	Rtt               float64
	ReprIndex         int
	//MpdIndex             int
	AdaptIndex      int
	SegmentIndex    int
	Played          bool
	SegReplace      string
	P1203           float64
	HTTPprotocol    string
	Clae            float64
	Duanmu          float64
	Yin             float64
	Yu              float64
	P1203Kbps       float64
	SegmentFileName string

	// QoE metrics
	SegmentRates   []float64
	SumSegRate     float64
	TotalStallDur  float64
	NumStalls      int
	NumSwitches    int
	RateDifference float64
	SumRateChange  float64
	RateChange     []float64
	MimeType       string
	Profile        string
}

func NewSegmentInfo() *SegmentInfo {
	return &SegmentInfo{NetDetails: network.FileMetadata{}}
}
