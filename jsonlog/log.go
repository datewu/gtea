package jsonlog

// Debug defaultLogger.Debug
func Debug(msg string, props ...map[string]interface{}) {
	defaultLogger.Debug(msg, props...)
}

// Info defaultLogger.Info
func Info(msg string, props ...map[string]interface{}) {
	defaultLogger.Info(msg, props...)
}

// Err defaultLogger.Err
func Err(err error, props ...map[string]interface{}) {
	defaultLogger.Err(err, props...)
}

// Trace defaultLogger.Trace
func Trace(err error, props ...map[string]interface{}) {
	defaultLogger.Trace(err, props...)
}

// Fatal defaultLogger.Fatal
func Fatal(err error, props ...map[string]interface{}) {
	defaultLogger.Fatal(err, props...)
}

// Write defaultLogger.Write
func Write(msg []byte) (int, error) {
	return defaultLogger.Write(msg)
}
