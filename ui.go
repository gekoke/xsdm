package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type updateTimeMsg struct {
	time time.Time
}

type performLoginMsg struct{}

type messages struct {
	infos  []string
	errors []string
}

type model struct {
	hostname string
	time     time.Time

	usernameInput textinput.Model
	passwordInput textinput.Model

	showingPassword bool

	messages     messages
	awaitingAuth bool

	width  int
	height int
}

func initialModel() model {
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("couldn't get hostname: %s", err)
	}
	var model = model{
		hostname:      hostname,
		time:          time.Now(),
		passwordInput: textinput.New(),
		usernameInput: textinput.New(),
	}
	model.usernameInput.Focus()
	return model
}

func (model model) Init() tea.Cmd {
	return nil
}

func (model model) swapFocus() model {
	if model.passwordInput.Focused() {
		model.usernameInput.Focus()
		model.passwordInput.Blur()
	} else {
		model.passwordInput.Focus()
		model.usernameInput.Blur()
	}
	return model
}

func (model model) doEnter() (model, tea.Cmd) {
	if model.usernameInput.Value() == "" {
		model.messages.infos = append(model.messages.infos, "Please enter a username")
		return model, nil
	}

	if !model.passwordInput.Focused() {
		return model.swapFocus(), nil
	} else {
		model.awaitingAuth = true
		return model, func() tea.Msg { return performLoginMsg{} }
	}
}

func (model model) performLogin() model {
	pamMessages := pamMessages{}
	err := authenticate(model.usernameInput.Value(), model.passwordInput.Value(), &pamMessages)

	if err.Error() != "Success" {
		model.messages.errors = append(model.messages.errors, err.Error())
		model.awaitingAuth = false
		return model
	}

	_ = tea.Quit()

	err = login()
	if err != nil {
		log.Fatalf("exec failed: %s", err.Error())
	}

	// Since we `exec` or exit fatally
	panic("unreachable")
}

func (model model) toggleShowPassword() model {
	model.showingPassword = !model.showingPassword
	return model
}

func (model model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model.messages.infos = []string{}
	model.messages.errors = []string{}

	switch msg := msg.(type) {
	case updateTimeMsg:
		model.time = msg.time
	case performLoginMsg:
		return model.performLogin(), nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return model, tea.Quit
		case tea.KeyTab, tea.KeyShiftTab:
			return model.swapFocus(), nil
		case tea.KeyEnter:
			return model.doEnter()
		case tea.KeyCtrlR:
			return model.toggleShowPassword(), nil
		default:
			if model.usernameInput.Focused() {
				model.usernameInput, _ = model.usernameInput.Update(msg)
			}
			if model.passwordInput.Focused() {
				model.passwordInput, _ = model.passwordInput.Update(msg)
			}
		}
	case tea.WindowSizeMsg:
		model.width = msg.Width
		model.height = msg.Height
	}
	return model, nil
}

func inputFieldStyle() lipgloss.Style {
	return lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).Width(60)
}

func hostnameTextStyle() lipgloss.Style {
	return lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#999999"))
}

func infoTextStyle() lipgloss.Style {
	return lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#5599FF"))
}

func errorTextStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#EE1111"))
}

func loadingTextStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#22EE22"))
}

func buildLogs(model model) string {
	var logs string
	if model.awaitingAuth {
		logs += loadingTextStyle().Render("Authenticating...")
	}
	if len(model.messages.infos) > 0 {
		if logs != "" {
			logs += "\n"
		}
		logs += infoTextStyle().Render(strings.Join(model.messages.infos, "\n"))
	}
	if len(model.messages.errors) > 0 {
		if logs != "" {
			logs += "\n"
		}
		logs += errorTextStyle().Render(strings.Join(model.messages.errors, "\n"))
	}
	return logs
}

func (model model) View() string {
	if !model.showingPassword {
		model.passwordInput.EchoMode = textinput.EchoPassword
		model.passwordInput.EchoCharacter = '•'
		model.passwordInput.Prompt = "> "
	} else {
		model.passwordInput.EchoMode = textinput.EchoPassword
		model.passwordInput.EchoMode = textinput.EchoNormal
		model.passwordInput.Prompt = errorTextStyle().Render("! ")
	}
	if model.passwordInput.Focused() {
		model.passwordInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9944BB"))
		model.usernameInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	} else {
		model.usernameInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9944BB"))
		model.passwordInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	}

	return lipgloss.Place(
		model.width,
		model.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			hostnameTextStyle().Render(fmt.Sprintf("💻 %s", model.hostname)),
			hostnameTextStyle().Render(fmt.Sprintf("🕙 %s", model.time.Local().Format(time.DateTime))),
			inputFieldStyle().Render(model.usernameInput.View()),
			inputFieldStyle().Render(model.passwordInput.View()),
			buildLogs(model),
			help.New().View(keys),
		),
	)
}
