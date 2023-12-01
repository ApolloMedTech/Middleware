package alertManager

type AlertType string

const (
	AlertPrimary   AlertType = "alert-primary"
	AlertSecondary AlertType = "alert-secondary"
	AlertSuccess   AlertType = "alert-success"
	AlertDanger    AlertType = "alert-danger"
	AlertWarning   AlertType = "alert-warning"
	AlertInfo      AlertType = "alert-info"
	AlertLight     AlertType = "alert-light"
	AlertDark      AlertType = "alert-dark"
)

type Alert struct {
	Type    AlertType
	Message string
}
