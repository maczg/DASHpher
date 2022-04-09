package models

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/massimo-gollo/DASHpher/constant"
	"github.com/massimo-gollo/DASHpher/network"
	"github.com/massimo-gollo/DASHpher/utils"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"strconv"
	"strings"
	"time"
)

//GetSupportedCodec - return an array of codec offered in MPD - assuming all repr are in one adp set
func (m *MPD) GetSupportedCodec() (codecs []string, hasAudioCodec bool) {
	//Evaluate only first Period - Assume video are Single-Period - assume single repr
	for _, currentAdpSet := range m.Periods[0].AdaptationSet {
		mpdCodec := currentAdpSet.Representation[0].Codecs
		var repRateCodec string
		switch {
		case strings.Contains(mpdCodec, "avc"):
			repRateCodec = constant.RepCodecAVC
		case strings.Contains(mpdCodec, "hev"):
			repRateCodec = constant.RepCodecHEVC
		case strings.Contains(mpdCodec, "mp4a"):
			repRateCodec = constant.RepCodecAudio
			hasAudioCodec = true
		case strings.Contains(mpdCodec, "ac-3"):
			repRateCodec = constant.RepCodecAudio
			hasAudioCodec = true
		default:
			repRateCodec = "Unknown"
		}
		codecs = append(codecs, repRateCodec)
	}
	return
}

// GetSegmentDetails returns NumberOfSegment and their durations in sec
func (m *MPD) GetSegmentDetails() (maxSegments int, segmentsDuration int, err error) {
	defer utils.HandleError()
	streamDuration, err := m.ParseDurationOf(m.MediaPresentationDuration)

	if err != nil {
		return 0, 0, err
	}

	//TODO imported from godash - refactor me
	//mpd.MaxSegmentDuration may not be the actual segment size (just the size of the last segment)
	//segmentDuration = splitMPDSegmentDuration(mpd.MaxSegmentDuration)
	duration := m.Periods[0].AdaptationSet[0].Representation[0].SegmentTemplate.Duration
	timeScale := m.Periods[0].AdaptationSet[0].Representation[0].SegmentTemplate.Timescale

	if duration == 0 {
		duration = m.Periods[0].AdaptationSet[0].SegmentTemplate[0].Duration
	}
	if timeScale == 0 {
		timeScale = m.Periods[0].AdaptationSet[0].SegmentTemplate[0].Timescale
	}

	// this might be a byte-range, so return empty if timeScale is empty
	if timeScale == 0 {
		timeScale = 1
	}

	//TODO divisions error prone

	// get segment duration
	segmentsDuration = duration / timeScale
	maxSegments = *streamDuration / segmentsDuration
	return maxSegments, segmentsDuration, nil
}

//ParseDurationOf take duration in MediaFormat and return duration in seconds
func (m *MPD) ParseDurationOf(duration string) (timeSeconds *int, err error) {
	timeSeconds = new(int)
	var streamDuration string

	// lets first determine the length of the file
	// remove the "PT"
	streamDurationHMS := strings.Replace(duration, "PT", "", -1)

	// if streamDurationHMS contains hours
	if strings.Contains(streamDurationHMS, "H") {
		// get the hours
		H := strings.Split(streamDurationHMS, "H")
		streamDurationH := H[0]
		// if there are hours, convert to seconds
		i3, err := strconv.Atoi(streamDurationH)
		if err != nil {
			logger.Errorf("Error parsing hours")
			return nil, err
		}
		if i3 > 0 {
			*timeSeconds += i3 * 60 * 60
		}
		streamDuration = H[1]
	} else {
		// remove the "PT0H"
		// this can't contain PT0H, as H was not found in the if check
		// streamDuration = strings.Replace(mpdSegDuration, "PT0H", "", -1)
		// PT feels better
		streamDuration = strings.Replace(duration, "PT", "", -1)
	}

	// split around the Minutes
	if strings.Contains(streamDuration, "M") {
		m := strings.Split(streamDuration, "M")
		// if there are minutes, convert to seconds
		i1, err := strconv.Atoi(m[0])
		if err != nil {
			logger.Errorf("Error parsing minutes")
			return nil, err
		}
		if i1 > 0 {
			*timeSeconds += i1 * 60
		}
		// get the seconds and convert to int
		streamDuration = m[1]
	}

	// split around the Seconds
	if strings.Contains(streamDuration, "S") {
		// get the seconds and convert to int
		s := strings.Split(streamDuration, ".")
		i2, err := strconv.Atoi(s[0])
		if err != nil {
			logger.Errorf("Error parsing seconds")
			return nil, err
		}
		if i2 > 0 {
			*timeSeconds += i2
		}
	}

	// return the hours, minutes and seconds (in seconds)
	return timeSeconds, nil
}

