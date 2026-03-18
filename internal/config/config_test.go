package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoadRequiresTargets(t *testing.T) {
	os.Unsetenv("INPUT_TARGETS")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error when targets is empty")
	}
}

func TestLoadValidConfig(t *testing.T) {
	t.Setenv("INPUT_TARGETS", "gitlab::https://gitlab.com/org/repo.git")
	t.Setenv("INPUT_GITLAB_TOKEN", "test-token")
	t.Setenv("INPUT_MIRROR_BRANCHES", "main,develop")
	t.Setenv("INPUT_MIRROR_TAGS", "true")
	t.Setenv("INPUT_FORCE_PUSH", "true")
	t.Setenv("INPUT_DRY_RUN", "false")
	t.Setenv("INPUT_DEBUG", "true")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Targets) != 1 {
		t.Fatalf("expected 1 target, got %d", len(cfg.Targets))
	}
	if cfg.Targets[0].Provider != ProviderGitLab {
		t.Errorf("expected gitlab provider, got %s", cfg.Targets[0].Provider)
	}
	if cfg.GitLabToken != "test-token" {
		t.Errorf("expected test-token, got %s", cfg.GitLabToken)
	}
	if cfg.MirrorAllBranches {
		t.Error("expected MirrorAllBranches to be false")
	}
	if len(cfg.MirrorBranches) != 2 {
		t.Errorf("expected 2 branches, got %d", len(cfg.MirrorBranches))
	}
	if !cfg.Debug {
		t.Error("expected debug to be true")
	}
}

func TestLoadAllBranches(t *testing.T) {
	t.Setenv("INPUT_TARGETS", "https://gitlab.com/org/repo.git")
	t.Setenv("INPUT_MIRROR_BRANCHES", "all")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.MirrorAllBranches {
		t.Error("expected MirrorAllBranches to be true")
	}
}

