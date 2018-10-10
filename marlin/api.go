package marlin

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Api struct {
	Host   string
	AppId  string
	ApiKey string
}

func getClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client
}

func (api *Api) formUrl(path string) string {
	var url = api.Host + "/1/" + path
	return url
}

func (api *Api) httpGet(path string) (resp *http.Response, err error) {
	client := getClient()
	url := api.formUrl(path)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Marlin-Application-Id", api.AppId)
	req.Header.Add("X-Marlin-REST-API-KEY", api.ApiKey)
	return client.Do(req)
}

func (api *Api) Connect() (bool, string) {
	resp, err := api.httpGet("marlin")
	if err != nil || resp.StatusCode != 200 {
		return false, "0"
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	defer resp.Body.Close()

	if version, ok := result["version"]; ok {
		return true, version.(string)
	}
	return false, "0"
}

func handleResponse(resp *http.Response, err error) (string, bool) {
	if err != nil || resp.StatusCode != 200 {
		return "", false
	}
	bodyBytes, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return "", false
	}
	bodyString := string(bodyBytes)
	return bodyString, true
}

func (api *Api) getInfo() (string, bool) {
	resp, err := api.httpGet("marlin")
	return handleResponse(resp, err)
}

func (api *Api) getApplications() (string, bool) {
	resp, err := api.httpGet("applications")
	return handleResponse(resp, err)
}
