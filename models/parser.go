package models

import (
	"encoding/xml"
	"errors"
	"github.com/massimo-gollo/DASHpher/network"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
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
func GetMPDFrom(requestedUrl string) (mpd *MPD, requestMetadata interface{}, err error) {
	//Get Custom http client - ulimit timeouts
	client := network.NewCustomHttp()

	url := strings.TrimSpace(requestedUrl)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := client.Do(req)

	if err != nil {
		return nil, nil, err
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New("can't get response or status code not OK")
	}
	//Resolve MPD
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	mpd, err = ParseMPDFrom(&body)
	if err != nil {
		return nil, nil, err
	}

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	return mpd, nil, nil
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
