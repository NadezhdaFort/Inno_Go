package main

import (
	"fmt"
	"os"
	"time"
)

type EffectiveLogger struct {
	buf []byte
}

func NewEffectiveLogger() EffectiveLogger {
	return EffectiveLogger{
		buf: make([]byte, 0, 256),
	}
}

func (l EffectiveLogger) Info(msg string) {
	l.buf = l.buf[:0]
	l.buf = append(time.Now().AppendFormat(l.buf, "2006-01-02 15:04:05"))
	l.buf = append(l.buf, ' ')
	l.buf = append(l.buf, msg...)
	l.buf = append(l.buf, '\n')
	_, err := os.Stdout.Write(l.buf)
	if err != nil {
		return
	}
}

// ------------------------------------------------------------------------

type IneffectiveLogger struct {
}

func NewIneffectiveLogger() *IneffectiveLogger {
	return &IneffectiveLogger{}
}

func (l *IneffectiveLogger) Info(msg string) {
	t := time.Now()
	str := fmt.Sprintf("%s %s", t.Format("2006-01-02 15:04:05"), msg)
	fmt.Println(str)
}

func main() {
	il := NewIneffectiveLogger()
	el := NewEffectiveLogger()

	il.Info("Сообщение 1")
	el.Info("Сообщение 2")
}
