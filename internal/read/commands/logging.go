package commands

import (
	"io"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/cli"
)

type cliWriter struct {
	ui    cli.Ui
	level hclog.Level
}

func newCliWriter(ui cli.Ui, level hclog.Level) io.Writer {
	return &cliWriter{
		ui:    ui,
		level: level,
	}
}

func (w *cliWriter) Write(p []byte) (n int, err error) {
	msg := string(p)
	msg = strings.Trim(msg, "\n")
	switch w.level {
	case hclog.Error:
		w.ui.Error(msg)
	case hclog.Warn:
		w.ui.Warn(msg)
	case hclog.Info, hclog.NoLevel:
		w.ui.Info(msg)
	case hclog.Debug, hclog.Trace:
		w.ui.Output(msg)
	default:
		// suppress log output
	}

	return len(p), nil
}

func InitLogging(ui cli.Ui, level hclog.Level) {
	hclog.SetDefault(hclog.New(&hclog.LoggerOptions{
		Level: level,
		Output: hclog.NewLeveledWriter(newCliWriter(ui, hclog.NoLevel), map[hclog.Level]io.Writer{
			hclog.Trace: newCliWriter(ui, hclog.Trace),
			hclog.Debug: newCliWriter(ui, hclog.Debug),
			hclog.Info:  newCliWriter(ui, hclog.Info),
			hclog.Warn:  newCliWriter(ui, hclog.Warn),
			hclog.Error: newCliWriter(ui, hclog.Error),
			hclog.Off:   newCliWriter(ui, hclog.Off),
		}),
		Color:       hclog.ColorOff,
		DisableTime: false,
	}))
}
