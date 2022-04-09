package models

import (
	"encoding/xml"
	"github.com/sirupsen/logrus"
	"os"
)

var logger = logrus.Logger{
	Out:       os.Stderr,
	Formatter: &logrus.TextFormatter{DisableColors: false, TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true},
	Level:     logrus.InfoLevel,
}

// MPD structure
type MPD struct {
	XMLName xml.Name `xml:"MPD"`

	Xmlns                     string `xml:"xmlns,attr"`
	MinBufferTime             string `xml:"minBufferTime,attr"`
	MediaPresentationDuration string `xml:"mediaPresentationDuration,attr"`
	MaxSegmentDuration        string `xml:"maxSegmentDuration,attr"`
	Profiles                  string `xml:"profiles,attr"`

	Periods            []Period           `xml:"Period"`
	ProgramInformation ProgramInformation `xml:"ProgramInformation"`

	AvailabilityStartTime string `xml:"availabilityStartTime,attr"`
	ID                    string `xml:"id,attr"`
	MinimumUpdatePeriod   string `xml:"minimumUpdatePeriod,attr"`
	PublishTime           string `xml:"publishTime,attr"`
	TimeShiftBufferDepth  string `xml:"timeShiftBufferDepth,attr"`
	Type                  string `xml:"type,attr"`
	NS1schemaLocation     string `xml:"ns1:schemaLocation,attr"`
	BaseURL               string `xml:"BaseURL"`

	ServiceDescription ServiceDescription `xml:"ServiceDescription"`
}

// AdaptationSet in MPD
type AdaptationSet struct {
	XMLName            xml.Name `xml:"AdaptationSet"`
	SegmentAlignment   bool     `xml:"segmentAlignment,attr"`
	BitstreamSwitching bool     `xml:"bitstreamSwitching,attr"`
	MaxWidth           int      `xml:"maxWidth,attr"`
	MaxHeight          int      `xml:"maxHeight"`
	MaxFrameRate       int      `xml:"maxFrameRate"`

	Par string `xml:"par,attr"`

	Lang                      string                    `xml:"lang,attr"`
	BaseURL                   string                    `xml:"BaseURL"`
	Representation            []Representation          `xml:"Representation"`
	SegmentTemplate           []SegmentTemplate         `xml:"SegmentTemplate"`
	SegmentList               SegmentList               `xml:"SegmentList"`
	SubsegmentStartsWithSAP   int                       `xml:"subsegmentStartsWithSAP"`
	AudioChannelConfiguration AudioChannelConfiguration `xml:"AudioChannelConfiguration"`
	Role                      Role                      `xml:"Role"`
	ContentType               string                    `xml:"contentType,attr"`
	MimeType                  string                    `xml:"mimeType,attr"`
	StartWithSAP              int                       `xml:"startWithSAP,attr"`

	FrameRate string `xml:"frameRate,attr"`
	Height    string `xml:"height,attr"`
	ScanType  string `xml:"scanType,attr"`
	Width     int    `xml:"width,attr"`

	StartWithSap int    `xml:"startWithSap,attr"`
	ID           string `xml:"id,attr"`
}

// Period in MPD
type Period struct {
	XMLName       xml.Name        `xml:"Period"`
	Duration      string          `xml:"duration,attr"`
	AdaptationSet []AdaptationSet `xml:"AdaptationSet"`
	ID            string          `xml:"id,attr"`
	Start         string          `xml:"start,attr"`
}

// ServiceDescription in MPD
type ServiceDescription struct {
	XMLName xml.Name `xml:"ServiceDescription"`
	ID      int      `xml:"id,attr"`
}

// ProgramInformation in MPD
type ProgramInformation struct {
	XMLName            xml.Name `xml:"ProgramInformation"`
	MoreInformationURL string   `xml:"moreInformationURL,attr"`
	Title              string   `xml:"Title"`
}

// Representation in MPD
type Representation struct {
	XMLName                   xml.Name                  `xml:"Representation"`
	ID                        string                    `xml:"id,attr"`
	MimeType                  string                    `xml:"mimeType,attr"`
	Codecs                    string                    `xml:"codecs,attr"`
	Width                     int                       `xml:"width,attr"`
	Height                    int                       `xml:"height,attr"`
	FrameRate                 int                       `xml:"frameRate,attr"`
	Sar                       string                    `xml:"sar,attr"`
	StartWithSap              int                       `xml:"startWithSap,attr"`
	BandWidth                 int                       `xml:"bandwidth,attr"`
	BaseURL                   string                    `xml:"BaseURL"`
	SegmentTemplate           SegmentTemplate           `xml:"SegmentTemplate"`
	SegmentList               SegmentList               `xml:"SegmentList"`
	SegmentBase               SegmentBase               `xml:"SegmentBase"`
	AudioSamplingRate         int                       `xml:"audioSamplingRate,attr"`
	AudioChannelConfiguration AudioChannelConfiguration `xml:"AudioChannelConfiguration"`
}

// SegmentTemplate in MPD
type SegmentTemplate struct {
	XMLName        xml.Name `xml:"SegmentTemplate"`
	Media          string   `xml:"media,attr"`
	Timescale      int      `xml:"timescale,attr"`
	StartNumber    int      `xml:"startNumber,attr"`
	Duration       int      `xml:"duration,attr"`
	Initialization string   `xml:"initialization,attr"`
}

// SegmentList in MPD
type SegmentList struct {
	XMLName            xml.Name       `xml:"SegmentList"`
	Timescale          int            `xml:"timescale,attr"`
	Duration           int            `xml:"duration,attr"`
	SegmentURL         []SegmentURL   `xml:"SegmentURL"`
	SegmentInitization Initialization `xml:"Initialization"`
}

// AudioChannelConfiguration in MPD
type AudioChannelConfiguration struct {
	XMLName     xml.Name `xml:"AudioChannelConfiguration"`
	SchemeIDURI string   `xml:"schemeIdUri,attr"`
	Value       int      `xml:"value,attr"`
}

// SegmentBase in MPD
type SegmentBase struct {
	XMLName            xml.Name       `xml:"SegmentBase"`
	IndexRangeExact    string         `xml:"indexRangeExact,attr"`
	IndexRange         string         `xml:"indexRange,attr"`
	SegmentInitization Initialization `xml:"Initialization"`
}

// Role in MPD
type Role struct {
	XMLName     xml.Name `xml:"Role"`
	SchemeIDURI string   `xml:"schemeIdUri,attr"`
	Value       string   `xml:"value,attr"`
}

// Initialization in MPD
type Initialization struct {
	XMLName   xml.Name `xml:"Initialization"`
	SourceURL string   `xml:"sourceURL,attr"`
}

// SegmentURL in MPD
type SegmentURL struct {
	XMLName    xml.Name `xml:"SegmentURL"`
	MediaRange string   `xml:"mediaRange,attr"`
	IndexRange string   `xml:"indexRange,attr"`
}

//TODO unused atm
type MpdCodec string
type MpdCodecIndex int
type RepRateCodec string
