package tui

import (
	"api/internal/admin/client"
	"api/internal/module/admin"
	"fmt"
	"strings"

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

	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("241")).
			Padding(1, 2)

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

	containers        []admin.ContainerSummary
	containersLoading bool

	logs        []admin.LogEntry
	logsLoading bool
	logCursor   int
	logDetail   bool
	searchMode  bool
	searchQuery string

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

		storage:           -1,
		statsLoading:      true,
		usersLoading:      true,
		containersLoading: true,
		logsLoading:       true,
	}
}

func (m DashboardModel) Init() tea.Cmd {
	return tea.Batch(fetchStatsCmd(), fetchUsersCmd(), fetchContainersCmd(), fetchLogsCmd())
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

type containersMsg struct {
	containers []admin.ContainerSummary
}

type containersErrMsg struct {
	err error
}

type restartContainersMsg struct{}

type restartContainersErrMsg struct {
	err error
}

type logsMsg struct {
	logs []admin.LogEntry
}

type logsErrMsg struct {
	err error
}

const logsLimit = 200

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

func fetchContainersCmd() tea.Cmd {
	return func() tea.Msg {
		containers, err := client.GetContainers()
		if err != nil {
			return containersErrMsg{err: err}
		}
		return containersMsg{containers: containers}
	}
}

func restartContainersCmd() tea.Cmd {
	return func() tea.Msg {
		if err := client.RestartContainers(); err != nil {
			return restartContainersErrMsg{err: err}
		}
		return restartContainersMsg{}
	}
}

func fetchLogsCmd() tea.Cmd {
	return func() tea.Msg {
		logs, err := client.GetLogs(logsLimit)
		if err != nil {
			return logsErrMsg{err: err}
		}
		return logsMsg{logs: logs}
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

	case containersMsg:
		m.containers = msg.containers
		m.containersLoading = false
		m.err = ""
		return m, nil

	case containersErrMsg:
		m.containersLoading = false
		m.err = msg.err.Error()
		return m, nil

	case restartContainersMsg:
		m.status = "Containers restarted."
		m.err = ""
		m.containersLoading = true
		return m, fetchContainersCmd()

	case restartContainersErrMsg:
		m.err = msg.err.Error()
		return m, nil

	case logsMsg:
		m.logs = msg.logs
		m.logsLoading = false
		m.err = ""
		m.logCursor = clampCursor(m.logCursor, len(filterLogs(m.logs, m.searchQuery)))
		return m, nil

	case logsErrMsg:
		m.logsLoading = false
		m.err = msg.err.Error()
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

		if m.searchMode {
			switch msg.Type {
			case tea.KeyEnter:
				m.searchMode = false
				m.logCursor = clampCursor(m.logCursor, len(filterLogs(m.logs, m.searchQuery)))
				return m, nil
			case tea.KeyEsc:
				m.searchMode = false
				m.searchQuery = ""
				m.logCursor = 0
				return m, nil
			case tea.KeyBackspace, tea.KeyDelete:
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					m.logCursor = clampCursor(m.logCursor, len(filterLogs(m.logs, m.searchQuery)))
				}
				return m, nil
			case tea.KeyRunes:
				m.searchQuery += string(msg.Runes)
				m.logCursor = clampCursor(m.logCursor, len(filterLogs(m.logs, m.searchQuery)))
				return m, nil
			}
			return m, nil
		}

		if m.logDetail {
			switch msg.String() {
			case "esc":
				m.logDetail = false
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
			} else if m.menu[m.cursor] == "Logs" {
				if m.logCursor > 0 {
					m.logCursor--
				}
			} else if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.menu[m.cursor] == "Users" {
				if m.userCursor < len(m.users)-1 {
					m.userCursor++
				}
			} else if m.menu[m.cursor] == "Logs" {
				if m.logCursor < len(filterLogs(m.logs, m.searchQuery))-1 {
					m.logCursor++
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
			m.containersLoading = true
			m.logsLoading = true
			return m, tea.Batch(fetchUsersCmd(), fetchStatsCmd(), fetchContainersCmd(), fetchLogsCmd())

		case "R":
			if m.menu[m.cursor] == "Containers" {
				m.status = "Restarting containers..."
				m.err = ""
				m.containersLoading = true
				return m, restartContainersCmd()
			}

		case "/":
			if m.menu[m.cursor] == "Logs" {
				m.searchMode = true
				return m, nil
			}

		case "enter":
			if m.menu[m.cursor] == "Logs" {
				if _, ok := selectedLog(m.logs, m.searchQuery, m.logCursor); ok {
					m.confirmDelete = false
					m.logDetail = true
					m.status = ""
					m.err = ""
					return m, nil
				}
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
	} else if m.menu[m.cursor] == "Containers" {
		s += renderContainersList(m.containers, m.containersLoading)
		s += "\n\n"
	} else if m.menu[m.cursor] == "Logs" {
		s += renderLogsList(m.logs, m.logCursor, m.logsLoading, m.searchQuery, m.searchMode)
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

	if m.logDetail {
		if logEntry, ok := selectedLog(m.logs, m.searchQuery, m.logCursor); ok {
			s += modalStyle.Render(renderLogDetail(logEntry))
			s += "\n\n"
		}
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
	} else if m.menu[m.cursor] == "Containers" {
		s += helpStyle.Render("[←/→] menu • [R] restart all • [r] refresh • [q] quit")
	} else if m.menu[m.cursor] == "Logs" {
		if m.searchMode {
			s += helpStyle.Render("Type to search • [enter] apply • [esc] clear")
		} else if m.logDetail {
			s += helpStyle.Render("[esc] close • [q] quit")
		} else {
			s += helpStyle.Render("[↑/↓] logs • [enter] details • [/] search • [r] refresh • [q] quit")
		}
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

func renderContainersList(containers []admin.ContainerSummary, loading bool) string {
	if loading {
		return "Loading containers..."
	}

	if len(containers) == 0 {
		return "No containers found."
	}

	header := fmt.Sprintf("%-28s  %s", "Name", "Status")
	output := header + "\n"

	for _, container := range containers {
		line := fmt.Sprintf(
			"%-28s  %s",
			truncate(container.Name, 28),
			container.Status,
		)
		output += "  " + line + "\n"
	}

	return output
}

func renderLogsList(logs []admin.LogEntry, cursor int, loading bool, query string, searchMode bool) string {
	if loading {
		return "Loading logs..."
	}

	filtered := filterLogs(logs, query)
	if len(filtered) == 0 {
		if query != "" {
			return fmt.Sprintf("Search: %s\nNo logs match the current filter.", query)
		}
		return "No logs available."
	}

	searchLine := "Search: /"
	if query != "" || searchMode {
		searchLine = fmt.Sprintf("Search: %s", query)
	}

	header := fmt.Sprintf(
		"%-19s  %-6s  %-3s  %-7s  %-20s  %-18s  %s",
		"Timestamp",
		"Method",
		"St",
		"Dur",
		"Error",
		"Location",
		"Path",
	)

	output := searchLine + "\n" + header + "\n"

	for i, entry := range filtered {
		errorText := entry.Error
		if errorText == "" {
			errorText = "-"
		}
		locationText := entry.Location
		if locationText == "" {
			locationText = "-"
		}

		line := fmt.Sprintf(
			"%-19s  %-6s  %-3d  %-7s  %-20s  %-18s  %s",
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Method,
			entry.Status,
			fmt.Sprintf("%dms", entry.DurationMs),
			truncate(errorText, 20),
			truncate(locationText, 18),
			truncate(entry.Path, 40),
		)

		if i == cursor {
			output += selectedStyle.Render("> "+line) + "\n"
		} else {
			output += "  " + line + "\n"
		}
	}

	return output
}

func renderLogDetail(entry admin.LogEntry) string {
	errorText := entry.Error
	if errorText == "" {
		errorText = "-"
	}
	locationText := entry.Location
	if locationText == "" {
		locationText = "-"
	}

	return fmt.Sprintf(
		"Timestamp: %s\nMethod: %s\nStatus: %d\nDuration: %dms\nPath: %s\nError: %s\nLocation: %s",
		entry.Timestamp.Format("2006-01-02 15:04:05"),
		entry.Method,
		entry.Status,
		entry.DurationMs,
		entry.Path,
		errorText,
		locationText,
	)
}

func filterLogs(logs []admin.LogEntry, query string) []admin.LogEntry {
	if query == "" {
		return logs
	}

	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return logs
	}

	filtered := make([]admin.LogEntry, 0, len(logs))
	for _, entry := range logs {
		if strings.Contains(strings.ToLower(entry.Method), q) ||
			strings.Contains(strings.ToLower(entry.Path), q) ||
			strings.Contains(strings.ToLower(entry.Error), q) ||
			strings.Contains(strings.ToLower(entry.Location), q) ||
			strings.Contains(strings.ToLower(fmt.Sprintf("%d", entry.Status)), q) {
			filtered = append(filtered, entry)
		}
	}

	return filtered
}

func selectedLog(logs []admin.LogEntry, query string, cursor int) (admin.LogEntry, bool) {
	filtered := filterLogs(logs, query)
	if cursor < 0 || cursor >= len(filtered) {
		return admin.LogEntry{}, false
	}

	return filtered[cursor], true
}

func clampCursor(cursor int, length int) int {
	if length <= 0 {
		return 0
	}

	if cursor < 0 {
		return 0
	}

	if cursor >= length {
		return length - 1
	}

	return cursor
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
