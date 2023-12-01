package alertManager

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const SessionKey = "alerts"
const MaxAlerts = 10 // define the maximum number of alerts

func AlertMiddleware(store sessions.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		alerts, err := getAlertsFromSession(c)
		if err != nil {
			// handle error
		}
		c.Set(SessionKey, alerts)
		c.Next()
	}
}

func AddAlert(c *gin.Context, msg string, typ AlertType) {
	session := sessions.Default(c)
	alerts, err := getAlertsFromSession(c)
	if err != nil {
		// handle error
	}

	// If the number of alerts has reached the limit, remove the oldest alert
	if len(alerts) >= MaxAlerts {
		alerts = alerts[1:]
	}

	alerts = append(alerts, Alert{Type: typ, Message: msg})
	session.Set(SessionKey, alerts)
	err = session.Save()
	if err != nil {
		// handle error
	}
}

func GetAlerts(c *gin.Context) []Alert {
	logrus.Debug("Getting alerts")
	alerts, err := getAlertsFromSession(c)
	if err != nil {
		// handle error
	}
	return alerts
}

func getAlertsFromSession(c *gin.Context) ([]Alert, error) {
	session := sessions.Default(c)
	rawAlerts := session.Get(SessionKey)
	if rawAlerts == nil {
		return nil, nil
	}
	alerts, ok := rawAlerts.([]Alert)
	if !ok {
		logrus.Error("Failed to assert alerts from session: ")
		return nil, fmt.Errorf("failed to assert alerts from session")
	}

	session.Delete(SessionKey)
	err := session.Save()
	if err != nil {
		logrus.Error("Failed to delete alerts from session: ", err)
		return nil, fmt.Errorf("failed to delete alerts from session")
	}
	return alerts, nil
}
