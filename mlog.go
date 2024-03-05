package mlog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"golang.org/x/term"
)

type Sgr string

const (
	SgrReset                 Sgr = "\x1b[0m"
	SgrBold                  Sgr = "\x1b[1m"
	SgrFaint                 Sgr = "\x1b[2m"
	SgrItalic                Sgr = "\x1b[3m"
	SgrUnderline             Sgr = "\x1b[4m"
	SgrSlowBlink             Sgr = "\x1b[5m"
	SgrRapidBlink            Sgr = "\x1b[6m"
	SgrInvert                Sgr = "\x1b[7m"
	SgrConceal               Sgr = "\x1b[8m"
	SgrCrossedOut            Sgr = "\x1b[9m"
	SgrDoublyUnderlined      Sgr = "\x1b[21m"
	SgrNoBoldOrFaint         Sgr = "\x1b[22m"
	SgrNoItalic              Sgr = "\x1b[23m"
	SgrNoUnderline           Sgr = "\x1b[24m"
	SgrNoBlink               Sgr = "\x1b[25m"
	SgrNoInvert              Sgr = "\x1b[27m"
	SgrNoConceal             Sgr = "\x1b[28m"
	SgrNoCrossedOut          Sgr = "\x1b[29m"
	SgrOverline              Sgr = "\x1b[53m"
	SgrNoOverline            Sgr = "\x1b[55m"
	SgrUnderlineColor        Sgr = "\x1b[58m"
	SgrUnderlineColorDefault Sgr = "\x1b[59m"

	SgrFgBlack         Sgr = "\x1b[30m"
	SgrFgRed           Sgr = "\x1b[31m"
	SgrFgGreen         Sgr = "\x1b[32m"
	SgrFgYellow        Sgr = "\x1b[33m"
	SgrFgBlue          Sgr = "\x1b[34m"
	SgrFgMagenta       Sgr = "\x1b[35m"
	SgrFgCyan          Sgr = "\x1b[36m"
	SgrFgWhite         Sgr = "\x1b[37m"
	SgrFgExtendedColor Sgr = "\x1b[38;5;%vm"
	SgrFgTrueColor     Sgr = "\x1b[38;2;%v;%v;%vm"
	SgrFgDefault       Sgr = "\x1b[39m"

	SgrFgBrightBlack   Sgr = "\x1b[90m"
	SgrFgBrightRed     Sgr = "\x1b[91m"
	SgrFgBrightGreen   Sgr = "\x1b[92m"
	SgrFgBrightYellow  Sgr = "\x1b[93m"
	SgrFgBrightBlue    Sgr = "\x1b[94m"
	SgrFgBrightMagenta Sgr = "\x1b[95m"
	SgrFgBrightCyan    Sgr = "\x1b[96m"
	SgrFgBrightWhite   Sgr = "\x1b[97m"

	SgrBgBlack         Sgr = "\x1b[40m"
	SgrBgRed           Sgr = "\x1b[41m"
	SgrBgGreen         Sgr = "\x1b[42m"
	SgrBgYellow        Sgr = "\x1b[43m"
	SgrBgBlue          Sgr = "\x1b[44m"
	SgrBgMagenta       Sgr = "\x1b[45m"
	SgrBgCyan          Sgr = "\x1b[46m"
	SgrBgWhite         Sgr = "\x1b[47m"
	SgrBgExtendedColor Sgr = "\x1b[48;5;%vm"
	SgrBgTrueColor     Sgr = "\x1b[48;2;%v;%v;%vm"
	SgrBgDefault       Sgr = "\x1b[49m"

	SgrBgBrightBlack   Sgr = "\x1b[100m"
	SgrBgBrightRed     Sgr = "\x1b[101m"
	SgrBgBrightGreen   Sgr = "\x1b[102m"
	SgrBgBrightYellow  Sgr = "\x1b[103m"
	SgrBgBrightBlue    Sgr = "\x1b[104m"
	SgrBgBrightMagenta Sgr = "\x1b[105m"
	SgrBgBrightCyan    Sgr = "\x1b[106m"
	SgrBgBrightWhite   Sgr = "\x1b[107m"
)

