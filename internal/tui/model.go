package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rsanzone/clawdbay/internal/tmux"
)

const refreshInterval = 3 * time.Second

// tickMsg triggers periodic refresh.
type tickMsg time.Time

// refreshMsg carries new data from a refresh.
type refreshMsg struct {
	Groups         []RepoGroup
	WindowStatuses map[string]tmux.Status
}

// NodeType represents what kind of tree node the cursor is on.
type NodeType int

const (
	// NodeRepo is a repository group node.
	NodeRepo NodeType = iota
	// NodeSession is a worktree session node.
	NodeSession
	// NodeWindow is a tmux window node.
	NodeWindow
)

// RepoGroup represents a repository with its worktree sessions.
type RepoGroup struct {
	Name     string
	Path     string
	Sessions []WorktreeSession
	Expanded bool
}

// WorktreeSession represents a tmux session tied to a worktree.
type WorktreeSession struct {
	Name     string
	Status   tmux.Status
	Windows  []tmux.Window
	Expanded bool
}

// TreeNode represents a flattened position in the tree for cursor navigation.
type TreeNode struct {
	Type         NodeType
	RepoIndex    int
	SessionIndex int
	WindowIndex  int
}

// Model is the Bubbletea model for the dashboard.
type Model struct {
	Groups         []RepoGroup
	Cursor         int
	Nodes          []TreeNode
	Quitting       bool
	TmuxClient     *tmux.Client
	SelectedName   string
	SelectedWindow string
	WindowStatuses map[string]tmux.Status
}

// RollupStatus returns the most active status from a slice.
// Priority: WORKING > IDLE > DONE
func RollupStatus(statuses []tmux.Status) tmux.Status {
	hasIdle := false
	for _, s := range statuses {
		if s == tmux.StatusWorking {
			return tmux.StatusWorking
		}
		if s == tmux.StatusIdle {
			hasIdle = true
		}
	}
	if hasIdle {
		return tmux.StatusIdle
	}
	return tmux.StatusDone
}

// GroupByRepo groups sessions by their repository name.
func GroupByRepo(
	sessions []tmux.Session,
	repoNames map[string]string,
	windows map[string][]tmux.Window,
	statuses map[string]tmux.Status,
) []RepoGroup {
	repoMap := make(map[string]*RepoGroup)
	var repoOrder []string

	for _, session := range sessions {
		repoName := repoNames[session.Name]
		if repoName == "" {
			repoName = "Unknown"
		}

		if _, exists := repoMap[repoName]; !exists {
			repoMap[repoName] = &RepoGroup{
				Name:     repoName,
				Expanded: true,
			}
			repoOrder = append(repoOrder, repoName)
		}

		wins := windows[session.Name]
		var windowStatuses []tmux.Status
		for _, w := range wins {
			key := session.Name + ":" + w.Name
			if status, ok := statuses[key]; ok {
				windowStatuses = append(windowStatuses, status)
			}
		}

		ws := WorktreeSession{
			Name:     session.Name,
			Status:   RollupStatus(windowStatuses),
			Windows:  wins,
			Expanded: true,
		}

		repoMap[repoName].Sessions = append(repoMap[repoName].Sessions, ws)
	}

	var groups []RepoGroup
	for _, name := range repoOrder {
		groups = append(groups, *repoMap[name])
	}
	return groups
}

// BuildNodes flattens the tree into a list of navigable nodes.
func BuildNodes(groups []RepoGroup) []TreeNode {
	var nodes []TreeNode

	for ri, repo := range groups {
		nodes = append(nodes, TreeNode{Type: NodeRepo, RepoIndex: ri})

		if !repo.Expanded {
			continue
		}

		for si, session := range repo.Sessions {
			nodes = append(nodes, TreeNode{Type: NodeSession, RepoIndex: ri, SessionIndex: si})

			if !session.Expanded {
				continue
			}

			for wi := range session.Windows {
				nodes = append(nodes, TreeNode{Type: NodeWindow, RepoIndex: ri, SessionIndex: si, WindowIndex: wi})
			}
		}
	}

	return nodes
}

// InitialModel creates the initial dashboard model.
func InitialModel(tmuxClient *tmux.Client) Model {
	return Model{
		Groups:         []RepoGroup{},
		TmuxClient:     tmuxClient,
		WindowStatuses: make(map[string]tmux.Status),
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.refreshCmd(), m.tickCmd())
}

