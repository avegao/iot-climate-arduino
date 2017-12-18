package controller

import (
	pb "github.com/avegao/iot-climate-arduino/proto"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/avegao/iot-climate-arduino/arduino"
	"github.com/avegao/gocondi"
)

type Controller struct {
	pb.ArduinoServer
}

func (c Controller) GetTemperature(ctx context.Context, request *pb.ArduinoRequest) (response *pb.TemperatureResponse, err error) {
	const logTag = "Controller.GetTemperature()"
	logger := gocondi.GetContainer().GetLogger()
	logger.WithField("request", request).Debugf("%s -> START", logTag)

	ar := arduino.Arduino{Address: "192.168.1.163", Port: 80}
	temperature, err := ar.GetTemperature()

	if err != nil {
		logger.WithError(err).WithField("arduino", ar).Errorf("%s -> Error getting temperature", logTag)
		logger.WithError(err).Debugf("%s -> STOP", logTag)
	} else {
		logger.WithFields(logrus.Fields{"arduino": ar, "temperature": *temperature}).Debugf("%s -> Temperature gotten", logTag)
		response = &pb.TemperatureResponse{Temperature: *temperature}
		logger.WithField("response", *response).Debugf("%s -> END", logTag)
	}

	return
}
