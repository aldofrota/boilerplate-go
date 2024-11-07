package usecases

// ValidateIfHealthy is a use case that validates if the service is healthy
type ValidateIfHealthy interface {
	Validate() (bool, error)
}
