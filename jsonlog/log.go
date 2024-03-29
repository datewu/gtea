package jsonlog

// Debug defaultLogger.Debug
func Debug(msg string, props ...map[string]any) {
	defaultLogger.Debug(msg, props...)
}

// Info defaultLogger.Info
func Info(msg string, props ...map[string]any) {
	defaultLogger.Info(msg, props...)
}

// Err defaultLogger.Err
func Err(err error, props ...map[string]any) {
	defaultLogger.Err(err, props...)
}

// Trace defaultLogger.Trace
func Trace(err error, props ...map[string]any) {
	defaultLogger.Trace(err, props...)
}

// Fatal defaultLogger.Fatal
func Fatal(err error, props ...map[string]any) {
	defaultLogger.Fatal(err, props...)
}

// Write defaultLogger.Write
func Write(msg []byte) (int, error) {
	return defaultLogger.Write(msg)
}
