package util

func BoldTerminalText(text string) string {
	return "\033[1m" + text + "\033[0m"
}
