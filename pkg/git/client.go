package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Client represents a Git client that interacts with a remote repository
type Client struct {
	RepoURL    string
	Branch     string
	WorkDir    string
	cloned     bool
	lastCommit string
}

// NewClient creates a new Git client
func NewClient(repoURL, branch string) *Client {
	if branch == "" {
		branch = "main" // Default to main branch
	}

	return &Client{
		RepoURL: repoURL,
		Branch:  branch,
		WorkDir: filepath.Join(os.TempDir(), "configsync-"+randString(8)),
		cloned:  false,
	}
}

// SyncRepository ensures the Git repository is cloned or updated
func (c *Client) SyncRepository() error {
	if c.cloned {
		return c.pullRepository()
	}
	return c.cloneRepository()
}

// cloneRepository clones the Git repository
func (c *Client) cloneRepository() error {
	// Create working directory
	if err := os.MkdirAll(c.WorkDir, 0755); err != nil {
		return fmt.Errorf("failed to create working directory: %w", err)
	}

	// Clone the repository
	cmd := exec.Command("git", "clone", "--branch", c.Branch, c.RepoURL, c.WorkDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone failed: %w, output: %s", err, string(output))
	}

	c.cloned = true
	return nil
}

// pullRepository updates the repository with the latest changes
func (c *Client) pullRepository() error {
	cmd := exec.Command("git", "-C", c.WorkDir, "pull", "origin", c.Branch)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git pull failed: %w, output: %s", err, string(output))
	}
	return nil
}

// HasChanges checks if there are changes in the Git repository since the last commit ID
func (c *Client) HasChanges(lastCommitID string) (bool, string, error) {
	if !c.cloned {
		return false, "", errors.New("repository not cloned yet")
	}

	// Get the latest commit ID
	cmd := exec.Command("git", "-C", c.WorkDir, "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return false, "", fmt.Errorf("git rev-parse failed: %w", err)
	}

	currentCommitID := strings.TrimSpace(string(output))

	// If no last commit ID is provided, assume there are changes
	if lastCommitID == "" {
		return true, currentCommitID, nil
	}

	// Check if the commit IDs are different
	if currentCommitID != lastCommitID {
		return true, currentCommitID, nil
	}

	return false, currentCommitID, nil
}

// GetDirectoryContent returns the content of a directory in the repository
func (c *Client) GetDirectoryContent(dirPath string) ([]string, error) {
	fullPath := filepath.Join(c.WorkDir, dirPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory %s does not exist in the repository", dirPath)
	}

	files, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}

	return filenames, nil
}

// GetFileContent returns the content of a file from the repository
func (c *Client) GetFileContent(filePath string) ([]byte, error) {
	fullPath := filepath.Join(c.WorkDir, filePath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return content, nil
}

// Cleanup removes the local copy of the repository
func (c *Client) Cleanup() error {
	if c.WorkDir != "" {
		return os.RemoveAll(c.WorkDir)
	}
	return nil
}

// Helper function to generate a random string
func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
