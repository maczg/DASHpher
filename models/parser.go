package models

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/massimo-gollo/DASHpher/network"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"path"
	"strings"
	"time"
)

//TODO requestMetadata - cool if we can store metrics about each request
// see httptrace.ClientTrace

func ParseMPDFrom(mpdBody *[]byte) (mpd *MPD, err error) {
	//extract everything from the file read in bytes to the structures
	if err = xml.Unmarshal(*mpdBody, &mpd); err != nil {
		return nil, err
	}
	return mpd, nil
}

//GetMPDFrom - get mpd from requested url
func GetMPDFrom(requestedUrl *string) (mpd *MPD, requestMetadata *network.FileMetadata, err error) {

	url := strings.TrimSpace(*requestedUrl)
	//Get Custom http client - ulimit timeouts
	client := network.NewCustomHttp()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		//if redirected prints
		//	INFO[2022-04-13 22:23:35] redirected - http://cloud.gollo1.particles.dieei.unict.it/vms/videos/624d99627f1af072aead0c47
		//	INFO[2022-04-13 22:23:35] redirected - http://cloud.gollo1.particles.dieei.unict.it/videofiles/624d99627f1af072aead0c47/video.mpd

		logger.Infof("redirected - %s", req.URL.String())
		*requestedUrl = req.URL.String()
		return nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var startTime time.Time

	fetchingInfo := network.FileMetadata{}
	tracer := network.GetTraceRequestFile(&fetchingInfo, &startTime)
	clientTraceCtx := httptrace.WithClientTrace(req.Context(), tracer)
	req = req.WithContext(clientTraceCtx)

	startTime = time.Now()
	resp, err := client.Do(req)

	if err != nil {
		return nil, &fetchingInfo, err
	}

	if resp.StatusCode == http.StatusFound { //status code 302
		fmt.Println(resp.Location())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &fetchingInfo, errors.New("NOT 200")
	}
	defer resp.Body.Close()

	//Resolve MPD
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &fetchingInfo, err
	}
	mpd, err = ParseMPDFrom(&body)
	if err != nil {
		return nil, &fetchingInfo, err
	}

	return mpd, &fetchingInfo, nil
}

// JoinURL
// Return full path of what we must download
func JoinURL(baseURL string, append string) string {
	// if "append" already contains "http", then do nothing
	if !(strings.Contains(append, "http")) {
		// get the base of the current url
		base := path.Base(baseURL)
		// replace this base url with the required file string
		urlHeaderString := strings.Replace(baseURL, base, append, -1)
		// return the new url
		return urlHeaderString
	}
	// return the new url
	return append
}
