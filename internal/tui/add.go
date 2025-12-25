package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ksysoev/authkeeper/internal/vault"
)

type AddClientModel struct {
	vault    *vault.Vault
	inputs   []textinput.Model
	focused  int
	step     int
	password string
	Err      error
	success  bool
}

const (
	stepPassword = iota
	stepClientInfo
	stepConfirm
	stepSaving
	stepDone
)

const (
	inputName = iota
	inputClientID
	inputClientSecret
	inputTokenURL
	inputScopes
)

func NewAddClientModel(v *vault.Vault) AddClientModel {
	inputs := make([]textinput.Model, 5)

	inputs[inputName] = textinput.New()
	inputs[inputName].Placeholder = "My OIDC Provider"
	inputs[inputName].Focus()
	inputs[inputName].CharLimit = 100
	inputs[inputName].Width = 50

	inputs[inputClientID] = textinput.New()
	inputs[inputClientID].Placeholder = "client_id_here"
	inputs[inputClientID].CharLimit = 200
	inputs[inputClientID].Width = 50

	inputs[inputClientSecret] = textinput.New()
	inputs[inputClientSecret].Placeholder = "client_secret_here"
	inputs[inputClientSecret].EchoMode = textinput.EchoPassword
	inputs[inputClientSecret].EchoCharacter = '•'
	inputs[inputClientSecret].CharLimit = 200
	inputs[inputClientSecret].Width = 50

	inputs[inputTokenURL] = textinput.New()
	inputs[inputTokenURL].Placeholder = "https://provider.com/oauth/token"
	inputs[inputTokenURL].CharLimit = 500
	inputs[inputTokenURL].Width = 50

	inputs[inputScopes] = textinput.New()
	inputs[inputScopes].Placeholder = "read write (optional)"
	inputs[inputScopes].CharLimit = 200
	inputs[inputScopes].Width = 50

	return AddClientModel{
		vault:   v,
		inputs:  inputs,
		focused: 0,
		step:    stepPassword,
	}
}

func (m AddClientModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m AddClientModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			return m.handleEnter()

		case "tab", "down":
			if m.step == stepClientInfo {
				m.nextInput()
			}

		case "shift+tab", "up":
			if m.step == stepClientInfo {
				m.prevInput()
			}
		}

	case saveCompleteMsg:
		if msg.err != nil {
			m.Err = msg.err
			m.step = stepClientInfo
		} else {
			m.success = true
			m.step = stepDone
		}
		return m, nil
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *AddClientModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case stepPassword:
		if len(m.inputs[inputName].Value()) > 0 {
			m.password = m.inputs[inputName].Value()
			m.inputs[inputName].SetValue("")
			m.inputs[inputName].Placeholder = "My OIDC Provider"
			m.step = stepClientInfo
			return m, nil
		}

	case stepClientInfo:
		if m.validateInputs() {
			m.step = stepConfirm
			return m, nil
		}

	case stepConfirm:
		m.step = stepSaving
		return m, m.saveClient()

	case stepDone:
		return m, tea.Quit
	}

	return m, nil
}

func (m *AddClientModel) validateInputs() bool {
	return len(m.inputs[inputName].Value()) > 0 &&
		len(m.inputs[inputClientID].Value()) > 0 &&
		len(m.inputs[inputClientSecret].Value()) > 0 &&
		len(m.inputs[inputTokenURL].Value()) > 0
}

func (m *AddClientModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
	m.updateFocus()
}

func (m *AddClientModel) prevInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
	m.updateFocus()
}

func (m *AddClientModel) updateFocus() {
	for i := range m.inputs {
		if i == m.focused {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}

func (m *AddClientModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

type saveCompleteMsg struct {
	err error
}

func (m *AddClientModel) saveClient() tea.Cmd {
	return func() tea.Msg {
		scopes := []string{}
		if scopesStr := m.inputs[inputScopes].Value(); len(scopesStr) > 0 {
			scopes = strings.Fields(scopesStr)
		}

		client := vault.Client{
			Name:         m.inputs[inputName].Value(),
			ClientID:     m.inputs[inputClientID].Value(),
			ClientSecret: m.inputs[inputClientSecret].Value(),
			TokenURL:     m.inputs[inputTokenURL].Value(),
			Scopes:       scopes,
		}

		err := m.vault.AddClient(client, m.password)
		return saveCompleteMsg{err: err}
	}
}

func (m AddClientModel) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("🔐 Add OIDC Client"))
	s.WriteString("\n")

	switch m.step {
	case stepPassword:
		s.WriteString(subtitleStyle.Render("Enter master password to unlock vault"))
		s.WriteString("\n\n")
		s.WriteString(labelStyle.Render("Master Password:"))
		s.WriteString("\n")
		input := m.inputs[inputName]
		input.EchoMode = textinput.EchoPassword
		input.EchoCharacter = '•'
		input.Placeholder = "Enter password"
		s.WriteString(focusedInputStyle.Render(input.View()))
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("Press Enter to continue • Esc to quit"))

	case stepClientInfo:
		s.WriteString(subtitleStyle.Render("Enter client credentials"))
		s.WriteString("\n\n")

		labels := []string{"Client Name:", "Client ID:", "Client Secret:", "Token URL:", "Scopes (optional):"}
		for i, label := range labels {
			s.WriteString(labelStyle.Render(label))
			s.WriteString("\n")
			if i == m.focused {
				s.WriteString(focusedInputStyle.Render(m.inputs[i].View()))
			} else {
				s.WriteString(inputStyle.Render(m.inputs[i].View()))
			}
			s.WriteString("\n")
		}

		s.WriteString("\n")
		s.WriteString(helpStyle.Render("Tab/↓: Next field • Shift+Tab/↑: Previous • Enter to continue • Esc to quit"))

	case stepConfirm:
		s.WriteString(subtitleStyle.Render("Confirm client details"))
		s.WriteString("\n\n")

		details := boxStyle.Render(fmt.Sprintf(
			"%s\n%s\n%s\n%s\n%s",
			labelStyle.Render("Name: ")+m.inputs[inputName].Value(),
			labelStyle.Render("Client ID: ")+m.inputs[inputClientID].Value(),
			labelStyle.Render("Client Secret: ")+strings.Repeat("•", len(m.inputs[inputClientSecret].Value())),
			labelStyle.Render("Token URL: ")+m.inputs[inputTokenURL].Value(),
			labelStyle.Render("Scopes: ")+m.inputs[inputScopes].Value(),
		))
		s.WriteString(details)
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("Press Enter to save • Esc to cancel"))

	case stepSaving:
		s.WriteString("\n\n")
		s.WriteString(spinnerStyle.Render("⠋ Encrypting and saving to vault..."))
		s.WriteString("\n\n")

	case stepDone:
		if m.success {
			s.WriteString("\n")
			s.WriteString(successStyle.Render("✓ Client added successfully!"))
			s.WriteString("\n\n")
			s.WriteString(mutedStyle.Render("Press Enter to exit"))
		}
	}

	if m.Err != nil {
		s.WriteString("\n\n")
		s.WriteString(errorStyle.Render("✗ Error: " + m.Err.Error()))
	}

	return s.String()
}
