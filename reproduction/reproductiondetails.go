package reproduction

import (
	"github.com/massimo-gollo/DASHpher/models"
	"time"
)

type StreamInfo struct {
	//Original MPD info plus copy of mpd
	TotalVideoDuration int
	TotalSegmentCount  int
	UrlResource        string
	//Unit duration of segment in seconds
	SingleSegmentDuration        int
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

	StartTimeDownloading   time.Time
	NextRunTimeDownloading time.Time
	WaitToPlayCount        int

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
	CurrentRepIdx                int
	CurrentBandwidth             int
	CurrentSegSize               int
	NextSegmentNumber            int
	ThroughputList               []int

	BufferLevel int
	//we consider only one adp set. If we want consider also audio or other adset
	// should take into account to change in []map[int]...
	SegmentInformation map[int]*models.SegmentInfo
	MimeTypes          []int

	//	NextRunTime           time.Time
	//ArrivalTime           int

	SegmentDurationTotal int
	BaseURL              string
	AudioContent         bool
	RepRate              int
}
