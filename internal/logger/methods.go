package logger

import (
	goctx "context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/smithy-go/logging"
	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/env"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var (
	sqlFilter = regexp.MustCompile(`'([^']{4,})'`)
)

// Base methods.
func (l *Writer) Write(p []byte) (n int, err error) {
	if l.filelogger != nil {
		l.filelogger.Println(string(p))
	}

	if l.syslogger != nil {
		return l.syslogger.Write(p)
	}

	if l.cloudwatchStream != nil {
		return l.cloudwatchStream.Write(p)
	}

	return len(p), nil
}

func (l *Writer) Fatal(v ...interface{}) {
	l.Write([]byte(fmt.Sprint(v...)))
	log.Fatal(v...)
}

func (l *Writer) Print(v ...interface{}) {
	l.Write([]byte(fmt.Sprint(v...)))
}

func (l *Writer) Printf(format string, v ...interface{}) {
	l.Write([]byte(fmt.Sprintf(format, v...)))
}

func (l *Writer) Println(v ...interface{}) {
	l.Write([]byte(fmt.Sprintln(v...)))
}

func (l *Writer) Log(v ...interface{}) {
	l.Write([]byte(fmt.Sprint(v...)))
}

func (l *Writer) Logf(classification logging.Classification, format string, v ...interface{}) {
	l.Write([]byte(fmt.Sprintf(format, v...)))
}

// Context logger methods.

// Gorm methods.

func (l *Writer) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.level = int(level) + 3
	return &newlogger
}

func (l *Writer) Error(ctx goctx.Context, format string, v ...interface{}) {
	l.Errorf(ctx, format, v...)
}

func (l *Writer) Warn(ctx goctx.Context, format string, v ...interface{}) {
	l.Warnf(ctx, format, v...)
}

func (l *Writer) Info(ctx goctx.Context, format string, v ...interface{}) {
	l.Infof(ctx, format, v...)
}

func (l *Writer) Debug(ctx goctx.Context, format string, v ...interface{}) {
	l.Debugf(ctx, format, v...)
}

func (l *Writer) Trace(ctx goctx.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	// Filter out sensitive information from SQL queries.
	// This is a hacky way to do it, since GORM doesn't provide a way to do this.
	if env.IsProd() {
		sqlFiltered := sqlFilter.ReplaceAllFunc([]byte(sql), func(match []byte) []byte {
			// Technically this can't happen since the regex
			// only matches strings with at least 4 characters.
			if len(match) < 4 {
				return match
			}
			return []byte(fmt.Sprintf("'%s***'", match[:4]))
		})
		sql = string(sqlFiltered)
	}

	ctx = context.WithValue(ctx, context.KeyDBQuery, strings.ReplaceAll(sql, `"`, `\"`))
	ctx = context.WithValue(ctx, context.KeyDBRows, rows)
	ctx = context.WithValue(ctx, context.KeyDBElapsed, elapsed)
	ctx = context.WithValue(ctx, context.KeyDBFile, utils.FileWithLineNum())

	if err != nil {
		ctx = context.WithValue(ctx, context.KeyError, err)
		l.Errorf(ctx, "")
		return
	}

	l.Infof(ctx, "")
}

// Formatted methods.

// Fatalf logs a message at fatal level. Does not take context since gorm
// expects a logger with this signature.
func (l *Writer) Fatalf(format string, v ...interface{}) {
	l.Fatal(fmt.Sprintf(format, v...))
}

// Emergf logs a message at emergency level.
func (l *Writer) Emerf(ctx context.Ctx, format string, v ...interface{}) {
	if l.level < LevelEmerg {
		return
	}

	ctx = context.WithValue(ctx, context.KeyLogLevel, LevelEmerg)

	if l.mode == ModeAWS {
		l.Printf(context.JSON(ctx, format, v...))
		return
	}

	ctx = context.WithValue(ctx, context.KeyLogLevel, LevelEmerg)
	l.Printf(context.String(ctx, format, v...))
}

// Errorf logs a message at error level.
func (l *Writer) Errorf(ctx context.Ctx, format string, v ...interface{}) {
	if l.level < LevelError {
		return
	}

	if l.mode == ModeAWS {
		l.Printf(context.JSON(ctx, format, v...))
		return
	}

	ctx = context.WithValue(ctx, context.KeyLogLevel, LevelError)
	l.Printf(context.String(ctx, format, v...))
}

// Warngf logs a message at warning level.
func (l *Writer) Warnf(ctx context.Ctx, format string, v ...interface{}) {
	if l.level < LevelWarng {
		return
	}

	if l.mode == ModeAWS {
		l.Printf(context.JSON(ctx, format, v...))
		return
	}

	ctx = context.WithValue(ctx, context.KeyLogLevel, LevelWarng)
	l.Printf(context.String(ctx, format, v...))
}

// Inforf logs a message at info level.
func (l *Writer) Infof(ctx context.Ctx, format string, v ...interface{}) {
	if l.level < LevelInfo {
		return
	}

	if l.mode == ModeAWS {
		l.Printf(context.JSON(ctx, format, v...))
		return
	}

	ctx = context.WithValue(ctx, context.KeyLogLevel, LevelInfo)
	l.Printf(context.String(ctx, format, v...))
}

// Debugf logs a message at debug level.
func (l *Writer) Debugf(ctx context.Ctx, format string, v ...interface{}) {
	if l.level < LevelDebug {
		return
	}

	if l.mode == ModeAWS {
		l.Printf(context.JSON(ctx, format, v...))
		return
	}

	ctx = context.WithValue(ctx, context.KeyLogLevel, LevelDebug)
	l.Printf(context.String(ctx, format, v...))
}
