package logger

import (
	"log"
	"os"
)

var (
	reset  = "\033[0m"
	bold   = "\033[1m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	purple = "\033[35m"
	cyan   = "\033[36m"
	gray   = "\033[37m"
	white  = "\033[97m"

	outfile, _ = os.Create("gobulk.log")
	l          = log.New(outfile, "", 0)
)

func LogWarning(s ...any) { l.Println(yellowc("[Warning]"), s) }
func LogError(s ...any)   { l.Println(redc("[Error]"), s) }
func LogFatal(s ...any)   { l.Fatal(redc("[Fatal] "), s) }
func LogSuccess(s ...any) { l.Println(greenc("[Success]"), s) }
func LogDebug(s ...any)   { l.Println(bluec("[Debug]"), s) }
func LogNormal(s ...any)  { l.Println("[Log]", s) }

func Warning(s ...any) {
	log.Println(yellowc("[Warning]"), s)
	l.Println(yellowc("[Warning]"), s)
}

func Error(s ...any) {
	log.Println(redc("[Error]"), s)
	l.Println(redc("[Error]"), s)
}

func Fatal(s ...any) {
	log.Fatal(redc("[Fatal] "), s)
	l.Fatal(redc("[Fatal] "), s)
}

func Success(s ...any) {
	log.Println(greenc("[Success]"), s)
	l.Println(greenc("[Success]"), s)
}

func Debug(s ...any) {
	log.Println(bluec("[Debug]"), s)
	l.Println(bluec("[Debug]"), s)
}

func Normal(s ...any) {
	log.Println("[Log]", s)
	l.Println("[Log]", s)
}

func colorize(color, s string) string { return color + s + reset }
func boldc(s string) string           { return colorize(bold, s) }
func redc(s string) string            { return colorize(red, s) }
func greenc(s string) string          { return colorize(green, s) }
func yellowc(s string) string         { return colorize(yellow, s) }
func bluec(s string) string           { return colorize(blue, s) }
func purplec(s string) string         { return colorize(purple, s) }
func cyanc(s string) string           { return colorize(cyan, s) }
func grayc(s string) string           { return colorize(gray, s) }
func whitec(s string) string          { return colorize(white, s) }
