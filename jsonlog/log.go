package jsonlog

// Debug proxy to defaultLogger.Debug
func Debug(msg string, props map[string]string) {
	defaultLogger.Debug(msg, props)
}

// Info proxy to defaultLogger.Info
func Info(msg string, props map[string]string) {
}

// Err proxy to defaultLogger.Err
func Err(err error, props map[string]string) {
	defaultLogger.Err(err, props)
}

// Fatal proxy to defaultLogger.Fatal
func Fatal(err error, props map[string]string) {
	defaultLogger.Fatal(err, props)
}

// Write proxy to defaultLogger.Write
func Write(msg []byte) (int, error) {
	return defaultLogger.Write(msg)
}