//GetIndexReprMaxHeight return indexes of max Representation Resolution in AdaptationSet based on MaxRes requested and min res (assuming 2160 is the max res)
func (m *MPD) GetIndexReprMaxHeight(maxHeight int) (maxHeightIdx, minHeightIdx int, err error) {
	//only one repr
	if len(m.Periods[0].AdaptationSet[0].Representation) == 1 {
		maxHeightIdx = 0
		minHeightIdx = 0
		return
	}
	maxVal := 0
	minVal := 2160

	//TODO we are assuming that same repr(with res) can't have different bitrate but it could be possible to have multiple rep_rates for given resolution
	for i, v := range m.Periods[0].AdaptationSet[0].Representation {
		if v.Height > maxHeight {
			continue
		}
		if v.Height > maxVal {
			maxVal = m.Periods[0].AdaptationSet[0].Representation[i].Height
			maxHeightIdx = i
		}
		if v.Height < minVal {
			minVal = m.Periods[0].AdaptationSet[0].Representation[i].Height
			minHeightIdx = i
		}
		//TODO fix me later
		//if our desired res is lowest of minimum res in mpd, set lowest index
	}
	return
}

// GetReproductionDetails return relevant info about reproduction. maxStream duration = 0 means reproduce entire mpd
// maxRequestedStreamDuration is in millisec
func (m *MPD) GetReproductionDetails(maxHeight, requestedStreamDuration int) (
	originalStreamDuration, originalTotalSegmentMPD,
	actualStreamDuration, actualTotalSegmentToStream,
	maxHeightReprIndex, minHeightIndex int, bandwidthList []int,
	baseUrl string,
	singleSegmentDurationSeconds int, err error) {

	//TODO skip - GetDetails in ByteRange format
	originalTotalSegmentMPD, singleSegmentDurationSeconds, err = m.GetSegmentDetails()
	lastSegDuration, err := m.ParseDurationOf(m.MaxSegmentDuration)
	originalStreamDuration = ((singleSegmentDurationSeconds * (originalTotalSegmentMPD - 1)) + (*lastSegDuration)) * 10000
	//compute cap duration
	if requestedStreamDuration != 0 {
		//we are specifying cap duration, so cap maxDuration
		actualStreamDuration = requestedStreamDuration
		if actualStreamDuration%singleSegmentDurationSeconds != 0 {
			actualTotalSegmentToStream = int(float64(actualStreamDuration)/float64(singleSegmentDurationSeconds*1000) + 1)
		}
		actualTotalSegmentToStream = int(float64(actualStreamDuration)/float64(singleSegmentDurationSeconds*1000) + 1)
	} else {
		if err != nil {
			return 0, 0, 0, 0, 0, 0, nil, "", 0, nil
		}

		actualStreamDuration = originalStreamDuration
		actualTotalSegmentToStream = originalTotalSegmentMPD
	}
	maxHeightReprIndex, minHeightIndex, err = m.GetIndexReprMaxHeight(maxHeight)
	bandwidthList = m.GetBandWidths()
	baseUrl = m.Periods[0].AdaptationSet[0].BaseURL
	return
}

func (m *MPD) GetBandWidths() (bandwidthList []int) {
	for i, _ := range m.Periods[0].AdaptationSet[0].Representation {
		bandwidthList = append(bandwidthList, m.Periods[0].AdaptationSet[0].Representation[i].BandWidth)
	}

	return
}

//ReverseRepr order representation of first AdaptionSet from the highest Resolution to the lower
func (m *MPD) ReverseRepr(url string) {
	// get the current adaptation set, number of representations and min and max index based on max resolution height
	mpdReprLength := len(m.Periods[0].AdaptationSet[0].Representation)
	lowestBandWidth := m.Periods[0].AdaptationSet[0].Representation[0].BandWidth
	highestBandWidth := m.Periods[0].AdaptationSet[0].Representation[mpdReprLength-1].BandWidth

	// if the MPD is reversed (index 0 for represenstion is the lowest rate)
	// then reverse the represenstions
	if lowestBandWidth < highestBandWidth {
		//TODO fix me - we can reduce one call working with pointer to reorder repr
		r, _, _ := GetMPDFrom(url)

		// create it with content

		// loop over the existing list and reverse the representations
		i := 0
		for j := mpdReprLength - 1; j >= 0; j-- {
			// save the lowest index of structList in the highest index of reversedStructList
			r.Periods[0].AdaptationSet[0].Representation[j] = m.Periods[0].AdaptationSet[0].Representation[i]
			// reset the ID number of reversedStructList
			r.Periods[0].AdaptationSet[0].Representation[j].ID = strconv.Itoa(j + 1)
			// increment i
			i = i + 1
		}
		//reset the structlist to the new rates
		m.Periods[0].AdaptationSet[0] = r.Periods[0].AdaptationSet[0]
	}
}

