package utils

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/massimo-gollo/DASHpher/models"
	"github.com/massimo-gollo/DASHpher/network"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"runtime/debug"
	"strconv"
	"time"
)

func HandleError() {
	if err := recover(); err != nil {
		logrus.Errorln(err)
		debug.PrintStack()
	}
}

// GetFile provide file from HTTP server
func GetFile(originalUrl, fileURI string, info *models.SegmentInfo, segmentDuration int) (err error) {
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
			return errors.New("status code != 200")
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

	return nil
}
