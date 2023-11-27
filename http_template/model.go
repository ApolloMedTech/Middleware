type AlertType string

const (
    AlertSuccess AlertType = "alert-success"
    AlertDanger  AlertType = "alert-danger"
    AlertWarning AlertType = "alert-warning"
    AlertInfo    AlertType = "alert-info"
)

type Alert struct {
    Type AlertType
    Text string
}
