package main

import (
	"net/http"
	"net/url"
	"errors"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

// Parameters needed for all calls of to the Majestic API are stored in this struct.
// Currently, that's just the API key.
type MajesticApi struct {
	apiKey string
}

// The Majestic API endpoint
const baseUrl string = "https://api.majestic.com/api/json"

// Majestic has two indices you can use.
type Datasource string
const(
	Historic Datasource = "historic"
	Fresh               = "fresh"
)

type DeletedLinkMode int
const(
	KeepDeleted DeletedLinkMode = 0
	RemoveDeleted 				= 1
)

// Majestic's JSON API returns errors in this format.
type MajesticError struct {
	Code string
	ErrorMessage string
	FullError string
}

// MajesticError implements the error interface
func (err *MajesticError) Error() string {
	return fmt.Sprintf(
		"Majestic API Error: Code='%s', ErrorMessage='%s', FullError='%s'",
		err.Code, err.ErrorMessage, err.FullError)
}

func bool2int(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func (api *MajesticApi) makeCall(cmd string, options map[string]string) (string, error) {
	query := make(url.Values)
	for k, v := range options {
		query.Add(k, v)
	}
	query.Add("app_api_key", api.apiKey)
	query.Add("cmd", cmd)

	url := baseUrl + "?" + query.Encode()

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("HTTP error code %d.", response.StatusCode))
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var apiErr MajesticError
	json.Unmarshal(body, &apiErr)
	if apiErr.Code != "OK" {
		return "", &apiErr
	}
	return string(body), nil
}

type GetBackLinkDataRequest struct {
	Datasource Datasource
	Item string
	Count int32
	Mode DeletedLinkMode
	ShowDomainInfo bool
	MaxSourceURLsPerRefDomain int32
	MaxSameSourceURLs int32
	RefDomain string
	FilterTopic string
	FilterTopicsRefDomainsMode bool
	UsePrefixScan bool
	api *MajesticApi
}

func (req *GetBackLinkDataRequest) Perform() (string, error) {
	if req.Datasource == "" {
		return "", errors.New("Datasource must be set.")
	}
	if req.Item == "" {
		return "", errors.New("Item must be set")
	}
	options := make(map[string]string)
	options["datasource"] = string(req.Datasource)
	options["item"] = req.Item
	options["Count"] = fmt.Sprint(req.Count)
	options["Mode"] = fmt.Sprint(int(req.Mode))
	options["ShowDomainInfo"] = fmt.Sprint(bool2int(req.ShowDomainInfo))
	options["MaxSourceURLsPerRefDomain"] = fmt.Sprint(req.MaxSourceURLsPerRefDomain)
	options["MaxSameSourceURLs"] = fmt.Sprint(req.MaxSameSourceURLs)
	if req.RefDomain != "" {
		options["RefDomain"] = req.RefDomain
	}
	if req.FilterTopic != "" {
		options["FilterTopic"] = req.FilterTopic
	}
	if req.FilterTopicsRefDomainsMode {
		options["FilterTopicsRefDomainsMode"] = "1"
	}
	if req.UsePrefixScan {
		options["UsePrefixScan"] = "1"
	}

	return req.api.makeCall("GetBackLinkData", options)
}

func (api *MajesticApi) GetBackLinkData() GetBackLinkDataRequest {
	var req GetBackLinkDataRequest
	req.Count = 100
	req.Mode = KeepDeleted
	req.ShowDomainInfo = false
	req.MaxSourceURLsPerRefDomain = -1
	req.MaxSameSourceURLs = -1
	req.FilterTopicsRefDomainsMode = false
	req.UsePrefixScan = false

	req.api = api
	return req
}