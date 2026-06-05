package utils

import (
	"net/http"

	"go_template/configs/common/constants"
	"go_template/internal/models"
)

func GetXDeviceId(r *http.Request) *models.Headers {
	xDeviceId := r.Header.Get(constants.DeviceId)
	return &models.Headers{
		XDeviceId: xDeviceId,
	}
}
