package models

import (
	"github.com/massimo-gollo/DASHpher/network"
	"time"
)

type ReproductionMetrics struct {
	Id           uint64
	Url          string
	FetchMpdInfo network.FileMetadata
	MPD          MPD
	SegmentsInfo map[int]SegmentInfo

	ReprStartTime time.Time
	ReprEndTime   time.Time
	ReprDuration  time.Duration

	CompleteWithSuccess bool
	CompletedWithError  bool
	ErrorCount          int
	SegmentErrorCount   int
	LastErrorReason     string
}

func NewReproductionMetrics() *ReproductionMetrics {
	return &ReproductionMetrics{SegmentsInfo: make(map[int]SegmentInfo)}
}
