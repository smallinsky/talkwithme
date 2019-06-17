package vt

import (
	"fmt"
)

// Control sequnce for vt100 terminal
// Reference: http://man7.org/linux/man-pages/man4/console_codes.4.html
var (
	ESP         string = "\x1b"     // starc an espace sequence
	CIS         string = "\x1b\x5b" // ESP + [
	EndColorSeq string = "\x1b\x5b0m"
)

func MoveCurDown(n int) string  { return fmt.Sprintf("%s%d%s", CIS, n, "B") }
func MoveCurUp(n int) string    { return fmt.Sprintf("%s%d%s", CIS, n, "A") }
func MoveCurRight(n int) string { return fmt.Sprintf("%s%d%s", CIS, n, "C") }
func MoveCurLeft(n int) string  { return fmt.Sprintf("%s%d%s", CIS, n, "D") }
func InsertLine(n int) string   { return fmt.Sprintf("%s%d%s", CIS, n, "L") }
func DeleteNLines(n int) string { return fmt.Sprintf("%s%d%s", CIS, n, "M") }
func SaveCurPos() string        { return fmt.Sprintf("%s%s", ESP, "7") }
func RestoreCurPos() string     { return fmt.Sprintf("%s%s", ESP, "8") }
func EraseLine() string         { return fmt.Sprintf("\r%s%s", CIS, "2K") }
func ClearLineCurLeft() string  { return fmt.Sprintf("\r%s%s", CIS, "1K") }
func Red(s string) string       { return fmt.Sprintf("%s%s%s%s", CIS, "31m", s, EndColorSeq) }
func Yellow(s string) string    { return fmt.Sprintf("%s%s%s%s", CIS, "33m", s, EndColorSeq) }
func Magenta(s string) string   { return fmt.Sprintf("%s%s%s%s", CIS, "35m", s, EndColorSeq) }
func Cyan(s string) string      { return fmt.Sprintf("%s%s%s%s", CIS, "36m", s, EndColorSeq) }
