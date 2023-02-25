package main

import "fmt"

const (
	redCode    = 1
	greenCode  = 2
	yellowCode = 3
)

var (
	red    = fmt.Sprintf("\x1b[38;5;%dm", redCode)
	green  = fmt.Sprintf("\x1b[38;5;%dm", greenCode)
	yellow = fmt.Sprintf("\x1b[38;5;%dm", yellowCode)
	reset  = "\x1b[0;0;0m"
)
