package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chamanbravo/upstat/queries"
	"github.com/chamanbravo/upstat/serializers"
	"github.com/chamanbravo/upstat/utils"
	"github.com/gofiber/fiber/v2"
)

// @Tags Monitors
// @Accept json
// @Produce json
// @Param body body serializers.AddMonitorIn true "Body"
// @Success 200 {object} serializers.SuccessResponse
// @Success 400 {object} serializers.ErrorResponse
// @Router /api/monitors/create [post]
func CreateMonitor(c *fiber.Ctx) error {
	newMonitor := new(serializers.AddMonitorIn)
	if err := c.BodyParser(newMonitor); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := utils.BodyValidator.Validate(newMonitor)
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	monitor, err := queries.CreateMonitor(newMonitor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = queries.NotificationMonitor(monitor.ID, newMonitor.NotificationChannels)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	fmt.Println(newMonitor.StatusPages)
	err = queries.StatusPageMonitor(monitor.ID, newMonitor.StatusPages)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	utils.StartGoroutine(monitor)

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
	})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} serializers.MonitorInfoOut
// @Success 400 {object} serializers.ErrorResponse
// @Router /api/monitors/info/{id} [get]
func MonitorInfo(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	monitor, err := queries.FindMonitorById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
		"monitor": monitor,
	})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Param body body serializers.AddMonitorIn true "Body"
// @Success 200 {object} serializers.SuccessResponse
// @Success 400 {object} serializers.ErrorResponse
// @Router /api/monitors/update/{id} [put]
func UpdateMonitor(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	monitor := new(serializers.AddMonitorIn)
	if err := c.BodyParser(monitor); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := utils.BodyValidator.Validate(monitor)
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	err = queries.UpdateMonitorById(id, monitor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = queries.UpdateNotificationMonitorById(id, monitor.NotificationChannels)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = queries.UpdateStatusPageMonitorById(id, monitor.StatusPages)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
	})

}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} serializers.SuccessResponse
// @Success 400 {object} serializers.ErrorResponse
// @Router /api/monitors/pause/{id} [put]
func PauseMonitor(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	utils.StopGoroutine(int(id))
	err = queries.UpdateMonitorStatus(id, "yellow")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
	})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} serializers.SuccessResponse
// @Success 400 {object} serializers.ErrorResponse
// @Router /api/monitors/resume/{id} [put]
func ResumeMonitor(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	monitor, err := queries.FindMonitorById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	utils.StartGoroutine(monitor)
	err = queries.UpdateMonitorStatus(id, "green")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
	})
}

// @Accept json
// @Produce json
// @Success 200 {object} serializers.MonitorsListOut
// @Success 400 {object} serializers.ErrorResponse
// @Router /api/monitors/list [get]
func MonitorsList(c *fiber.Ctx) error {
	monitors, err := queries.RetrieveMonitors()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var monitorsList []fiber.Map
	for _, v := range monitors {
		heartbeat, err := queries.RetrieveHeartbeats(v.ID, 10)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		monitorItem := fiber.Map{
			"id":        v.ID,
			"name":      v.Name,
			"url":       v.Url,
			"frequency": v.Frequency,
			"status":    v.Status,
			"heartbeat": heartbeat,
		}
		monitorsList = append(monitorsList, monitorItem)
	}

	return c.Status(200).JSON(fiber.Map{
		"message":  "success",
		"monitors": monitorsList,
	})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Param startTime query time.Time true "Start Time" format(json)
// @Success 200 {object} serializers.HeartbeatsOut
// @Success 400 {object} serializers.ErrorResponse
// @Router /api/monitors/heartbeat/{id} [get]
func RetrieveHeartbeat(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	query := new(serializers.RetrieveHeartbeatIn)
	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := utils.BodyValidator.Validate(query)
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	heartbeat, err := queries.RetrieveHeartbeatsByTime(id, query.StartTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{"message": "success", "heartbeat": heartbeat})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} serializers.MonitorSummaryOut
// @Success 400 {object} serializers.ErrorResponse
// @Router /api/monitors/summary/{id} [get]
func MonitorSummary(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	averageLatency, err := queries.RetrieveAverageLatency(id, time.Now().Add(-time.Hour*24))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	dayUptime, err := queries.RetrieveUptime(id, time.Now().Add(-time.Hour*24))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	monthUptime, err := queries.RetrieveUptime(id, time.Now().Add(-time.Hour*30*24))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
		"summary": fiber.Map{
			"averageLatency": averageLatency,
			"dayUptime":      dayUptime,
			"monthUptime":    monthUptime,
		},
	})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} serializers.SuccessResponse
// @Success 400 {object} serializers.ErrorResponse
// @Router /api/monitors/delete/{id} [delete]
func DeleteMonitor(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	utils.StopGoroutine(id)
	err = queries.DeleteMonitorById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{"message": "success"})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} serializers.CertificateExpiryCountDown
// @Success 400 {object} serializers.ErrorResponse
// @Router /api/monitors/cert-exp-countdown/{id} [get]
func CertificateExpiryCountDown(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	monitor, err := queries.FindMonitorById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	response, err := http.Get(monitor.Url)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	tlsInfo := response.TLS
	cert := tlsInfo.PeerCertificates[0]
	expirationDate := cert.NotAfter
	daysUnitlExp := int(expirationDate.Sub(time.Now()).Hours() / 24)

	return c.Status(200).JSON(fiber.Map{
		"message":             "success",
		"daysUntilExpiration": daysUnitlExp,
	})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} serializers.NotificationListOut
// @Success 400 {object} serializers.ErrorResponse
// @Router /api/monitors/{id}/notifications [get]
func NotificationChannelListOfMonitor(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	notification, err := queries.FindNotificationChannelsByMonitorId(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":       "success",
		"notifications": notification,
	})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} serializers.ListStatusPagesOut
// @Success 400 {object} serializers.ErrorResponse
// @Router /api/monitors/{id}/status-pages [get]
func StatusPagesListOfMonitor(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	statusPages, err := queries.FindStatusPageByMonitorId(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":     "success",
		"statusPages": statusPages,
	})
}