func TestParseTargetsMultiple(t *testing.T) {
	raw := `gitlab::https://gitlab.com/org/repo.git
codecommit::https://git-codecommit.us-east-1.amazonaws.com/v1/repos/repo
https://bitbucket.org/org/repo.git`

	targets, err := parseTargets(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 3 {
		t.Fatalf("expected 3 targets, got %d", len(targets))
	}
	if targets[0].Provider != ProviderGitLab {
		t.Errorf("target[0]: expected gitlab, got %s", targets[0].Provider)
	}
	if targets[1].Provider != ProviderCodeCommit {
		t.Errorf("target[1]: expected codecommit, got %s", targets[1].Provider)
	}
	if targets[2].Provider != ProviderBitbucket {
		t.Errorf("target[2]: expected bitbucket, got %s", targets[2].Provider)
	}
}

func TestParseTargetsEmpty(t *testing.T) {
	_, err := parseTargets("")
	if err == nil {
		t.Fatal("expected error for empty targets")
	}
}

func TestDetectProvider(t *testing.T) {
	tests := []struct {
		url      string
		expected Provider
	}{
		{"https://gitlab.com/org/repo.git", ProviderGitLab},
		{"https://github.com/org/repo.git", ProviderGitHub},
		{"https://bitbucket.org/org/repo.git", ProviderBitbucket},
		{"https://git-codecommit.us-east-1.amazonaws.com/v1/repos/repo", ProviderCodeCommit},
		{"https://custom-git.example.com/repo.git", ProviderGeneric},
	}

	for _, tt := range tests {
		got := detectProvider(tt.url)
		if got != tt.expected {
			t.Errorf("detectProvider(%q) = %s, want %s", tt.url, got, tt.expected)
		}
	}
}

func TestEnvBool(t *testing.T) {
	t.Setenv("TEST_BOOL_TRUE", "true")
	t.Setenv("TEST_BOOL_FALSE", "false")
	t.Setenv("TEST_BOOL_YES", "yes")

	if !envBool("TEST_BOOL_TRUE", false) {
		t.Error("expected true")
	}
	if envBool("TEST_BOOL_FALSE", true) {
		t.Error("expected false")
	}
	if !envBool("TEST_BOOL_YES", false) {
		t.Error("expected true for 'yes'")
	}
	if !envBool("NONEXISTENT", true) {
		t.Error("expected default true")
	}
}

func TestParseTargetsUnknownProvider(t *testing.T) {
	_, err := parseTargets("typo::https://example.com/repo.git")
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
	if !strings.Contains(err.Error(), "unknown provider") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseTargetsInvalidURL(t *testing.T) {
	_, err := parseTargets("gitlab::not-a-valid-url")
	if err == nil {
		t.Fatal("expected error for invalid URL format")
	}
	if !strings.Contains(err.Error(), "invalid URL format") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseTargetsSSHURL(t *testing.T) {
	targets, err := parseTargets("git@github.com:org/repo.git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 1 {
		t.Fatalf("expected 1 target, got %d", len(targets))
	}
	if targets[0].Provider != ProviderGitHub {
		t.Errorf("expected github provider, got %s", targets[0].Provider)
	}
}

func TestValidateEmptyBranches(t *testing.T) {
	cfg := &Config{
		MirrorAllBranches: false,
		MirrorBranches:    nil,
		Targets:           []Target{{Provider: ProviderGeneric, URL: "https://example.com/repo.git"}},
	}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for empty branches with MirrorAllBranches=false")
	}
	if !strings.Contains(err.Error(), "mirror_branches") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateWithBranches(t *testing.T) {
	cfg := &Config{
		MirrorAllBranches: false,
		MirrorBranches:    []string{"main"},
		Targets:           []Target{{Provider: ProviderGeneric, URL: "https://example.com/repo.git"}},
	}
	err := cfg.Validate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateAllBranches(t *testing.T) {
	cfg := &Config{
		MirrorAllBranches: true,
		Targets:           []Target{{Provider: ProviderGeneric, URL: "https://example.com/repo.git"}},
	}
	err := cfg.Validate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateWarnsNoToken(t *testing.T) {
	// Warnings are logged to stderr but don't return error
	cfg := &Config{
		MirrorAllBranches: true,
		Targets: []Target{
			{Provider: ProviderGitLab, URL: "https://gitlab.com/org/repo.git"},
			{Provider: ProviderGitHub, URL: "https://github.com/org/repo.git"},
			{Provider: ProviderBitbucket, URL: "https://bitbucket.org/org/repo.git"},
			{Provider: ProviderCodeCommit, URL: "https://git-codecommit.us-east-1.amazonaws.com/v1/repos/repo"},
		},
	}
	err := cfg.Validate()
	if err != nil {
		t.Fatalf("expected no error (only warnings): %v", err)
	}
}

func TestLoadEmptyBranchesFails(t *testing.T) {
	t.Setenv("INPUT_TARGETS", "gitlab::https://gitlab.com/org/repo.git")
	t.Setenv("INPUT_GITLAB_TOKEN", "token")
	t.Setenv("INPUT_MIRROR_BRANCHES", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for empty branches")
	}
}

func TestParseTargetsEmptyProvider(t *testing.T) {
	_, err := parseTargets("::https://example.com/repo.git")
	if err == nil {
		t.Fatal("expected error for empty provider")
	}
}

func TestParseTargetsValidProviders(t *testing.T) {
	tests := []struct {
		input    string
		provider Provider
	}{
		{"gitlab::https://gitlab.com/repo.git", ProviderGitLab},
		{"github::https://github.com/repo.git", ProviderGitHub},
		{"bitbucket::https://bitbucket.org/repo.git", ProviderBitbucket},
		{"codecommit::https://git-codecommit.us-east-1.amazonaws.com/v1/repos/repo", ProviderCodeCommit},
		{"generic::https://example.com/repo.git", ProviderGeneric},
	}

	for _, tt := range tests {
		targets, err := parseTargets(tt.input)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tt.input, err)
		}
		if targets[0].Provider != tt.provider {
			t.Errorf("input %q: expected provider %s, got %s", tt.input, tt.provider, targets[0].Provider)
		}
	}
}

func TestEnvInt(t *testing.T) {
	t.Setenv("TEST_INT_VALID", "3")
	t.Setenv("TEST_INT_ZERO", "0")
	t.Setenv("TEST_INT_INVALID", "abc")
	t.Setenv("TEST_INT_EMPTY", "")

	if got := envInt("TEST_INT_VALID", 0); got != 3 {
		t.Errorf("expected 3, got %d", got)
	}
	if got := envInt("TEST_INT_ZERO", 5); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
	if got := envInt("TEST_INT_INVALID", 7); got != 7 {
		t.Errorf("expected default 7 for invalid, got %d", got)
	}
	if got := envInt("TEST_INT_EMPTY", 10); got != 10 {
		t.Errorf("expected default 10 for empty, got %d", got)
	}
	if got := envInt("NONEXISTENT_INT", 42); got != 42 {
		t.Errorf("expected default 42 for nonexistent, got %d", got)
	}
}

func TestLoadWithRetryAndExclude(t *testing.T) {
	t.Setenv("INPUT_TARGETS", "gitlab::https://gitlab.com/org/repo.git")
	t.Setenv("INPUT_GITLAB_TOKEN", "token")
	t.Setenv("INPUT_MIRROR_BRANCHES", "all")
	t.Setenv("INPUT_RETRY_COUNT", "3")
	t.Setenv("INPUT_RETRY_DELAY", "10")
	t.Setenv("INPUT_EXCLUDE_BRANCHES", "staging, hotfix")
	t.Setenv("INPUT_PARALLEL", "true")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.RetryCount != 3 {
		t.Errorf("expected RetryCount 3, got %d", cfg.RetryCount)
	}
	if cfg.RetryDelay != 10 {
		t.Errorf("expected RetryDelay 10, got %d", cfg.RetryDelay)
	}
	if len(cfg.ExcludeBranches) != 2 {
		t.Fatalf("expected 2 exclude branches, got %d", len(cfg.ExcludeBranches))
	}
	if cfg.ExcludeBranches[0] != "staging" || cfg.ExcludeBranches[1] != "hotfix" {
		t.Errorf("unexpected exclude branches: %v", cfg.ExcludeBranches)
	}
	if !cfg.Parallel {
		t.Error("expected Parallel to be true")
	}
}

func TestLoadDefaultRetryValues(t *testing.T) {
	t.Setenv("INPUT_TARGETS", "gitlab::https://gitlab.com/org/repo.git")
	t.Setenv("INPUT_GITLAB_TOKEN", "token")
	t.Setenv("INPUT_MIRROR_BRANCHES", "all")
	os.Unsetenv("INPUT_RETRY_COUNT")
	os.Unsetenv("INPUT_RETRY_DELAY")
	os.Unsetenv("INPUT_EXCLUDE_BRANCHES")
	os.Unsetenv("INPUT_PARALLEL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.RetryCount != 0 {
		t.Errorf("expected default RetryCount 0, got %d", cfg.RetryCount)
	}
	if cfg.RetryDelay != 5 {
		t.Errorf("expected default RetryDelay 5, got %d", cfg.RetryDelay)
	}
	if len(cfg.ExcludeBranches) != 0 {
		t.Errorf("expected no exclude branches, got %v", cfg.ExcludeBranches)
	}
	if cfg.Parallel {
		t.Error("expected Parallel to be false by default")
	}
}