func (m Model) tickCmd() tea.Cmd {
	return tea.Tick(refreshInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) refreshCmd() tea.Cmd {
	return func() tea.Msg {
		groups, statuses := fetchGroups(m.TmuxClient)
		return refreshMsg{Groups: groups, WindowStatuses: statuses}
	}
}

// fetchGroups queries tmux for all data.
func fetchGroups(tmuxClient *tmux.Client) ([]RepoGroup, map[string]tmux.Status) {
	if tmuxClient == nil {
		return nil, nil
	}

	sessions, err := tmuxClient.ListSessions()
	if err != nil {
		return nil, nil
	}

	repoNames := make(map[string]string)
	windowMap := make(map[string][]tmux.Window)
	statusMap := make(map[string]tmux.Status)

	for _, s := range sessions {
		repoNames[s.Name] = tmuxClient.GetRepoName(s.Name)

		wins, winErr := tmuxClient.ListWindows(s.Name)
		if winErr != nil {
			continue
		}
		windowMap[s.Name] = wins

		for _, w := range wins {
			if strings.HasPrefix(w.Name, "claude") {
				statusMap[s.Name+":"+w.Name] = tmuxClient.GetPaneStatus(s.Name, w.Name)
			}
		}
	}

	return GroupByRepo(sessions, repoNames, windowMap, statusMap), statusMap
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case refreshMsg:
		m.Groups = mergeExpandState(m.Groups, msg.Groups)
		m.WindowStatuses = msg.WindowStatuses
		m.Nodes = BuildNodes(m.Groups)
		if m.Cursor >= len(m.Nodes) {
			m.Cursor = max(0, len(m.Nodes)-1)
		}
		return m, nil

	case tickMsg:
		return m, tea.Batch(m.refreshCmd(), m.tickCmd())

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Nodes)-1 {
				m.Cursor++
			}
		case "enter":
			return m.handleEnter()
		case "l", "right":
			return m.handleExpand()
		case "h", "left":
			return m.handleCollapse()
		case "r":
			return m, m.refreshCmd()
		}
	}
	return m, nil
}

// mergeExpandState preserves expand/collapse state across refreshes.
func mergeExpandState(old, updated []RepoGroup) []RepoGroup {
	oldState := make(map[string]bool)
	oldSessionState := make(map[string]bool)

	for _, g := range old {
		oldState[g.Name] = g.Expanded
		for _, s := range g.Sessions {
			oldSessionState[s.Name] = s.Expanded
		}
	}

	for i := range updated {
		if expanded, ok := oldState[updated[i].Name]; ok {
			updated[i].Expanded = expanded
		}
		for j := range updated[i].Sessions {
			if expanded, ok := oldSessionState[updated[i].Sessions[j].Name]; ok {
				updated[i].Sessions[j].Expanded = expanded
			}
		}
	}
	return updated
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	if m.Cursor >= len(m.Nodes) {
		return m, nil
	}
	node := m.Nodes[m.Cursor]

	switch node.Type {
	case NodeRepo:
		m.Groups[node.RepoIndex].Expanded = !m.Groups[node.RepoIndex].Expanded
		m.Nodes = BuildNodes(m.Groups)
	case NodeSession:
		session := m.Groups[node.RepoIndex].Sessions[node.SessionIndex]
		m.SelectedName = session.Name
		return m, tea.Quit
	case NodeWindow:
		session := m.Groups[node.RepoIndex].Sessions[node.SessionIndex]
		window := session.Windows[node.WindowIndex]
		m.SelectedName = session.Name
		m.SelectedWindow = window.Name
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleExpand() (tea.Model, tea.Cmd) {
	if m.Cursor >= len(m.Nodes) {
		return m, nil
	}
	node := m.Nodes[m.Cursor]

	switch node.Type {
	case NodeRepo:
		m.Groups[node.RepoIndex].Expanded = true
		m.Nodes = BuildNodes(m.Groups)
	case NodeSession:
		m.Groups[node.RepoIndex].Sessions[node.SessionIndex].Expanded = true
		m.Nodes = BuildNodes(m.Groups)
	}
	return m, nil
}

func (m Model) handleCollapse() (tea.Model, tea.Cmd) {
	if m.Cursor >= len(m.Nodes) {
		return m, nil
	}
	node := m.Nodes[m.Cursor]

	switch node.Type {
	case NodeRepo:
		m.Groups[node.RepoIndex].Expanded = false
		m.Nodes = BuildNodes(m.Groups)
	case NodeSession:
		m.Groups[node.RepoIndex].Sessions[node.SessionIndex].Expanded = false
		m.Nodes = BuildNodes(m.Groups)
	case NodeWindow:
		// Collapse parent session
		m.Groups[node.RepoIndex].Sessions[node.SessionIndex].Expanded = false
		m.Nodes = BuildNodes(m.Groups)
	}
	return m, nil
}
