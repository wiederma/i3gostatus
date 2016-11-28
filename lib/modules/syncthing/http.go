package syncthing

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

var csrfName string
var csrfToken string

type noActiveSessionError string

func (e noActiveSessionError) Error() string {
	return fmt.Sprintf("No active session for: %s", e)
}

func initHTTPSession(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Println("initHTTPSession failed. Is Syncthing running?")
		return
	}
	defer resp.Body.Close()

	csrfName = "X-" + strings.Split(resp.Header.Get("Set-Cookie"), "=")[0]
	csrfToken = strings.Split(resp.Header.Get("Set-Cookie"), "=")[1]
}

func stGet(baseUrl string, endpoint string) (string, error) {
	reqUrl, _ := url.Parse(baseUrl)
	reqUrl.Path = path.Join(reqUrl.Path, endpoint)
	req, err := http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add(csrfName, csrfToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return "", noActiveSessionError(reqUrl.String())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body), nil
}
