package usecase

// A Logger belong to the usecases layer.
type Logger interface {
	LogError(string, ...interface{})
	LogAccess(string, ...interface{})
}
