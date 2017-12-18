package arduino

import (
	"fmt"
	"github.com/avegao/gocondi"
	"encoding/json"
	"net/http"
	"io/ioutil"
)

const (
	getTemperatureUrl = "http://%s:%d/arduino/temp"
	name              = "%s:%d"
)

type arduinoInterface interface {
	GetTemperature() (float32, error)
	String() string
}

type Arduino struct {
	arduinoInterface
	Address string
	Port    int
}

func (ar Arduino) String() string {
	return fmt.Sprintf(name, ar.Address, ar.Port)
}

func (ar Arduino) GetTemperature() (*float32, error) {
	const logTag = "Arduino.GetTemperature()"
	logger := gocondi.GetContainer().GetLogger()
	logger.Debugf("%s -> START", logTag)

	url := fmt.Sprintf(getTemperatureUrl, ar.Address, ar.Port)
	response, err := doRequest(url)

	if nil != err {
		logger.WithError(err).Errorf("%s -> Error doing request", logTag)
		logger.WithError(err).Debugf("%s -> STOP", logTag)

		return nil, err
	}

	logger.WithField("temperature", response.Temperature).Debugf("%s -> END", logTag)

	return &response.Temperature, nil
}

func doRequest(url string) (*temperatureResponse, error) {
	const logTag = "Arduino.doRequest()"
	logger := gocondi.GetContainer().GetLogger()
	logger.WithField("url", url).Debugf("%s -> START", logTag)

	response, err := http.Get(url)

	if err != nil {
		logger.WithError(err).Errorf("%s -> Error with HTTP request", logTag)
		logger.WithError(err).Debugf("%s -> STOP", logTag)

		return nil, err
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		logger.WithError(err).Errorf("%s -> Error reading body response", logTag)
		logger.WithError(err).Debugf("%s -> STOP", logTag)

		return nil, err
	}

	logger.WithField("body", string(contents)).Debugf("%s -> Response gotten", logTag)

	jsonResponse := new(temperatureResponse)
	err = json.Unmarshal(contents, jsonResponse)

	if err != nil {
		logger.WithError(err).Errorf("%s -> Error parsing body response to JSON", logTag)
		logger.WithError(err).Debugf("%s -> STOP", logTag)

		return nil, err
	}

	logger.WithField("response", jsonResponse).Debugf("%s -> END", logTag)

	return jsonResponse, nil
}
