package jsonlog

type Level int8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelError
	LevelTrace
	LevelFatal
	LevelOff
)

// String ...
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelTrace:
		return "Trace"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}