func SetSgr(sgr Sgr) {
	fmt.Printf("%v", string(sgr))
}

func ExtendedFgColor(number uint8) Sgr {
	return Sgr(fmt.Sprintf(string(SgrFgExtendedColor), number))
}

func ExtendedBgColor(number uint8) Sgr {
	return Sgr(fmt.Sprintf(string(SgrBgExtendedColor), number))
}

func TrueFgColor(r, g, b uint8) Sgr {
	return Sgr(fmt.Sprintf(string(SgrFgTrueColor), r, g, b))
}

func TrueBgColor(r, g, b uint8) Sgr {
	return Sgr(fmt.Sprintf(string(SgrBgTrueColor), r, g, b))
}

func PrintLine() {
	if !term.IsTerminal(0) {
		return
	}

	width, _, err := term.GetSize(0)

	if err != nil {
		return
	}

	line := strings.Repeat("─", width)
	SetSgr(SgrFgBrightBlack)
	fmt.Printf("%v\n", line)
	SetSgr(SgrReset)
}

type LogLevel uint8

const (
	LogLevelFatal = iota
	LogLevelError
	LogLevelWarning
	LogLevelInfo
	LogLevelDebug
)

func printInternal(level LogLevel, text string, err error, params ...any) {
	SetSgr(SgrFgBrightBlack)
	fmt.Printf("%v ", time.Now().Format("02 Jan 2006 15:04:05 MST"))
	SetSgr(SgrReset)

	switch level {
	case LogLevelFatal:
		SetSgr(SgrFgBrightRed)
		fmt.Printf("FATAL")
	case LogLevelError:
		SetSgr(SgrFgBrightRed)
		fmt.Printf("ERROR")
	case LogLevelWarning:
		SetSgr(SgrFgBrightYellow)
		fmt.Printf("WARNING")
	case LogLevelInfo:
		SetSgr(SgrFgBrightGreen)
		fmt.Printf("INFO")
	case LogLevelDebug:
		SetSgr(SgrFgBrightCyan)
		fmt.Printf("DEBUG")
	}

	SetSgr(SgrReset)

	_, filename, line, _ := runtime.Caller(2)

	filename = filepath.Base(filename)

	SetSgr(SgrFgBrightBlue)
	fmt.Printf(" %v:%v ", filename, line)
	SetSgr(SgrReset)

	fmt.Printf("%v\n", text)

	switch level {

	case LogLevelFatal,
		LogLevelError,
		LogLevelWarning:

		SetSgr(SgrFgBrightMagenta)

		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			stack := debug.Stack()
			fmt.Printf("%v", string(stack))
		}

		SetSgr(SgrReset)
	}

	for _, param := range params {
		// XXX: dumper
		fmt.Printf("%+v\n", param)
	}

	PrintLine()
}

func Fatal(text string, params ...any) {
	printInternal(LogLevelFatal, text, nil, params...)
	os.Exit(1)
}

func FatalWithErr(text string, err error, params ...any) {
	printInternal(LogLevelFatal, text, err, params...)
	os.Exit(1)
}

func Error(text string, params ...any) {
	printInternal(LogLevelError, text, nil, params...)
}

func ErrorWithErr(text string, err error, params ...any) {
	printInternal(LogLevelError, text, err, params...)
}

func Warning(text string, params ...any) {
	printInternal(LogLevelWarning, text, nil, params...)
}

func Info(text string, params ...any) {
	printInternal(LogLevelInfo, text, nil, params...)
}

func Debug(text string, params ...any) {
	printInternal(LogLevelDebug, text, nil, params...)
}

func init() {
	PrintLine()
}
