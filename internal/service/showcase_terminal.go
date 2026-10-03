package service

import (
	"bytes"

	"github.com/gorilla/websocket"
)

// ShowcaseLine is one keystroke buffer for the simulated exhibit shell.
type ShowcaseLine struct {
	buf string
}

// Feed echoes typed characters and returns a completed line on Enter.
func (s *ShowcaseLine) Feed(msg []byte) (echo string, line string, ok bool) {
	if bytes.HasPrefix(bytes.TrimSpace(msg), []byte(`{"type":"resize"`)) {
		return "", "", false
	}
	var echoed stringsBuilder
	for _, b := range msg {
		switch b {
		case '\r', '\n':
			line = s.buf
			s.buf = ""
			echoed.WriteString("\r\n")
			return echoed.String(), line, true
		case 127, 8:
			if s.buf != "" {
				s.buf = s.buf[:len(s.buf)-1]
				echoed.WriteString("\b \b")
			}
		default:
			if b < 32 {
				continue
			}
			s.buf += string(b)
			echoed.WriteByte(b)
		}
	}
	return echoed.String(), "", false
}

// stringsBuilder avoids importing strings just to accumulate a short echo.
type stringsBuilder struct{ b []byte }

func (s *stringsBuilder) WriteString(v string) { s.b = append(s.b, v...) }
func (s *stringsBuilder) WriteByte(v byte)     { s.b = append(s.b, v) }
func (s *stringsBuilder) String() string       { return string(s.b) }

// RunShowcaseTerminal serves a line-oriented simulated shell over the websocket.
func RunShowcaseTerminal(ws *websocket.Conn, banner, prompt string, reply func(line string) (string, bool)) {
	_ = ws.WriteMessage(websocket.TextMessage, []byte(banner+prompt))
	var line ShowcaseLine
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			return
		}
		echo, entered, ok := line.Feed(msg)
		if echo != "" {
			_ = ws.WriteMessage(websocket.TextMessage, []byte(echo))
		}
		if !ok {
			continue
		}
		out, done := reply(entered)
		if out != "" {
			_ = ws.WriteMessage(websocket.TextMessage, []byte(out))
		}
		if done {
			return
		}
		_ = ws.WriteMessage(websocket.TextMessage, []byte(prompt))
	}
}
