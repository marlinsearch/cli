package marlin

import (
	"bytes"
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

func (api *Api) httpPost(path string, data string) (resp *http.Response, err error) {
	client := getClient()
	url := api.formUrl(path)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Marlin-Application-Id", api.AppId)
	req.Header.Add("X-Marlin-REST-API-KEY", api.ApiKey)
	return client.Do(req)
}

func (api *Api) httpDelete(path string) (resp *http.Response, err error) {
	client := getClient()
	url := api.formUrl(path)
	req, err := http.NewRequest("DELETE", url, nil)
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

func (api *Api) getIndexes() (string, bool) {
	resp, err := api.httpGet("indexes")
	return handleResponse(resp, err)
}

func (api *Api) getApplication(appName string) (string, bool) {
	resp, err := api.httpGet("applications/" + appName)
	return handleResponse(resp, err)
}

func (api *Api) createApplication(name string) (string, bool) {
	path := "applications"
	var dat map[string]interface{}
	s := name
	if err := json.Unmarshal([]byte(name), &dat); err != nil {
		dat = make(map[string]interface{})
		dat["name"] = name
		sb, _ := json.Marshal(dat)
		s = string(sb)
	}
	resp, err := api.httpPost(path, s)
	return handleResponse(resp, err)
}

func (api *Api) createIndex(name string, numShards int) (string, bool) {
	path := "indexes"
	var dat map[string]interface{}
	s := name
	if err := json.Unmarshal([]byte(name), &dat); err != nil {
		dat = make(map[string]interface{})
		dat["name"] = name
		dat["numShards"] = numShards
		sb, _ := json.Marshal(dat)
		s = string(sb)
	}
	resp, err := api.httpPost(path, s)
	return handleResponse(resp, err)
}

func (api *Api) deleteIndex(name string) (string, bool) {
	path := "indexes/" + name
	resp, err := api.httpDelete(path)
	return handleResponse(resp, err)
}

func (api *Api) getNumIndexJobs() float64 {
	resp, err := api.httpGet("indexes/" + CliState.ActiveIndex + "/info")
	body, success := handleResponse(resp, err)
	if success {
		var dat map[string]interface{}
		if err = json.Unmarshal([]byte(body), &dat); err == nil {
			return dat["numJobs"].(float64)
		}
	}
	return 1
}

func (api *Api) addObjectsToIndex(s string) (string, bool) {
	resp, err := api.httpPost("indexes/"+CliState.ActiveIndex, s)
	return handleResponse(resp, err)
}
