package reproduction

import (
	"github.com/massimo-gollo/DASHpher/algo"
	"github.com/massimo-gollo/DASHpher/constant"
	"github.com/massimo-gollo/DASHpher/models"
	"github.com/massimo-gollo/DASHpher/utils"
	"strings"
	"time"
)

func EndWithErr(metrics *models.ReproductionMetrics, t *time.Time, e error) (err error) {
	metrics.ReprEndTime = time.Now()
	metrics.Status = models.Aborted
	metrics.ReprDuration = time.Since(*t)
	metrics.LastErrorReason = e.Error()
	return err
}

func Stream(reproductionDetails *models.ReproductionMetrics, codec, adaptAlgorithm string, maxHeightRes, requestedStreamDuration, initBuffSeconds, MaxBufferSeconds int, nrequest uint64) (err error) {
	startTimeExecution := time.Now()
	urlResource := reproductionDetails.ContentUrl

	mpd, fetchInfo, err := models.GetMPDFrom(urlResource)
	if err != nil {
		//TODO handle properly
		return EndWithErr(reproductionDetails, &startTimeExecution, err)
	}
	reproductionDetails.FetchMpdInfo = *fetchInfo

	//TODO - check if MaxBuffer < MBT (minimum buffer time from mpd)

	//prepare streamInfo - contains parameters of streaming and info about segments
	//segmentInfo := make(map[int]*models.SegmentInfo)
	//TODO - check mimetypes and consider multiple adpset - for now assuming only video

	streamInfo := StreamInfo{
		UrlResource:       urlResource,
		MPD:               *mpd,
		Codec:             codec,
		AudioContent:      false,
		IsByteRangeMPD:    false,
		AdaptionAlgorithm: adaptAlgorithm,
		//		SegmentInformation: segmentInfo,

		StartTimeReproduction: &startTimeExecution,
		InitBuffer:            initBuffSeconds,
		MaxBuffer:             MaxBufferSeconds,
		MaxReqHeight:          maxHeightRes,
		//initializer
		CurrentSegmentInReproduction: 0,
	}

	//NOTICE ordering representations adpSet[0] (only one atm) from highest (0) to lower(n-1) - goDASH compliant
	err = mpd.ReverseRepr(urlResource)
	if err != nil {
		return EndWithErr(reproductionDetails, &startTimeExecution, err)
	}

	//omitted baseUrl := is for byteRange
	totalVideoDuration,
		totalSegmentCount,
		actualStreamDuration,
		actualSegmentCount,
		highestResIndex,
		lowestResIndex,
		bandwidthList,
		_,
		segmentDuration, err := mpd.GetReproductionDetails(maxHeightRes, requestedStreamDuration)

	if err != nil {
		return EndWithErr(reproductionDetails, &startTimeExecution, err)
	}

	streamInfo.TotalVideoDuration = totalVideoDuration
	streamInfo.TotalSegmentCount = totalSegmentCount
	streamInfo.ActualStreamDuration = actualStreamDuration
	streamInfo.ActualTotalSegmentToStream = actualSegmentCount
	streamInfo.MaxHeightReprIdx = highestResIndex
	streamInfo.MinHeightReprIdx = lowestResIndex
	streamInfo.BandwidthList = bandwidthList
	streamInfo.SingleSegmentDuration = segmentDuration

	profiles := strings.Split(mpd.Profiles, ":")
	profile := profiles[len(profiles)-2]
	streamInfo.Profile = profile
	streamInfo.CurrentRepIdx = highestResIndex
	streamInfo.MimeTypes = []int{0} //video only - adpset 0
	streamInfo.ThroughputList = []int{}

	//GET initializer file - we start at highestRate (usefull if we want evaluate initial repRate)
	//we are assuming for sure we have 1 adpSet (video only) and is not byterange or audioByteRange
	initUrl := mpd.GetFullStreamHeader(0, highestResIndex, false, false)
	target := models.JoinURL(urlResource, initUrl)

	//TODO check if we can start with highest res - based on initializers (rep 0 && 1)

	//init struct segment informations
	//consider segment init as seg 0 in map
	segmentInfo := make(map[int]*models.SegmentInfo)
	segmentInfo[0] = models.NewSegmentInfo()
	streamInfo.SegmentInformation = segmentInfo

	streamInfo.StartTimeDownloading = time.Now()
	streamInfo.NextRunTimeDownloading = time.Now()
	streamInfo.WaitToPlayCount = 0
	//GETFILE adapt algoritm and save info about initializer
	switch adaptAlgorithm {
	case constant.ConventionalAlg:
		err = models.GetFile(urlResource, target, streamInfo.SegmentInformation[0], segmentDuration)
		if err != nil {
			return EndWithErr(reproductionDetails, &startTimeExecution, err)
		}
	}

	//all set to init gettin segment -> Reproduce
	streamInfo.CurrentSegmentInReproduction += 1
	streamInfo.BufferLevel = 0

	reproductionDetails.SegmentsInfo[0] = *segmentInfo[0]

	err = Reproduce(&streamInfo, nrequest, reproductionDetails)

	return err
}

