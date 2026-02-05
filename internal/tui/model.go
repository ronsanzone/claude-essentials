package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rsanzone/clawdbay/internal/tmux"
)

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
func InitialModel() Model {
	return Model{
		Groups: []RepoGroup{},
		Cursor: 0,
		Nodes:  []TreeNode{},
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
			if len(m.Nodes) > 0 && m.Cursor < len(m.Nodes)-1 {
				m.Cursor++
			}
		case "enter":
			// Would attach to session
			return m, tea.Quit
		}
	}
	return m, nil
}
