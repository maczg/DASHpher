package reproduction

import (
	"github.com/massimo-gollo/DASHpher/models"
	"time"
)

type StreamStruct struct {
	//Original MPD info plus copy of mpd
	OriginalStreamDuration  int
	OriginalTotalSegmentMPD int
	OriginalUrl             string
	//Unit duration of segment in seconds
	OriginalDurationPerSegment   int
	MaxHeightReprIdx             int
	MinHeightReprIdx             int
	BandwidthList                []int
	Profile                      string
	MPD                          models.MPD
	Codec                        string
	IsByteRangeMPD               bool
	StartTimeReproduction        *time.Time
	EndTimeReproduction          *time.Time
	ReproductionCompleteDuration time.Duration

	//Fine-tuned reproduction info
	ActualStreamDuration       int
	ActualTotalSegmentToStream int
	MaxReqHeight               int
	InitBuffer                 int
	MaxBuffer                  int
	AdaptionAlgorithm          string

	//Current parameters of streaming
	CurrentURLSegToStream        string
	CurrentSegmentInReproduction int
	CurrentHeightReprIdx         int
	CurrentBandwidth             int
	CurrentSegSize               int
	NextSegmentNumber            int
	ThroughputList               []int

	BufferLevel    int
	MapSegmentInfo map[int]*models.SegmentInfo

	MimeTypes []int

	//	NextRunTime           time.Time
	//ArrivalTime           int

	SegmentDurationTotal int
	BaseURL              string
	AudioContent         bool
	RepRate              int
}
