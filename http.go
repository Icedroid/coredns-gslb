package gslb

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

const ipipNetBaseURL = "http://freeapi.ipip.net/"
const defaultZoneIP = "14.215.109.226"

var cityMap = map[string]string{
	"广州": "14.215.109.226",
	"内蒙": "93.184.216.33",
}

// city to zone ip
// TODO: change to use external services
func getIPNetURL(url, ip string) string {
	if url == "" {
		url = ipipNetBaseURL
	}
	return strings.TrimRight(url, "/") + "/" + ip
}

// get ip location info from ipip.net
func requestURL(url string) (ip string, err error) {
	ip = defaultZoneIP
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		// TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	response, err := netClient.Get(url)
	if err != nil || response.StatusCode != 200 {
		log.Errorf("http request url=%s, error=%s", url, err)
		return
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return "", err
	}
	log.Infof("ipip.net response data=%v", v)

	if sl, ok := v.([]interface{}); ok && len(sl) > 3 { //ipip.net
		if city, ok := sl[2].(string); ok && city != "" {
			if c, ok := cityMap[city]; ok {
				log.Infof("zone ip=%s", c)
				ip = c
			}
		}
	}

	return
}
