package logger

import "fmt"

type Level struct {
	str string
	i   int
}

var (
	LevelPanic = Level{"PAN", 1}
	LevelError = Level{"ERR", 2}
	LevelWarn  = Level{"WRN", 3}
	LevelInfo  = Level{"INF", 4}
	LevelDebug = Level{"DBG", 5}
)

var (
	ColorReset  = "\033[m"
	ColorRed    = "\033[0;31m"
	ColorYellow = "\033[0;33m"
	ColorBlue   = "\033[0;34m"
	ColorGray   = "\033[38;5;239m"
)

func (lv Level) Fmt() string {
	switch lv {
	case LevelDebug:
		return fmt.Sprintf("%s%s%s", ColorGray, LevelDebug.str, ColorReset)
	case LevelInfo:
		return fmt.Sprintf("%s%s%s", ColorBlue, LevelInfo.str, ColorReset)
	case LevelWarn:
		return fmt.Sprintf("%s%s%s", ColorYellow, LevelWarn.str, ColorReset)
	case LevelError:
		return fmt.Sprintf("%s%s%s", ColorRed, LevelError.str, ColorReset)
	case LevelPanic:
		return fmt.Sprintf("%s%s%s", ColorRed, LevelPanic.str, ColorReset)
	default:
		panic(fmt.Errorf("invalid log Level %v", lv))
	}
}

func (lv Level) Gt(l Level) bool {
	return lv.i > l.i
}
