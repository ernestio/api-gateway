package usages

import (
	"net/http"
	"time"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
)

// Report : ...
func Report(au models.User, fromStr, toStr string) (int, []byte) {
	var err error
	var usage models.Usage
	var reportables []models.Usage
	var body []byte
	var from, to int64

	layout := "2006-01-02"

	if fromStr != "" {
		fromTime, err := time.Parse(layout, fromStr)
		if err != nil {
			h.L.Warning(err.Error())
			return http.StatusBadRequest, []byte("Invalid from parameter")
		}
		from = fromTime.Unix()
	}
	if toStr != "" {
		toTime, err := time.Parse(layout, toStr)
		if err != nil {
			h.L.Warning(err.Error())
			return http.StatusBadRequest, []byte("Invalid to parameter")
		}
		to = toTime.Unix()
	}

	if err = usage.FindAllInRange(from, to, &reportables); err != nil {
		return 500, []byte("Internal server error")
	}

	if body, err = views.RenderUsageReport(reportables); err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
