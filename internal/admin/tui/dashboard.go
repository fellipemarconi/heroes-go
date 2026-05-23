package tui

import (
	"api/internal/admin/client"
	"api/internal/module/admin"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true)

	statsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))
)

type DashboardModel struct {
	menu   []string
	cursor int

	users   int32
	heroes  int32
	images  int32
	storage int64

	loading bool
	err     string
}

func NewDashboard() DashboardModel {
	return DashboardModel{
		menu: []string{
			"Dashboard",
			"Users",
			"Containers",
			"Logs",
		},

		storage: -1,
		loading: true,
	}
}

func (m DashboardModel) Init() tea.Cmd {
	return fetchStatsCmd()
}

type statsMsg struct {
	stats *admin.Stats
}

type statsErrMsg struct {
	err error
}

func fetchStatsCmd() tea.Cmd {
	return func() tea.Msg {
		stats, err := client.GetStats()
		if err != nil {
			return statsErrMsg{err: err}
		}
		return statsMsg{stats: stats}
	}
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statsMsg:
		m.users = msg.stats.Users
		m.heroes = msg.stats.Heroes
		m.images = msg.stats.Images
		m.storage = msg.stats.Storage
		m.loading = false
		m.err = ""
		return m, nil

	case statsErrMsg:
		m.loading = false
		m.err = msg.err.Error()
		return m, nil

	case tea.KeyMsg:

		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.menu)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m DashboardModel) View() string {
	s := titleStyle.Render("Heroes Admin") + "\n\n"

	for i, item := range m.menu {

		line := fmt.Sprintf("  %s", item)

		if m.cursor == i {
			line = selectedStyle.Render("> " + item)
		}

		s += line + "\n"
	}

	s += "\n"

	stats := "Loading stats..."
	if !m.loading {
		stats = fmt.Sprintf(
			"Users: %d\nHeroes: %d\nImages: %d\nStorage: %s",
			m.users,
			m.heroes,
			m.images,
			formatStorage(m.storage),
		)
	}

	s += statsStyle.Render(stats)

	s += "\n\n"

	if m.err != "" {
		s += errorStyle.Render("Error: " + m.err)
		s += "\n\n"
	}

	s += helpStyle.Render("[↑/↓] navigate • [q] quit")

	return s
}

func formatStorage(bytes int64) string {
	if bytes < 0 {
		return "--"
	}

	return humanize.Bytes(uint64(bytes))
}
