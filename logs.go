package logs

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var Log *Logger = NewLogger(Warn)

var DefaultFormatterMap = map[Level]string{
	Debug: "[dbg] %s \n",
	Warn:  "[wrn] %s \n",
	Info:  "[inf] %s {{suffix}}\n",
	Error: "[err] %s {{suffix}}\n",
}

var Levels = map[Level]string{
	Debug: "dbg",
	Warn:  "wrn",
	Info:  "inf",
	Error: "err",
}

func NewLogger(level Level) *Logger {
	log := &Logger{
		level:     level,
		writer:    os.Stdout,
		levels:    Levels,
		formatter: DefaultFormatterMap,
		SuffixFunc: func() string {
			return ", " + getCurtime()
		},
		PrefixFunc: func() string {
			return ""
		},
	}

	return log
}

const (
	Debug Level = iota
	Warn
	Info
	Error
)

type Level int

func (l Level) Name() string {
	if name, ok := Levels[l]; ok {
		return name
	} else {
		return strconv.Itoa(int(l))
	}
}

func (l Level) Formatter() string {
	if formatter, ok := DefaultFormatterMap[l]; ok {
		return formatter
	} else {
		return "[" + l.Name() + "] %s"
	}
}

type Logger struct {
	Quiet bool // is enable Print

	onLogger  func(level Level, s interface{})
	writer    io.Writer
	level     Level
	levels    map[Level]string
	formatter map[Level]string

	// can set suffix and prefix
	SuffixFunc func() string
	PrefixFunc func() string
}

func (log *Logger) SetQuiet(q bool) {
	log.Quiet = q
}

func (log *Logger) SetLevel(l Level) {
	log.level = l
}

func (log *Logger) SetOutput(w io.Writer) {
	log.writer = w
}

func (log *Logger) SetFormatter(formatter map[Level]string) {
	log.formatter = formatter
}

func (log *Logger) logInterface(w io.Writer, level Level, s interface{}) {
	if log.Quiet {
		return
	}
	if level < log.level {
		return
	}
	line := log.Format(level, s)
	fmt.Fprint(w, line)

	if log.onLogger != nil {
		log.onLogger(level, line)
	}
}

func (log *Logger) logInterfacef(w io.Writer, level Level, format string, s ...interface{}) {
	log.logInterface(w, level, fmt.Sprintf(format, s...))
}

func (log *Logger) Print(level Level, s interface{}) {
	log.logInterface(log.writer, level, s)
}

func (log *Logger) Printf(level Level, format string, s ...interface{}) {
	log.logInterfacef(log.writer, level, format, s...)
}

func (log *Logger) Println(writer io.Writer, level Level, s ...interface{}) {
	log.logInterface(writer, level, fmt.Sprintln(s...))
}

func (log *Logger) Info(s interface{}) {
	log.logInterface(log.writer, Info, s)
}

func (log *Logger) Infof(format string, s ...interface{}) {
	log.logInterfacef(log.writer, Info, format, s...)
}

func (log *Logger) FInfof(writer io.Writer, format string, s ...interface{}) {
	log.logInterfacef(writer, Info, format, s...)
}

func (log *Logger) Error(s interface{}) {
	log.logInterface(log.writer, Error, s)
}

func (log *Logger) Errorf(format string, s ...interface{}) {
	log.logInterfacef(log.writer, Error, format, s...)
}

func (log *Logger) FErrorf(writer io.Writer, format string, s ...interface{}) {
	log.logInterfacef(writer, Error, format, s...)
}

func (log *Logger) Warn(s interface{}) {
	log.logInterface(log.writer, Warn, s)
}

func (log *Logger) Warnf(format string, s ...interface{}) {
	log.logInterfacef(log.writer, Warn, format, s...)
}

func (log *Logger) FWarnf(writer io.Writer, format string, s ...interface{}) {
	log.logInterfacef(writer, Warn, format, s...)
}

func (log *Logger) Debug(s interface{}) {
	log.logInterface(log.writer, Debug, s)

}

func (log *Logger) Debugf(format string, s ...interface{}) {
	log.logInterfacef(log.writer, Debug, format, s...)
}

func (log *Logger) FDebugf(writer io.Writer, format string, s ...interface{}) {
	log.logInterfacef(writer, Debug, format, s...)
}

func (log *Logger) Format(level Level, s ...interface{}) string {
	var line string
	if f, ok := log.formatter[level]; ok {
		line = fmt.Sprintf(f, s...)
	} else if f, ok := DefaultFormatterMap[level]; ok {
		line = fmt.Sprintf(f, s...)
	} else {
		line = fmt.Sprintf("[%s] %s ", append([]interface{}{level.Name()}, s...)...)
	}
	line = strings.Replace(line, "{{suffix}}", log.SuffixFunc(), -1)
	line = strings.Replace(line, "{{prefix}}", log.PrefixFunc(), -1)
	return line
}

// 获取当前时间
func getCurtime() string {
	curtime := time.Now().Format("2006-01-02 15:04.05")
	return curtime
}
