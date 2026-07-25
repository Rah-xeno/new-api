package controller

import (
	"strconv"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
)

func adminAnalyticsTimezoneOffset(c *gin.Context) int {
	raw := strings.TrimSpace(c.Query("tz_offset_minutes"))
	if raw != "" {
		offset, err := strconv.Atoi(raw)
		if err == nil && offset >= -14*60 && offset <= 14*60 {
			return offset
		}
	}
	_, offsetSeconds := time.Now().Zone()
	return offsetSeconds / 60
}

func GetAdminAnalyticsReport(c *gin.Context) {
	startTimestamp, startErr := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, endErr := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	if startErr != nil || endErr != nil || startTimestamp <= 0 || endTimestamp <= 0 {
		common.ApiErrorMsg(c, "invalid analytics time range")
		return
	}

	report, err := service.GetAdminAnalyticsReport(
		startTimestamp,
		endTimestamp,
		adminAnalyticsTimezoneOffset(c),
	)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, report)
}