//GetFullStreamHeader - return header file for current videoclip - FULL indicate for the full profile
func (m *MPD) GetFullStreamHeader(currentAdpSet, SegQuality int, isByteRange, isAudioByteRange bool) (initUrl string) {
	if isAudioByteRange {
		return m.Periods[0].AdaptationSet[currentAdpSet].Representation[(SegQuality)].BaseURL
	} else if isByteRange {
		return m.Periods[0].AdaptationSet[currentAdpSet].SegmentList.SegmentInitization.SourceURL
	}

	if m.Periods[0].AdaptationSet[currentAdpSet].SegmentTemplate != nil {
		return m.Periods[0].AdaptationSet[currentAdpSet].SegmentTemplate[0].Initialization
	} else {
		return m.Periods[0].AdaptationSet[currentAdpSet].Representation[(SegQuality)].SegmentTemplate.Initialization
	}
}

func GetNextSegUrl(segmentNumber int, mpd MPD, SegQUALITY int) string {
	var repRateBaseUrl string
	mimType := mpd.Periods[0].AdaptationSet[0].Representation[(SegQUALITY)].MimeType

	// check if this call is for audio
	if mimType == constant.RepCodecVideo {
		repRateBaseUrl = mpd.Periods[0].AdaptationSet[0].Representation[(SegQUALITY)].SegmentTemplate.Media
	} else {
		repRateBaseUrl = mpd.Periods[0].AdaptationSet[0].SegmentTemplate[0].Media
	}

	// convert the segment int to string
	nb := strconv.Itoa(segmentNumber)
	//Replace $Number$ by the segment number in the url and return it
	return strings.Replace(repRateBaseUrl, "$Number$", nb, -1)

}

// GetFile provide file from HTTP server
func GetFile(originalUrl, fileURI string, info *SegmentInfo, segmentDuration int) (err error) {
	fullUrl := JoinURL(originalUrl, fileURI)

	client := network.NewCustomHttp()

	request, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return err
	}
	var startTime time.Time

	tracer := network.GetTraceRequestFile(&info.NetDetails, &startTime)
	clientTraceCtx := httptrace.WithClientTrace(request.Context(), tracer)
	request = request.WithContext(clientTraceCtx)

	startTime = time.Now()
	resp, err := client.Do(request)

	if not200 := resp.StatusCode != http.StatusOK; err != nil || not200 {
		if not200 {
			s := fmt.Sprintf("Status code: %s", resp.Status)
			return errors.New(s)
		}
		return err
	}
	defer resp.Body.Close()

	//parse resp in mpd
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	info.HTTPprotocol = resp.Proto
	info.SegmentFileName = fileURI

	//get body size
	size := strconv.FormatInt(int64(len(body)), 10)
	segSize, err := strconv.Atoi(size)
	if err != nil {
		return err
	}

	// get the P.1203 segSize (less the header)
	withoutHeaderVal := int64(segSize)

	src := []byte("0000000468EFBC80")
	dst := make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(dst, src)
	if err != nil {
		return err
	}

	if bytes.Contains(body, dst[:n]) {
		// get the index for our dst value
		mdatValueInt := bytes.Index(body, dst[:n])
		// add 8 bits for header
		mdatValueInt += 8
		// get the file byte size less the header
		withoutHeaderVal = int64(segSize) - int64(mdatValueInt)
	}

	kbpsInt := (withoutHeaderVal * 8) / int64(segmentDuration)
	// convert kbps to a float
	kbpsFloat := float64(kbpsInt) / 1024
	// convert to sn easier string value
	kbpsFloatStringVal := fmt.Sprintf("%3f", kbpsFloat)
	_ = kbpsFloatStringVal

	info.SegmentSize = segSize

	return nil
}
