package colour

import (
	"fmt"
	"strconv"
	"strings"
)

const escape = "\x1b"

const (
	Reset Colour = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground text colors
const (
	FgBlack Colour = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground Hi-Intensity text colors
const (
	FgHiBlack Colour = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack Colour = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity text colors
const (
	BgHiBlack Colour = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

type Colour int

type Colours []Colour

func (c Colours) Sprint(a ...any) string {
	return c.Format() + fmt.Sprint(a...) + Unformat()
}

func (c Colours) Format() string {
	return fmt.Sprintf("%s[%sm", escape, c.sequence())
}

func Unformat() string {
	return fmt.Sprintf("%s[%dm", escape, Reset)
}

func (c Colours) sequence() string {
	format := make([]string, len(c))
	for i, colour := range c {
		format[i] = strconv.Itoa(int(colour))
	}

	return strings.Join(format, ";")
}