func Reproduce(si *StreamInfo, nreq uint64, repDetails *models.ReproductionMetrics) (err error) {

	var currentSegInfo *models.SegmentInfo
	for si.CurrentSegmentInReproduction <= si.ActualTotalSegmentToStream {
		segNum := si.CurrentSegmentInReproduction

		si.CurrentURLSegToStream =
			models.GetNextSegUrl(segNum, si.MPD, si.CurrentRepIdx)

		///StartTime for this seg
		currenTimeCurrentSeg := time.Now()
		switch si.AdaptionAlgorithm {
		case "conventional":
			currentSegInfo = models.NewSegmentInfo()
			err := models.GetFile(si.UrlResource, si.CurrentURLSegToStream, currentSegInfo, si.SingleSegmentDuration)
			if err != nil {
				//TODO Retry?
				logger.Errorf("[Req#%d] Error getting segment %d reason: %s", nreq, si.CurrentSegmentInReproduction, err.Error())
				currentSegInfo.NotPlayable = true
				repDetails.Status = models.Error
				repDetails.LastErrorReason = err.Error()
				si.CurrentSegmentInReproduction++
				continue
			}
		}

		//SegInfo[num].SegSize settet in GetFile
		si.CurrentSegSize = currentSegInfo.SegmentSize

		//compute times
		arrivalTime := int(time.Since(si.StartTimeDownloading).Nanoseconds() / (1000 * 1000))
		deliveryTime := int(time.Since(currenTimeCurrentSeg).Nanoseconds() / (1000 * 1000))
		thisRunTime := int(time.Since(si.NextRunTimeDownloading).Nanoseconds() / (1000 * 1000))

		//fmt.Println("ArrivalTime: ", arrivalTime)
		//fmt.Println("ThisRunTine: ", thisRunTime)

		si.NextRunTimeDownloading = time.Now()

		//we want to wait for an initial number of segments before stream begins
		if si.WaitToPlayCount >= si.InitBuffer {
			// * print the play_out logs only when the current time is >= play_out time
			//TODO notice, at the moment function return before print all segment playout
			PrintPlayout(arrivalTime, si.InitBuffer, si.SegmentInformation)

			//get current buffer
			currentBuffer := si.BufferLevel - thisRunTime
			if currentBuffer >= 0 {
				//	si.SegmentInformation[segNum].StallTime = 0
				currentSegInfo.StallTime = 0
			} else {
				/*logger.Infoln("current buffer: ", currentBuffer)*/
				//si.SegmentInformation[segNum].StallTime = currentBuffer
				currentSegInfo.StallTime = currentBuffer
				repDetails.StallCount++

			}

			// To have the bufferLevel we take the max between the remaining buffer and 0, we add the duration of the segment we downloaded
			si.BufferLevel = utils.Max(si.BufferLevel-thisRunTime, 0) + (si.SingleSegmentDuration * 1000)
			si.WaitToPlayCount += 1
		} else {
			// add to the current buffer before we start to play
			si.BufferLevel += si.SingleSegmentDuration * 1000
			// increment the waitToPlayCounter
			si.WaitToPlayCount += 1
		}

		//compare buffer level (millisec, duration of seg*1000) with MaxBuffer (from sec to millis)
		if si.BufferLevel > si.MaxBuffer*1000 {
			//sleep until maxBuffer level is reached
			sleepTime := si.BufferLevel - si.MaxBuffer*1000
			time.Sleep(time.Duration(sleepTime) * time.Millisecond)
			si.BufferLevel -= sleepTime
		}
		//notSure yet if we need this - base the play out position on the buffer level

		si.SegmentDurationTotal += si.SingleSegmentDuration * 1000

		currentSegInfo.PlayStartPosition = si.SegmentDurationTotal
		si.SegmentInformation[segNum] = currentSegInfo

		throughput := algo.CalculateThroughput(si.CurrentSegSize, deliveryTime)
		_ = throughput

		//TODO select bitrate and switch if needed && save the bitrate from the input segment (less the header info)

		if si.SegmentDurationTotal > si.ActualStreamDuration {
			//logger.Infoln("Finish Reproduction")
			break
		} else {
			repDetails.SegmentsInfo[segNum] = *currentSegInfo
			si.CurrentSegmentInReproduction += 1
		}
	}
	if repDetails.StallCount > 10 {
		repDetails.Status = models.Error
	}
	return nil
}

func PrintPlayout(currentTime, initBuffer int, segInfo map[int]*models.SegmentInfo) {
	for i, _ := range segInfo {
		if i == 0 {
			//skip initializer segment, start from 1
			continue
		}

		if currentTime >= (segInfo[i].PlayStartPosition+segInfo[initBuffer].PlayStartPosition) && !segInfo[i].Played {
			//logger.Warningf("passing seg to decoder to playout seg %s", segInfo[i].SegmentFileName)
			segInfo[i].Played = true
		}
	}

}
