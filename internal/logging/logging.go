package logging

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type Logger struct {
	instance *log.Logger
}

func GetLogger() (*Logger, error) {
	styles := log.DefaultStyles()
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERROR!!").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("204")).
		Foreground(lipgloss.Color("0"))
	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styles.Values["err"] = lipgloss.NewStyle().Bold(true)
	logger := Logger{}
	logger.instance = log.New(os.Stdout)
	logger.instance.SetStyles(styles)
	return &logger, nil
}

func GetResultLogger() (*Logger, error) {
	styles := log.DefaultStyles()
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("204")).
		Foreground(lipgloss.Color("0"))
	styles.Keys["info"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styles.Values["info"] = lipgloss.NewStyle().Bold(true)
	logger := Logger{}
	logger.instance = log.NewWithOptions(os.Stdout, log.Options{
		Prefix: "",
	})
	logger.instance.SetStyles(styles)
	return &logger, nil
}

func (l Logger) Info(msg interface{}) {
	l.instance.Info(msg)
}

func (l Logger) Print(msg interface{}) {
	l.instance.Print(msg)
}

func (l Logger) Error(msg interface{}) {
	l.instance.Error(msg)
}
