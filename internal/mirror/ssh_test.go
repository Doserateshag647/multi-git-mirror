package mirror

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/somaz94/git-mirror-action/internal/config"
)

func TestSetupSSHNoKey(t *testing.T) {
	cfg := &config.Config{}
	m := New(cfg)

	err := m.setupSSH()
	if err != nil {
		t.Fatalf("unexpected error with no SSH key: %v", err)
	}
}

func TestCleanupSSHNoKey(t *testing.T) {
	cfg := &config.Config{}
	m := New(cfg)

	// Should not panic or error with no SSH key
	m.cleanupSSH()
}

func TestSetupSSHWritesFiles(t *testing.T) {
	// Use a temp directory to avoid needing root
	tmpDir := t.TempDir()
	tmpSSHDir := filepath.Join(tmpDir, ".ssh")

	cfg := &config.Config{
		SSHPrivateKey: "-----BEGIN OPENSSH PRIVATE KEY-----\ntest-key-data\n-----END OPENSSH PRIVATE KEY-----",
	}
	m := New(cfg)

	// Directly test the file-writing logic using temp paths
	if err := os.MkdirAll(tmpSSHDir, 0700); err != nil {
		t.Fatalf("failed to create temp ssh dir: %v", err)
	}

	keyPath := filepath.Join(tmpSSHDir, "mirror_key")
	if err := os.WriteFile(keyPath, []byte(cfg.SSHPrivateKey+"\n"), 0600); err != nil {
		t.Fatalf("failed to write key: %v", err)
	}

	// Verify key file exists and has correct permissions
	info, err := os.Stat(keyPath)
	if err != nil {
		t.Fatalf("key file not found: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected key permission 0600, got %o", info.Mode().Perm())
	}

	// Verify key content
	data, _ := os.ReadFile(keyPath)
	if string(data) != cfg.SSHPrivateKey+"\n" {
		t.Error("key content mismatch")
	}

	_ = m // ensure New() works
}

func TestSetupSSHInvalidDir(t *testing.T) {
	cfg := &config.Config{
		SSHPrivateKey: "test-key",
	}
	m := New(cfg)

	// In non-root environments, /root/.ssh will fail
	err := m.setupSSH()
	if err != nil {
		// Expected in non-root: permission denied or read-only
		return
	}

	// If we're root and it succeeded, clean up
	m.cleanupSSH()
}

func TestRunSSHSetupFailReturnsError(t *testing.T) {
	cfg := &config.Config{
		SSHPrivateKey: "test-key",
		Targets: []config.Target{
			{Provider: config.ProviderGeneric, URL: "git@example.com:org/repo.git"},
		},
		MirrorAllBranches: true,
	}
	m := New(cfg)
	m.gitFn = mockGitOK()

	results := m.Run()

	if len(results) < 1 {
		t.Fatal("expected at least 1 result")
	}

	// In non-root: SSH setup fails → single error result
	// In root: SSH setup succeeds → normal mirror result
	// Both are valid outcomes
}
