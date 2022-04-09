package reproduction

import (
	"github.com/massimo-gollo/DASHpher/constant"
	"github.com/massimo-gollo/DASHpher/models"
	"github.com/massimo-gollo/DASHpher/utils"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func Stream(mpd models.MPD,
	codec, adaptiveAlgorithm, originalUrl string,
	maxHeight, streamDuration, maxNumBufferSeg, initNumBufferSeg, streamBuffer int, nreq uint64) (err error) {

	startTimeReproduction := time.Now()

	//NOTICE: Order repr of first adapSet from Highest to lower - goDASH compliant
	mpd.ReverseRepr(originalUrl)

	//k: segm number v:segmentInfo
	var segmentInfo map[int]*models.SegmentInfo

	// MPD INFO
	var originalStreamDuration int
	var originalTotalSegmentMPD int
	var actualStreamDuration int
	var actualTotalSegmentToStream int

	var lowestMPDRestIndex int
	var highestMPDResIndex int
	var originalSingleSegmentDuration int
	var bandwidthList []int
	var mpdProfile string

	//REPRODUCTION VARIABLES - start with lowestRepRate
	var currentRepRate int
	var currentSegment = 0

	//index values for types of MPD types
	var mimeTypes []int

	var isAudioOnly = false
	//TODO check if is only audio, assume is VideoCodec
	_ = isAudioOnly

	// loops over Adpset currently one adaptation set per video and audio
	//TODO i'm assuming h264 codec only - static 1 adpSet
	currentAdaptionSetIndex := 0
	mimeTypes = append(mimeTypes, currentAdaptionSetIndex)

	//TODO skip check byteRange - assumig it's not ByteRange
	//get relevant details about current MPD
	originalStreamDuration, originalTotalSegmentMPD, actualStreamDuration, actualTotalSegmentToStream,
		highestMPDResIndex, lowestMPDRestIndex, bandwidthList, _, originalSingleSegmentDuration, err = mpd.GetReproductionDetails(maxHeight, streamDuration)

	segmentInfo = make(map[int]*models.SegmentInfo)
	segmentInfo[currentSegment] = models.NewSegmentInfo()

	profiles := strings.Split(mpd.Profiles, ":")
	mpdProfile = profiles[len(profiles)-2]

	//must set to highest?
	currentRepRate = lowestMPDRestIndex

	//get initfile.mp4 at lowest rate
	initUrl := mpd.GetFullStreamHeader(0, currentRepRate, false, false)
	targetUrl := utils.JoinURL(originalUrl, initUrl)

	//GETFILE adapt algoritm
	switch adaptiveAlgorithm {
	case constant.ConventionalAlg:
		err = utils.GetFile(originalUrl, targetUrl, segmentInfo[0], originalSingleSegmentDuration)
		if err != nil {
			return err
		}
	}
	//we have saved initSegment info, let's start with other segments from 1 to..
	currentSegment++

	st := models.StreamStruct{
		//global info about reproduction
		OriginalStreamDuration:  originalStreamDuration,
		OriginalTotalSegmentMPD: originalTotalSegmentMPD,
		OriginalUrl:             originalUrl,
		OriginalSegSize:         originalSingleSegmentDuration,
		MaxHeightReprIdx:        highestMPDResIndex,
		MinHeightReprIdx:        lowestMPDRestIndex,
		BandwidthList:           bandwidthList,
		Profile:                 mpdProfile,
		MPD:                     mpd,
		Codec:                   codec,
		IsByteRangeMPD:          false,
		StartTimeReproduction:   &startTimeReproduction,

		ActualStreamDuration:         actualStreamDuration,
		ActualTotalSegmentToStream:   actualTotalSegmentToStream,
		MaxReqHeight:                 maxHeight,
		InitBuffer:                   initNumBufferSeg,
		MaxBuffer:                    maxNumBufferSeg,
		AdaptionAlgorithm:            adaptiveAlgorithm,
		CurrentSegmentInReproduction: currentSegment,
		//start with Lower repr - alias RepRate godash
		CurrentHeightReprIdx: lowestMPDRestIndex,
		//map with info about all segments
		MapSegmentInfo: segmentInfo,
		MimeTypes:      mimeTypes,

		NextSegmentNumber: 0,

		BufferLevel:          0,
		SegmentDurationTotal: 0,
		BaseURL:              "",
		AudioContent:         false,
		RepRate:              0,
	}

	ReproduceSegments(&st)

	endrep := time.Since(startTimeReproduction)

	logrus.Infof("[REQ#%d] Total duration reproduction %s", nreq, endrep.String())

	//TODO return metrics?
	return nil
}

func ReproduceSegments(streamStruct *models.StreamStruct) {

	//current milliseconds
	//var currentStreamDuration int = 0
	var stopPlay bool = false

	var waitToPlayerCounter int = 0
	_ = waitToPlayerCounter

	//iterate over all segment to reproduce
	for segNum := streamStruct.CurrentSegmentInReproduction; segNum <= streamStruct.ActualTotalSegmentToStream; {
		if !stopPlay {
			streamStruct.MapSegmentInfo[segNum] = models.NewSegmentInfo()
			streamStruct.MapSegmentInfo[segNum].SegmentIndex = segNum

			//GET SegmentUrl
			streamStruct.CurrentURLSegToStream =
				models.GetNextSegUrl(streamStruct.CurrentSegmentInReproduction, streamStruct.MPD, streamStruct.CurrentHeightReprIdx)
			currentTime := time.Now()
			switch streamStruct.AdaptionAlgorithm {
			case constant.ConventionalAlg:
				err := utils.GetFile(streamStruct.OriginalUrl, streamStruct.CurrentURLSegToStream, streamStruct.MapSegmentInfo[segNum], streamStruct.OriginalSegSize)
				if err != nil {
					logrus.Fatalf("Error getting segment %d with error: %s", segNum, err.Error())
				}
				//	logrus.Infof("Downloaded seg #%d with RTT %d", segNum, streamStruct.MapSegmentInfo[segNum].NetDetails.RTT2FirstByte)
			}

			//compute ArrTime and DeliveryTime - useless compute but i trust in godash
			arrTime := int(time.Since(*streamStruct.StartTimeReproduction).Nanoseconds() / (1000 * 1000))
			deliveryTime := int(time.Since(currentTime).Nanoseconds() / (1000 * 1000))
			//Not sure about this
			//thisRunTimeVal := int(time.Since(nextRunTime).Nanoseconds() / (glob.Conversion1000 * glob.Conversion1000))

			_, _ = arrTime, deliveryTime

			// check if the buffer level is higher than the max buffer
			if sleepTime := streamStruct.BufferLevel - streamStruct.MaxBuffer*1000; sleepTime > 0 {
				//should sleep a bit
				// retrieve the time it is going to sleep from the buffer level
				// sleep until the max buffer level is reached
				time.Sleep(time.Duration(sleepTime))
				// reset the buffer to the new value less sleep time - should equal maxBuffer
				streamStruct.BufferLevel -= sleepTime
			}

			//Pass to next Segment
			streamStruct.CurrentSegmentInReproduction++
		}
	}
	logrus.Infoln("Ending reproduction")

}
