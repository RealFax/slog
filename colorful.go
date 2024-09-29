package log

import (
	"log/slog"
	"strconv"
)

type Color int32

const (
	reset = "\033[0m"

	black        Color = 30
	red          Color = 31
	green        Color = 32
	yellow       Color = 33
	blue         Color = 34
	magenta      Color = 35
	cyan         Color = 36
	lightGray    Color = 37
	darkGray     Color = 90
	lightRed     Color = 91
	lightGreen   Color = 92
	lightYellow  Color = 93
	lightBlue    Color = 94
	lightMagenta Color = 95
	lightCyan    Color = 96
	white        Color = 97
)

func (c Color) String() string {
	switch c {
	case black:
		return "30"
	case red:
		return "31"
	case green:
		return "32"
	case yellow:
		return "33"
	case blue:
		return "34"
	case magenta:
		return "35"
	case cyan:
		return "36"
	case lightGray:
		return "37"
	case darkGray:
		return "90"
	case lightRed:
		return "91"
	case lightGreen:
		return "92"
	case lightYellow:
		return "93"
	case lightBlue:
		return "94"
	case lightMagenta:
		return "95"
	case lightCyan:
		return "96"
	case white:
		return "97"
	default:
		return strconv.Itoa(int(c))
	}
}

var (
	colorHeader      = []byte("\033[")
	colorEnd         = []byte(reset)
	colorAgent  byte = 'm'
)

func colorize(c Color, v string) []byte {
	b := make([]byte, 0, 7+len(c.String())+len(v))
	return append(append(append(append(append(b, colorHeader...), []byte(c.String())...), colorAgent), []byte(v)...), colorEnd...)
}

func getColor(level slog.Level) Color {
	switch level {
	case slog.LevelDebug:
		return darkGray
	case slog.LevelInfo:
		return cyan
	case slog.LevelWarn:
		return lightYellow
	case slog.LevelError:
		return lightRed
	default:
		return black
	}
}
