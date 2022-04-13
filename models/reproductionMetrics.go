package models

import (
	"fmt"
	"github.com/massimo-gollo/DASHpher/network"
	"time"
)

type RepStatus int

func (r RepStatus) String() string {
	switch r {
	case Success:
		return "success"
	case Error:
		return "error"
	case Aborted:
		return "aborted"
	default:
		return fmt.Sprintf("%d", int(r))
	}
}

const (
	Success RepStatus = 0
	Error             = 1
	Aborted           = 2
)

type ReproductionMetrics struct {
	ReproductionID uint64
	ContentUrl     string
	FetchMpdInfo   network.FileMetadata
	MPD            MPD
	SegmentsInfo   map[int]SegmentInfo

	ReprStartTime time.Time
	ReprEndTime   time.Time
	ReprDuration  time.Duration

	Status            RepStatus
	ErrorCount        int
	StallCount        int
	SegmentErrorCount int
	LastErrorReason   string
}

func NewReproductionMetrics() *ReproductionMetrics {
	return &ReproductionMetrics{SegmentsInfo: make(map[int]SegmentInfo)}
}
