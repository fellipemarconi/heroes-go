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

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("70"))

	confirmStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))
)

type DashboardModel struct {
	menu   []string
	cursor int

	usersCount int32
	heroes     int32
	images     int32
	storage    int64

	users        []admin.UserSummary
	userCursor   int
	usersLoading bool

	confirmDelete bool
	confirmUser   admin.UserSummary

	statsLoading bool
	err          string
	status       string
}

func NewDashboard() DashboardModel {
	return DashboardModel{
		menu: []string{
			"Users",
			"Containers",
			"Logs",
		},

		storage:      -1,
		statsLoading: true,
		usersLoading: true,
	}
}

func (m DashboardModel) Init() tea.Cmd {
	return tea.Batch(fetchStatsCmd(), fetchUsersCmd())
}

type statsMsg struct {
	stats *admin.Stats
}

type statsErrMsg struct {
	err error
}

type usersMsg struct {
	users []admin.UserSummary
}

type usersErrMsg struct {
	err error
}

type deleteUserMsg struct {
	id string
}

type deleteUserErrMsg struct {
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

func fetchUsersCmd() tea.Cmd {
	return func() tea.Msg {
		users, err := client.GetUsers()
		if err != nil {
			return usersErrMsg{err: err}
		}
		return usersMsg{users: users}
	}
}

func deleteUserCmd(userID string) tea.Cmd {
	return func() tea.Msg {
		if err := client.DeleteUser(userID); err != nil {
			return deleteUserErrMsg{err: err}
		}
		return deleteUserMsg{id: userID}
	}
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statsMsg:
		m.usersCount = msg.stats.Users
		m.heroes = msg.stats.Heroes
		m.images = msg.stats.Images
		m.storage = msg.stats.Storage
		m.statsLoading = false
		m.err = ""
		return m, nil

	case statsErrMsg:
		m.statsLoading = false
		m.err = msg.err.Error()
		return m, nil

	case usersMsg:
		m.users = msg.users
		m.usersLoading = false
		m.err = ""
		if m.userCursor >= len(m.users) {
			m.userCursor = max(0, len(m.users)-1)
		}
		return m, nil

	case usersErrMsg:
		m.usersLoading = false
		m.err = msg.err.Error()
		return m, nil

	case deleteUserMsg:
		m.status = "User deleted."
		m.err = ""
		m.usersLoading = true
		m.statsLoading = true
		m.confirmDelete = false
		return m, tea.Batch(fetchUsersCmd(), fetchStatsCmd())

	case deleteUserErrMsg:
		m.err = msg.err.Error()
		m.confirmDelete = false
		return m, nil

	case tea.KeyMsg:
		if m.confirmDelete {
			switch msg.String() {
			case "y", "Y":
				m.status = ""
				m.err = ""
				m.usersLoading = true
				m.confirmDelete = false
				return m, deleteUserCmd(m.confirmUser.ID)
			case "n", "esc":
				m.confirmDelete = false
				return m, nil
			case "q", "ctrl+c":
				return m, tea.Quit
			}
			return m, nil
		}

		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

		case "left", "h":
			if m.cursor > 0 {
				m.cursor--
			}

		case "right", "l":
			if m.cursor < len(m.menu)-1 {
				m.cursor++
			}

		case "up", "k":
			if m.menu[m.cursor] == "Users" {
				if m.userCursor > 0 {
					m.userCursor--
				}
			} else if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.menu[m.cursor] == "Users" {
				if m.userCursor < len(m.users)-1 {
					m.userCursor++
				}
			} else if m.cursor < len(m.menu)-1 {
				m.cursor++
			}

		case "d":
			if m.menu[m.cursor] == "Users" && len(m.users) > 0 {
				m.status = ""
				m.err = ""
				m.confirmDelete = true
				m.confirmUser = m.users[m.userCursor]
				return m, nil
			}

		case "r":
			m.status = ""
			m.err = ""
			m.usersLoading = true
			m.statsLoading = true
			return m, tea.Batch(fetchUsersCmd(), fetchStatsCmd())
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
	if !m.statsLoading {
		stats = fmt.Sprintf(
			"Users: %d\nHeroes: %d\nImages: %d\nStorage: %s",
			m.usersCount,
			m.heroes,
			m.images,
			formatStorage(m.storage),
		)
	}

	s += statsStyle.Render(stats)

	s += "\n\n"

	if m.menu[m.cursor] == "Users" {
		s += renderUsersList(m.users, m.userCursor, m.usersLoading)
		s += "\n\n"
	} else {
		s += helpStyle.Render("Section not available yet.")
		s += "\n\n"
	}

	if m.confirmDelete {
		s += confirmStyle.Render(
			fmt.Sprintf(
				"Delete user %s (%s)? [y]es / [n]o",
				m.confirmUser.Email,
				m.confirmUser.ID,
			),
		)
		s += "\n\n"
	}

	if m.status != "" {
		s += statusStyle.Render(m.status)
		s += "\n\n"
	}

	if m.err != "" {
		s += errorStyle.Render("Error: " + m.err)
		s += "\n\n"
	}

	if m.confirmDelete {
		s += helpStyle.Render("[y] confirm • [n] cancel • [q] quit")
	} else if m.menu[m.cursor] == "Users" {
		s += helpStyle.Render("[↑/↓] users • [←/→] menu • [d] delete • [r] refresh • [q] quit")
	} else {
		s += helpStyle.Render("[←/→] menu • [r] refresh • [q] quit")
	}

	return s
}

func formatStorage(bytes int64) string {
	if bytes < 0 {
		return "--"
	}

	return humanize.Bytes(uint64(bytes))
}

func renderUsersList(users []admin.UserSummary, cursor int, loading bool) string {
	if loading {
		return "Loading users..."
	}

	if len(users) == 0 {
		return "No users found."
	}

	header := fmt.Sprintf(
		"%-36s  %-24s  %-20s  %s",
		"ID",
		"Email",
		"Name",
		"Created At",
	)

	output := header + "\n"

	for i, user := range users {
		line := fmt.Sprintf(
			"%-36s  %-24s  %-20s  %s",
			truncate(user.ID, 36),
			truncate(user.Email, 24),
			truncate(user.Name, 20),
			user.CreatedAt.Format("2006-01-02 15:04"),
		)

		if i == cursor {
			output += selectedStyle.Render("> "+line) + "\n"
		} else {
			output += "  " + line + "\n"
		}
	}

	return output
}

func truncate(value string, max int) string {
	if len(value) <= max {
		return value
	}
	if max <= 3 {
		return value[:max]
	}
	return value[:max-3] + "..."
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
