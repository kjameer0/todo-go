package ansi

const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Dim    = "\033[2m"
	Italic = "\033[3m"
	Underline = "\033[4m"
	Blink  = "\033[5m"
	Reverse = "\033[7m"
	Hidden  = "\033[8m"
	Strikethrough = "\033[9m"

	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"

	Up    = "\033[A"
	Down  = "\033[B"
	Right = "\033[C"
	Left  = "\033[D"
	Backspace = "\033[3~"
	Delete = "\033[3~"
	Home   = "\033[H"
	End    = "\033[F"
	PageUp   = "\033[5~"
	PageDown = "\033[6~"
	Insert   = "\033[2~"

	CursorSave   = "\033[s"
	CursorRestore = "\033[u"
	CursorHide   = "\033[?25l"
	CursorShow   = "\033[?25h"

	ClearScreen   = "\033[2J"
	ClearLine     = "\033[K"
	ClearToEnd    = "\033[J"

	AltBufferOn  = "\033[?1049h"
	AltBufferOff = "\033[?1049l"

	Fg256        = "\033[38;5;"
	Bg256        = "\033[48;5;"
	CtrlC = 3
	CtrlD = 4
)
