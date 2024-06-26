package messages

import "github.com/jedib0t/go-pretty/v6/text"

func SuccessMessage(message string) string {
	return text.BgHiGreen.Sprintf("[success]: %s\n", message)
}

func ErrorMessage(err error) string {
	return text.BgRed.Sprintf("[error]: %s\n", err)
}

func WarningMessage(message string) string {
	return text.BgHiYellow.Sprintf("[warning]: %s\n", message)
}
