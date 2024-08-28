package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const lastFieldIdx = 1

type authInfo struct {
	infos  []string
	errors []string
}

type model struct {
	authInfo        authInfo
	awaitingAuth    bool
	hostname        string
	focusIndex      int
	username        string
	password        string
	width           int
	height          int
	showingPassword bool
}

func initialModel() model {
	hostname, err := os.Hostname()
	if err != nil {
		panic(fmt.Sprintf("couldn't get hostname: %s", err))
	}
	return model{
		hostname: hostname,
	}
}

func (model model) isUsernameFieldSelected() bool {
	return model.focusIndex == 0
}

func (model model) isPasswordFieldSelected() bool {
	return model.focusIndex == 1
}

func (model model) Init() tea.Cmd {
	return nil
}

func (model model) focusPreviousField(cycle bool) model {
	model.focusIndex--
	if model.focusIndex < 0 {
		if cycle {
			model.focusIndex = lastFieldIdx
		} else {
			model.focusIndex = 0
		}
	}
	return model
}

func (model model) focusNextField(cycle bool) model {
	model.focusIndex++
	if model.focusIndex > lastFieldIdx {
		if cycle {
			model.focusIndex = 0
		} else {
			model.focusIndex = lastFieldIdx
		}
	}
	return model
}

func typeInto(value string, input tea.KeyMsg) string {
	textInput := textinput.New()
	textInput.SetValue(value)
	textInput.Focus()
	model, _ := textInput.Update(input)
	return model.Value()
}

type performLogin struct{}

func (model model) doEnter() (model, tea.Cmd) {
	if model.username == "" {
		model.authInfo.infos = append(model.authInfo.infos, "Please enter a username")
		return model, nil
	}

	if !model.isPasswordFieldSelected() {
		return model.focusNextField(false), nil
	}

	model.awaitingAuth = true
	return model, func() tea.Msg { return performLogin{} }
}

func (model model) performLogin() model {
	err := login(model.username, model.password, &model.authInfo)
	if err != nil {
		model.authInfo.errors = append(model.authInfo.errors, err.Error())
	} else {
		panic("login was successful, but the rest is not yet implemented")
	}
	model.awaitingAuth = false
	return model
}

func (model model) toggleShowPassword() model {
	model.showingPassword = !model.showingPassword
	return model
}

func (model model) doType(msg tea.KeyMsg) model {
	if model.isUsernameFieldSelected() {
		model.username = typeInto(model.username, msg)
	} else if model.isPasswordFieldSelected() {
		model.password = typeInto(model.password, msg)
	}
	return model
}

func (model model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model.authInfo.infos = []string{}
	model.authInfo.errors = []string{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.width = msg.Width
		model.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return model, tea.Quit
		case tea.KeyTab:
			return model.focusNextField(true), nil
		case tea.KeyShiftTab:
			return model.focusPreviousField(true), nil
		case tea.KeyEnter:
			return model.doEnter()
		case tea.KeyCtrlR:
			return model.toggleShowPassword(), nil
		default:
			return model.doType(msg), nil
		}
	case performLogin:
		return model.performLogin(), nil
	}
	return model, nil
}

func setAppropriateFocus(index int, input *textinput.Model, model model) {
	if model.focusIndex == index {
		input.Focus()
		input.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9944BB"))
	} else {
		input.Blur()
		input.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	}
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
	if len(model.authInfo.infos) > 0 {
		if logs != "" {
			logs += "\n"
		}
		logs += infoTextStyle().Render(strings.Join(model.authInfo.infos, "\n"))
	}
	if len(model.authInfo.errors) > 0 {
		if logs != "" {
			logs += "\n"
		}
		logs += errorTextStyle().Render(strings.Join(model.authInfo.errors, "\n"))
	}
	return logs
}

func (model model) View() string {
	usernameInput := textinput.New()
	usernameInput.SetValue(model.username)
	usernameInput.Placeholder = "Username"
	setAppropriateFocus(0, &usernameInput, model)

	passwordInput := textinput.New()
	passwordInput.SetValue(model.password)
	passwordInput.Placeholder = "Password"
	setAppropriateFocus(1, &passwordInput, model)
	if !model.showingPassword {
		passwordInput.EchoMode = textinput.EchoPassword
		passwordInput.EchoCharacter = '•'
		passwordInput.Prompt = "> "
	} else {
		passwordInput.EchoMode = textinput.EchoPassword
		passwordInput.EchoMode = textinput.EchoNormal
		passwordInput.Prompt = errorTextStyle().Render("! ")
	}

	return lipgloss.Place(
		model.width,
		model.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			hostnameTextStyle().Render(fmt.Sprintf("💻 %s", model.hostname)),
			inputFieldStyle().Render(usernameInput.View()),
			inputFieldStyle().Render(passwordInput.View()),
			buildLogs(model),
			help.New().View(keys),
		),
	)
}
