package helper

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func ReloadAlertmanager() error {
	logrus.Infof("reload alertmanager...")

	resp, err := http.Post("http://alertmanager:9093/-/reload", "text/html", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}
