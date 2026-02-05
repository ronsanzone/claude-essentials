package tui

import (
	"testing"

	"github.com/rsanzone/clawdbay/internal/tmux"
)

func TestGroupByRepo(t *testing.T) {
	sessions := []tmux.Session{
		{Name: "cb_feat-auth"},
		{Name: "cb_refactor"},
		{Name: "cb_fix-login"},
	}

	repoNames := map[string]string{
		"cb_feat-auth": "my-project",
		"cb_refactor":  "my-project",
		"cb_fix-login": "other-project",
	}

	windows := map[string][]tmux.Window{
		"cb_feat-auth": {
			{Index: 0, Name: "shell", Active: true},
			{Index: 1, Name: "claude", Active: false},
			{Index: 2, Name: "claude:research", Active: false},
		},
		"cb_refactor": {
			{Index: 0, Name: "shell", Active: true},
		},
		"cb_fix-login": {
			{Index: 0, Name: "shell", Active: true},
			{Index: 1, Name: "claude", Active: false},
		},
	}

	statuses := map[string]tmux.Status{
		"cb_feat-auth:claude":          tmux.StatusWorking,
		"cb_feat-auth:claude:research": tmux.StatusIdle,
		"cb_fix-login:claude":          tmux.StatusDone,
	}

	groups := GroupByRepo(sessions, repoNames, windows, statuses)

	if len(groups) != 2 {
		t.Fatalf("got %d groups, want 2", len(groups))
	}

	// Verify ordering preserved
	if groups[0].Name != "my-project" {
		t.Errorf("first group = %q, want %q", groups[0].Name, "my-project")
	}
	if groups[1].Name != "other-project" {
		t.Errorf("second group = %q, want %q", groups[1].Name, "other-project")
	}

	// my-project should have 2 sessions
	if len(groups[0].Sessions) != 2 {
		t.Errorf("my-project has %d sessions, want 2", len(groups[0].Sessions))
	}

	// Check status rollup: feat-auth has WORKING and IDLE, should roll up to WORKING
	for _, s := range groups[0].Sessions {
		if s.Name == "cb_feat-auth" {
			if s.Status != tmux.StatusWorking {
				t.Errorf("cb_feat-auth status = %q, want %q", s.Status, tmux.StatusWorking)
			}
		}
	}

	// other-project has 1 session with DONE status
	if len(groups[1].Sessions) != 1 {
		t.Errorf("other-project has %d sessions, want 1", len(groups[1].Sessions))
	}
	if groups[1].Sessions[0].Status != tmux.StatusDone {
		t.Errorf("fix-login status = %q, want %q", groups[1].Sessions[0].Status, tmux.StatusDone)
	}
}

func TestGroupByRepo_UnknownRepo(t *testing.T) {
	sessions := []tmux.Session{{Name: "cb_orphan"}}
	repoNames := map[string]string{} // empty — no repo detected
	windows := map[string][]tmux.Window{}
	statuses := map[string]tmux.Status{}

	groups := GroupByRepo(sessions, repoNames, windows, statuses)

	if len(groups) != 1 {
		t.Fatalf("got %d groups, want 1", len(groups))
	}
	if groups[0].Name != "Unknown" {
		t.Errorf("group name = %q, want %q", groups[0].Name, "Unknown")
	}
}

func TestBuildNodes(t *testing.T) {
	groups := []RepoGroup{
		{
			Name:     "my-project",
			Expanded: true,
			Sessions: []WorktreeSession{
				{
					Name:     "cb_feat-auth",
					Status:   tmux.StatusWorking,
					Expanded: true,
					Windows: []tmux.Window{
						{Index: 0, Name: "shell"},
						{Index: 1, Name: "claude"},
					},
				},
				{
					Name:     "cb_refactor",
					Status:   tmux.StatusIdle,
					Expanded: false,
					Windows:  []tmux.Window{{Index: 0, Name: "shell"}},
				},
			},
		},
		{
			Name:     "other-project",
			Expanded: false,
			Sessions: nil,
		},
	}

	nodes := BuildNodes(groups)

	// Expected:
	// 0: Repo "my-project" (expanded)
	// 1: Session "cb_feat-auth" (expanded)
	// 2: Window "shell"
	// 3: Window "claude"
	// 4: Session "cb_refactor" (collapsed — no window children)
	// 5: Repo "other-project" (collapsed — no session children)
	if len(nodes) != 6 {
		t.Fatalf("got %d nodes, want 6", len(nodes))
	}

	if nodes[0].Type != NodeRepo {
		t.Errorf("node 0 type = %v, want NodeRepo", nodes[0].Type)
	}
	if nodes[1].Type != NodeSession {
		t.Errorf("node 1 type = %v, want NodeSession", nodes[1].Type)
	}
	if nodes[2].Type != NodeWindow {
		t.Errorf("node 2 type = %v, want NodeWindow", nodes[2].Type)
	}
	if nodes[3].Type != NodeWindow {
		t.Errorf("node 3 type = %v, want NodeWindow", nodes[3].Type)
	}
	if nodes[4].Type != NodeSession {
		t.Errorf("node 4 type = %v, want NodeSession", nodes[4].Type)
	}
	if nodes[5].Type != NodeRepo {
		t.Errorf("node 5 type = %v, want NodeRepo", nodes[5].Type)
	}
}

func TestBuildNodes_AllCollapsed(t *testing.T) {
	groups := []RepoGroup{
		{Name: "repo-a", Expanded: false},
		{Name: "repo-b", Expanded: false},
	}

	nodes := BuildNodes(groups)

	if len(nodes) != 2 {
		t.Fatalf("got %d nodes, want 2", len(nodes))
	}
}

func TestStatusRollup(t *testing.T) {
	tests := []struct {
		name     string
		statuses []tmux.Status
		expected tmux.Status
	}{
		{"all idle", []tmux.Status{tmux.StatusIdle, tmux.StatusIdle}, tmux.StatusIdle},
		{"one working", []tmux.Status{tmux.StatusIdle, tmux.StatusWorking}, tmux.StatusWorking},
		{"all done", []tmux.Status{tmux.StatusDone, tmux.StatusDone}, tmux.StatusDone},
		{"mixed", []tmux.Status{tmux.StatusDone, tmux.StatusIdle, tmux.StatusWorking}, tmux.StatusWorking},
		{"empty", []tmux.Status{}, tmux.StatusDone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RollupStatus(tt.statuses)
			if result != tt.expected {
				t.Errorf("RollupStatus() = %q, want %q", result, tt.expected)
			}
		})
	}
}
