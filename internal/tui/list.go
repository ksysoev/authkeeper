package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ksysoev/authkeeper/internal/vault"
)

type ListModel struct {
	vault         *vault.Vault
	passwordInput textinput.Model
	clients       []string
	vaultData     *vault.VaultData
	step          int
	password      string
	Err           error
	spinnerTick   int
}

const (
	listStepPassword = iota
	listStepLoading
	listStepDisplay
)

func NewListModel(v *vault.Vault) ListModel {
	input := textinput.New()
	input.Placeholder = "Enter password"
	input.EchoMode = textinput.EchoPassword
	input.EchoCharacter = '•'
	input.CharLimit = 200
	input.Width = 50
	input.Focus()

	return ListModel{
		vault:         v,
		passwordInput: input,
		step:          listStepPassword,
	}
}

func (m ListModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit

		case "enter":
			if m.step == listStepPassword {
				m.password = m.passwordInput.Value()
				m.step = listStepLoading
				return m, tea.Batch(m.loadVault(), tickCmd())
			} else if m.step == listStepDisplay {
				return m, tea.Quit
			}
		}

	case loadVaultMsg:
		if msg.err != nil {
			m.Err = msg.err
			return m, tea.Quit
		}
		m.vaultData = msg.data
		m.step = listStepDisplay
		return m, nil

	case tickMsg:
		m.spinnerTick++
		if m.step == listStepLoading {
			return m, tickCmd()
		}
	}

	var cmd tea.Cmd
	m.passwordInput, cmd = m.passwordInput.Update(msg)
	return m, cmd
}

type loadVaultMsg struct {
	data *vault.VaultData
	err  error
}

func (m *ListModel) loadVault() tea.Cmd {
	return func() tea.Msg {
		data, err := m.vault.Load(m.password)
		return loadVaultMsg{data: data, err: err}
	}
}

func (m ListModel) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("📋 OIDC Clients"))
	s.WriteString("\n")

	switch m.step {
	case listStepPassword:
		s.WriteString(subtitleStyle.Render("Enter master password"))
		s.WriteString("\n\n")
		s.WriteString(labelStyle.Render("Master Password:"))
		s.WriteString("\n")
		s.WriteString(focusedInputStyle.Render(m.passwordInput.View()))
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("Press Enter to continue • Esc to quit"))

	case listStepLoading:
		s.WriteString("\n")
		s.WriteString(spinnerStyle.Render(fmt.Sprintf("%s Loading vault...", getSpinnerFrame(m.spinnerTick))))
		s.WriteString("\n")

	case listStepDisplay:
		if m.vaultData == nil || len(m.vaultData.Clients) == 0 {
			s.WriteString(subtitleStyle.Render("No clients found"))
			s.WriteString("\n\n")
			s.WriteString(mutedStyle.Render("Use 'authkeeper add' to add your first client"))
		} else {
			s.WriteString(subtitleStyle.Render(fmt.Sprintf("%d client(s) in vault", len(m.vaultData.Clients))))
			s.WriteString("\n\n")

			for _, client := range m.vaultData.Clients {
				clientBox := boxStyle.Render(fmt.Sprintf(
					"%s\n%s\n%s\n%s\n%s",
					labelStyle.Render("Name: ")+client.Name,
					labelStyle.Render("Client ID: ")+client.ClientID,
					labelStyle.Render("Token URL: ")+client.TokenURL,
					labelStyle.Render("Scopes: ")+strings.Join(client.Scopes, ", "),
					mutedStyle.Render("Created: "+client.CreatedAt.Format("2006-01-02 15:04:05")),
				))
				s.WriteString(clientBox)
				s.WriteString("\n")
			}
		}

		s.WriteString("\n")
		s.WriteString(helpStyle.Render("Press q or Esc to exit"))
	}

	if m.Err != nil {
		s.WriteString("\n\n")
		s.WriteString(errorStyle.Render("✗ Error: " + m.Err.Error()))
	}

	return s.String()
}
