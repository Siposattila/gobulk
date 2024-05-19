package console

import "log"

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
)

func colorize(color, s string) string {
	return color + s + reset
}

func Warning(s ...any) {
	log.Println(yellowc("[Warning]"), s)
}

func Error(s ...any) {
	log.Println(redc("[Error]"), s)
}

func Fatal(s ...any) {
	log.Fatal(redc("[Fatal] "), s)
}

func Success(s ...any) {
	log.Println(greenc("[Success]"), s)
}

func Debug(s ...any) {
	log.Println(bluec("[Debug]"), s)
}

func Normal(s ...any) {
	log.Println("[Log]", s)
}

func boldc(s string) string {
	return colorize(bold, s)
}

func redc(s string) string {
	return colorize(red, s)
}

func greenc(s string) string {
	return colorize(green, s)
}

func yellowc(s string) string {
	return colorize(yellow, s)
}

func bluec(s string) string {
	return colorize(blue, s)
}

func purplec(s string) string {
	return colorize(purple, s)
}

func cyanc(s string) string {
	return colorize(cyan, s)
}

func grayc(s string) string {
	return colorize(gray, s)
}

func whitec(s string) string {
	return colorize(white, s)
}
