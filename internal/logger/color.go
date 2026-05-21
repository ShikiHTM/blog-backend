package logger

import "fmt"

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

func Info(format string, v ...any) string {
	return fmt.Sprintf(colorGreen+"[INFO] "+colorReset+format, v...)
}

func Warn(format string, v ...any) string {
	return fmt.Sprintf(colorYellow+"[WARN] "+colorReset+format, v...)
}

func Error(format string, v ...any) string {
	return fmt.Sprintf(colorRed+"[ERROR] "+colorReset+format, v...)
}

func System(format string, v ...any) string {
	return fmt.Sprintf(colorCyan+"[MAIN] "+colorReset+format, v...)
}
