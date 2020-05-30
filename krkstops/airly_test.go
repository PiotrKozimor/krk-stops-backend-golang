package krkstops

import (
	"net/http"
	"strings"
	"testing"

	krk_stops "github.com/PiotrKozimor/krk-stops-backend-golang/krkstops-grpc"
)

func TestGetAirly(test *testing.T) {
	app := App{}
	app.HTTPClient = &http.Client{}
	airly, err := app.GetAirly(&krk_stops.Installation{Id: 9914})
	if err != nil {
		test.Error(err)
	}
	if airly.Humidity < 0 {
		test.Errorf("Humdidity %d < 0", airly.Humidity)
	}
	if airly.Caqi < 0 {
		test.Errorf("CAQI %d < 0", airly.Caqi)
	}
	if strings.Contains(airly.Color, "#") != true {
		test.Errorf("Got improper color string: %s", airly.Color)
	}
}