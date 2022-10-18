package contracts

//
// Contract for logging
//
type ILogger interface {
	Error(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
}
