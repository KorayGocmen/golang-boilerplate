package logger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/syslog"
	"os"
	"time"

	"github.com/koraygocmen/golang-boilerplate/internal/aws"
	"github.com/koraygocmen/golang-boilerplate/internal/context"
)

type Config struct {
	Level  int
	Mode   string
	SHASUM string

	Syslog struct {
		Addr     string
		Protocol string
		Tag      string
	}

	File struct {
		Path string
	}

	AWS struct {
		LogGroup struct {
			Name   string
			Region string
		}
	}
}

type Writer struct {
	mode  Mode
	level int

	syslogger        *syslog.Writer
	filelogger       *log.Logger
	cloudwatchStream *aws.Stream
}

var (
	// Logger is the global logger.
	// It is initialized in cmd/api/main.go.
	// Do not let it be nil.
	Logger = &Writer{}
)

func New(config Config) (*Writer, error) {
	var (
		file             *os.File
		filelogger       *log.Logger
		syslogger        *syslog.Writer
		cloudwatchStream *aws.Stream
	)

	mode := ToMode(config.Mode)
	switch mode {
	case ModeConsole:
		{
			file = os.Stdout
		}
	case ModeFile:
		{
			logfile, err := os.Open(config.File.Path)
			if err != nil {
				err = fmt.Errorf("failed to open log file: %w", err)
				return nil, err
			}

			file = logfile
		}
	case ModeSyslog:
		{
			var (
				syslogProtocol = config.Syslog.Protocol
				syslogAddr     = config.Syslog.Addr
				syslogTag      = config.Syslog.Tag
			)

			var err error
			syslogger, err = syslog.Dial(syslogProtocol, syslogAddr, syslog.LOG_LOCAL0, syslogTag)
			if err != nil {
				err = fmt.Errorf("failed to dial syslog: %w", err)
				return nil, err
			}
		}
	case ModeAWS:
		{
			var err error
			cloudwatchStream, err = aws.Client.LogStream(
				context.Background(),
				config.AWS.LogGroup.Name,
				fmt.Sprintf("%s-%d", config.SHASUM, time.Now().Unix()),
			)
			if err != nil {
				err = fmt.Errorf("failed to create cloudwatch stream: %w", err)
				return nil, err
			}
		}
	case ModeNone:
		{
			filelogger = log.New(io.Discard, "", log.LstdFlags)
		}

	default:
		err := fmt.Errorf("unknown log mode: %s", config.Mode)
		return nil, err
	}

	if mode != ModeNone && file == nil && syslogger == nil && cloudwatchStream == nil {
		return nil, errors.New("no log output specified")
	}

	if file != nil {
		filelogger = log.New(file, "", log.LstdFlags|log.Lmicroseconds)
	}

	logger := &Writer{
		mode:             mode,
		level:            config.Level,
		syslogger:        syslogger,
		filelogger:       filelogger,
		cloudwatchStream: cloudwatchStream,
	}

	return logger, nil
}
