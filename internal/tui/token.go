package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ksysoev/authkeeper/internal/oauth"
	"github.com/ksysoev/authkeeper/internal/vault"
)

type TokenModel struct {
	vault        *vault.Vault
	passwordInput textinput.Model
	clients      []string
	selectedIdx  int
	step         int
	password     string
	token        *oauth.TokenResponse
	client       *vault.Client
	Err          error
	spinnerTick  int
}

const (
	tokenStepPassword = iota
	tokenStepLoading
	tokenStepSelect
	tokenStepFetching
	tokenStepDisplay
)

type loadClientsMsg struct {
	clients []string
	err     error
}

type tokenFetchedMsg struct {
	token *oauth.TokenResponse
	err   error
}

type tickMsg time.Time

func NewTokenModel(v *vault.Vault) TokenModel {
	input := textinput.New()
	input.Placeholder = "Enter password"
	input.EchoMode = textinput.EchoPassword
	input.EchoCharacter = '•'
	input.CharLimit = 200
	input.Width = 50
	input.Focus()

	return TokenModel{
		vault:         v,
		passwordInput: input,
		step:          tokenStepPassword,
	}
}

func (m TokenModel) Init() tea.Cmd {
	return textinput.Blink
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m TokenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.step == tokenStepDisplay {
				return m, tea.Quit
			}
			return m, tea.Quit

		case "enter":
			return m.handleEnter()

		case "up", "k":
			if m.step == tokenStepSelect && m.selectedIdx > 0 {
				m.selectedIdx--
			}

		case "down", "j":
			if m.step == tokenStepSelect && m.selectedIdx < len(m.clients)-1 {
				m.selectedIdx++
			}

		case "q":
			if m.step == tokenStepDisplay {
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
		m.step = tokenStepSelect
		return m, nil

	case tokenFetchedMsg:
		if msg.err != nil {
			m.Err = msg.err
			m.step = tokenStepSelect
			return m, nil
		}
		m.token = msg.token
		m.step = tokenStepDisplay
		return m, nil

	case tickMsg:
		m.spinnerTick++
		if m.step == tokenStepLoading || m.step == tokenStepFetching {
			return m, tickCmd()
		}
	}

	var cmd tea.Cmd
	m.passwordInput, cmd = m.passwordInput.Update(msg)
	return m, cmd
}

func (m TokenModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case tokenStepPassword:
		m.password = m.passwordInput.Value()
		m.step = tokenStepLoading
		return m, tea.Batch(m.loadClients(), tickCmd())

	case tokenStepSelect:
		if len(m.clients) > 0 {
			m.step = tokenStepFetching
			return m, tea.Batch(m.fetchToken(), tickCmd())
		}

	case tokenStepDisplay:
		return m, tea.Quit
	}

	return m, nil
}

func (m *TokenModel) loadClients() tea.Cmd {
	return func() tea.Msg {
		clients, err := m.vault.ListClients(m.password)
		return loadClientsMsg{clients: clients, err: err}
	}
}

func (m *TokenModel) fetchToken() tea.Cmd {
	return func() tea.Msg {
		client, err := m.vault.GetClient(m.clients[m.selectedIdx], m.password)
		if err != nil {
			return tokenFetchedMsg{err: err}
		}

		m.client = client
		oauthClient := oauth.NewClient()
		token, err := oauthClient.GetToken(context.Background(), client.TokenURL, client.ClientID, client.ClientSecret, client.Scopes)

		return tokenFetchedMsg{token: token, err: err}
	}
}

func (m TokenModel) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("🎫 Issue Access Token"))
	s.WriteString("\n")

	switch m.step {
	case tokenStepPassword:
		s.WriteString(subtitleStyle.Render("Enter master password"))
		s.WriteString("\n\n")
		s.WriteString(labelStyle.Render("Master Password:"))
		s.WriteString("\n")
		s.WriteString(focusedInputStyle.Render(m.passwordInput.View()))
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("Press Enter to continue • Esc to quit"))

	case tokenStepLoading:
		s.WriteString("\n")
		s.WriteString(spinnerStyle.Render(fmt.Sprintf("%s Loading clients from vault...", getSpinnerFrame(m.spinnerTick))))
		s.WriteString("\n")

	case tokenStepSelect:
		s.WriteString(subtitleStyle.Render(fmt.Sprintf("Select client (%d available)", len(m.clients))))
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

	case tokenStepFetching:
		s.WriteString("\n")
		s.WriteString(spinnerStyle.Render(fmt.Sprintf("%s Fetching access token...", getSpinnerFrame(m.spinnerTick))))
		s.WriteString("\n")

	case tokenStepDisplay:
		s.WriteString(successStyle.Render("✓ Token issued successfully!"))
		s.WriteString("\n\n")

		if m.client != nil {
			s.WriteString(labelStyle.Render("Client: "))
			s.WriteString(m.client.Name)
			s.WriteString("\n")
		}

		if m.token != nil {
			details := boxStyle.Render(fmt.Sprintf(
				"%s\n\n%s\n\n%s\n%s\n%s",
				labelStyle.Render("Access Token:"),
				m.token.AccessToken,
				labelStyle.Render("Token Type: ")+m.token.TokenType,
				labelStyle.Render("Expires In: ")+fmt.Sprintf("%d seconds", m.token.ExpiresIn),
				labelStyle.Render("Scope: ")+m.token.Scope,
			))
			s.WriteString(details)
		}

		s.WriteString("\n\n")
		s.WriteString(mutedStyle.Render("💡 Tip: Copy the access token to use in your API requests"))
		s.WriteString("\n")
		s.WriteString(helpStyle.Render("Press Enter or q to exit"))
	}

	if m.Err != nil {
		s.WriteString("\n\n")
		s.WriteString(errorStyle.Render("✗ Error: " + m.Err.Error()))
		s.WriteString("\n")
		s.WriteString(helpStyle.Render("Press Esc to quit"))
	}

	return s.String()
}
