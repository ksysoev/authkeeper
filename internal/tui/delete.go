package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ksysoev/authkeeper/internal/vault"
)

type DeleteModel struct {
	vault         *vault.Vault
	passwordInput textinput.Model
	clients       []string
	selectedIdx   int
	step          int
	password      string
	confirmDelete bool
	Err           error
	spinnerTick   int
	deleted       bool
}

const (
	deleteStepPassword = iota
	deleteStepLoading
	deleteStepSelect
	deleteStepConfirm
	deleteStepDeleting
	deleteStepDone
)

type deleteCompleteMsg struct {
	err error
}

func NewDeleteModel(v *vault.Vault) DeleteModel {
	input := textinput.New()
	input.Placeholder = "Enter password"
	input.EchoMode = textinput.EchoPassword
	input.EchoCharacter = '•'
	input.CharLimit = 200
	input.Width = 50
	input.Focus()

	return DeleteModel{
		vault:         v,
		passwordInput: input,
		step:          deleteStepPassword,
	}
}

func (m DeleteModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m DeleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.step == deleteStepConfirm {
				m.step = deleteStepSelect
				return m, nil
			}
			return m, tea.Quit

		case "enter":
			return m.handleEnter()

		case "up", "k":
			if m.step == deleteStepSelect && m.selectedIdx > 0 {
				m.selectedIdx--
			}

		case "down", "j":
			if m.step == deleteStepSelect && m.selectedIdx < len(m.clients)-1 {
				m.selectedIdx++
			}

		case "y", "Y":
			if m.step == deleteStepConfirm {
				m.step = deleteStepDeleting
				return m, tea.Batch(m.deleteClient(), tickCmd())
			}

		case "n", "N":
			if m.step == deleteStepConfirm {
				m.step = deleteStepSelect
				return m, nil
			}

		case "q":
			if m.step == deleteStepDone {
				return m, tea.Quit
			}
		}

	case loadClientsMsg:
		if msg.err != nil {
			m.Err = msg.err
			return m, tea.Quit
		}
		m.clients = msg.clients
		if len(m.clients) == 0 {
			m.Err = fmt.Errorf("no clients found in vault")
			return m, tea.Quit
		}
		m.step = deleteStepSelect
		return m, nil

	case deleteCompleteMsg:
		if msg.err != nil {
			m.Err = msg.err
			m.step = deleteStepSelect
			return m, nil
		}
		m.deleted = true
		m.step = deleteStepDone
		return m, nil

	case tickMsg:
		m.spinnerTick++
		if m.step == deleteStepLoading || m.step == deleteStepDeleting {
			return m, tickCmd()
		}
	}

	var cmd tea.Cmd
	m.passwordInput, cmd = m.passwordInput.Update(msg)
	return m, cmd
}

func (m DeleteModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case deleteStepPassword:
		m.password = m.passwordInput.Value()
		m.step = deleteStepLoading
		return m, tea.Batch(m.loadClients(), tickCmd())

	case deleteStepSelect:
		if len(m.clients) > 0 {
			m.step = deleteStepConfirm
			return m, nil
		}

	case deleteStepDone:
		return m, tea.Quit
	}

	return m, nil
}

func (m *DeleteModel) loadClients() tea.Cmd {
	return func() tea.Msg {
		clients, err := m.vault.ListClients(m.password)
		return loadClientsMsg{clients: clients, err: err}
	}
}

func (m *DeleteModel) deleteClient() tea.Cmd {
	clientName := m.clients[m.selectedIdx]
	return func() tea.Msg {
		err := m.vault.DeleteClient(clientName, m.password)
		return deleteCompleteMsg{err: err}
	}
}

func (m DeleteModel) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("🗑️  Delete OIDC Client"))
	s.WriteString("\n")

	switch m.step {
	case deleteStepPassword:
		s.WriteString(subtitleStyle.Render("Enter master password"))
		s.WriteString("\n\n")
		s.WriteString(labelStyle.Render("Master Password:"))
		s.WriteString("\n")
		s.WriteString(focusedInputStyle.Render(m.passwordInput.View()))
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("Press Enter to continue • Esc to quit"))

	case deleteStepLoading:
		s.WriteString("\n")
		s.WriteString(spinnerStyle.Render(fmt.Sprintf("%s Loading clients from vault...", getSpinnerFrame(m.spinnerTick))))
		s.WriteString("\n")

	case deleteStepSelect:
		s.WriteString(subtitleStyle.Render(fmt.Sprintf("Select client to delete (%d available)", len(m.clients))))
		s.WriteString("\n\n")

		for i, client := range m.clients {
			if i == m.selectedIdx {
				s.WriteString(selectedItemStyle.Render("▶ " + client))
			} else {
				s.WriteString(listItemStyle.Render("  " + client))
			}
			s.WriteString("\n")
		}

		s.WriteString("\n")
		s.WriteString(helpStyle.Render("↑/↓: Navigate • Enter: Select • Esc: Quit"))

	case deleteStepConfirm:
		s.WriteString(subtitleStyle.Render("⚠️  Confirm deletion"))
		s.WriteString("\n\n")
		s.WriteString(errorStyle.Render(fmt.Sprintf("Are you sure you want to delete '%s'?", m.clients[m.selectedIdx])))
		s.WriteString("\n\n")
		s.WriteString(mutedStyle.Render("This action cannot be undone."))
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("Press Y to confirm • N or Esc to cancel"))

	case deleteStepDeleting:
		s.WriteString("\n")
		s.WriteString(spinnerStyle.Render(fmt.Sprintf("%s Deleting client...", getSpinnerFrame(m.spinnerTick))))
		s.WriteString("\n")

	case deleteStepDone:
		if m.deleted {
			s.WriteString("\n")
			s.WriteString(successStyle.Render("✓ Client deleted successfully!"))
			s.WriteString("\n\n")
			s.WriteString(mutedStyle.Render("Press Enter or q to exit"))
		}
	}

	if m.Err != nil {
		s.WriteString("\n\n")
		s.WriteString(errorStyle.Render("✗ Error: " + m.Err.Error()))
		s.WriteString("\n")
		s.WriteString(helpStyle.Render("Press Esc to quit"))
	}

	return s.String()
}
