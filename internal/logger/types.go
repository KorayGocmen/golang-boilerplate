package logger

import "strings"

type Mode string

const (
	LevelFatal = iota
	LevelEmerg
	LevelAlert
	LevelCrit
	LevelError
	LevelWarng
	LevelNotice
	LevelInfo
	LevelDebug

	ModeFile    Mode = "FILE"
	ModeAWS     Mode = "AWS"
	ModeSyslog  Mode = "SYSLOG"
	ModeConsole Mode = "CONSOLE"
	ModeNone    Mode = "NONE"
)

func ToMode(mode string) Mode {
	return Mode(strings.TrimSpace(strings.ToUpper(mode)))
}
