package config

import (
	"fmt"
	"os"
	"strings"
)

// Provider represents a git hosting provider.
type Provider string

const (
	ProviderGitLab     Provider = "gitlab"
	ProviderGitHub     Provider = "github"
	ProviderBitbucket  Provider = "bitbucket"
	ProviderCodeCommit Provider = "codecommit"
	ProviderGeneric    Provider = "generic"
)

// Target represents a single mirror target.
type Target struct {
	Provider Provider
	URL      string
}

// Config holds all configuration for the mirror action.
type Config struct {
	Targets            []Target
	GitLabToken        string
	GitHubToken        string
	BitbucketUsername   string
	BitbucketPassword  string
	SSHPrivateKey      string
	MirrorBranches     []string
	MirrorAllBranches  bool
	MirrorTags         bool
	ForcePush          bool
	DryRun             bool
	Debug              bool
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	targetsRaw := os.Getenv("INPUT_TARGETS")
	if targetsRaw == "" {
		return nil, fmt.Errorf("targets input is required")
	}

	targets, err := parseTargets(targetsRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse targets: %w", err)
	}

	branches := os.Getenv("INPUT_MIRROR_BRANCHES")
	mirrorAll := strings.TrimSpace(strings.ToLower(branches)) == "all"

	var branchList []string
	if !mirrorAll && branches != "" {
		for _, b := range strings.Split(branches, ",") {
			if trimmed := strings.TrimSpace(b); trimmed != "" {
				branchList = append(branchList, trimmed)
			}
		}
	}

	return &Config{
		Targets:           targets,
		GitLabToken:       os.Getenv("INPUT_GITLAB_TOKEN"),
		GitHubToken:       os.Getenv("INPUT_GITHUB_TOKEN"),
		BitbucketUsername:  os.Getenv("INPUT_BITBUCKET_USERNAME"),
		BitbucketPassword:  os.Getenv("INPUT_BITBUCKET_PASSWORD"),
		SSHPrivateKey:     os.Getenv("INPUT_SSH_PRIVATE_KEY"),
		MirrorBranches:    branchList,
		MirrorAllBranches: mirrorAll,
		MirrorTags:        envBool("INPUT_MIRROR_TAGS", true),
		ForcePush:         envBool("INPUT_FORCE_PUSH", true),
		DryRun:            envBool("INPUT_DRY_RUN", false),
		Debug:             envBool("INPUT_DEBUG", false),
	}, nil
}

// parseTargets parses the newline-separated targets input.
// Format: "provider::url" or just "url" (auto-detect provider).
func parseTargets(raw string) ([]Target, error) {
	var targets []Target

	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var t Target
		if parts := strings.SplitN(line, "::", 2); len(parts) == 2 {
			t.Provider = Provider(strings.ToLower(strings.TrimSpace(parts[0])))
			t.URL = strings.TrimSpace(parts[1])
		} else {
			t.URL = line
			t.Provider = detectProvider(line)
		}

		if t.URL == "" {
			return nil, fmt.Errorf("empty URL in target: %q", line)
		}

		targets = append(targets, t)
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("no valid targets found")
	}

	return targets, nil
}

// detectProvider auto-detects the provider from the URL.
func detectProvider(url string) Provider {
	lower := strings.ToLower(url)
	switch {
	case strings.Contains(lower, "gitlab.com") || strings.Contains(lower, "gitlab"):
		return ProviderGitLab
	case strings.Contains(lower, "github.com") || strings.Contains(lower, "github"):
		return ProviderGitHub
	case strings.Contains(lower, "bitbucket.org") || strings.Contains(lower, "bitbucket"):
		return ProviderBitbucket
	case strings.Contains(lower, "codecommit"):
		return ProviderCodeCommit
	default:
		return ProviderGeneric
	}
}

func envBool(key string, defaultVal bool) bool {
	val := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	switch val {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	default:
		return defaultVal
	}
}
