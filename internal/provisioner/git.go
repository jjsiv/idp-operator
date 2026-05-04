package provisioner

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/client"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/transport/ssh"
)

type Git struct {
	repository *GitRepository
	files      []GitFile
	keys       *ssh.PublicKeys
}

func (g *Git) Provision() error {
	repo, err := g.cloneRepository()
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	repoPath := w.Filesystem.Root()
	for _, file := range g.files {
		fileFullPath := repoPath + "/" + file.Path
		if err := os.MkdirAll(filepath.Dir(fileFullPath), 0755); err != nil {
			return fmt.Errorf("failed to create directories for file %s: %w", file.Path, err)
		}

		if err := os.WriteFile(fileFullPath, []byte(file.Content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", file.Path, err)
		}
	}

	if err := g.pushChanges(repo); err != nil {
		return err
	}

	return nil
}

func (g *Git) Deprovision() error {
	return nil
}

func (g *Git) cloneRepository() (*git.Repository, error) {
	repoPath := os.TempDir() + "/" + generateRepoName()
	repo, err := git.PlainClone(repoPath, &git.CloneOptions{
		URL:           g.repository.URL,
		ReferenceName: plumbing.ReferenceName(g.repository.Ref),
		Depth:         1,
		ClientOptions: []client.Option{
			client.WithSSHAuth(g.keys),
		},
	})

	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (g *Git) pushChanges(repo *git.Repository) error {
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	if err := w.AddWithOptions(&git.AddOptions{All: true}); err != nil {
		return err
	}

	_, err = w.Commit("Provisioning resources", &git.CommitOptions{
		Author: &object.Signature{
			Name:  g.repository.CommitAuthor.Name,
			Email: g.repository.CommitAuthor.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	repo.Push(&git.PushOptions{
		ClientOptions: []client.Option{
			client.WithSSHAuth(g.keys),
		},
	})

	return nil
}

func NewGitProvisioner(repo *GitRepository, files []GitFile) (*Git, error) {
	if repo.URL == "" {
		return nil, errors.New("repositoryURL must not be empty")
	}

	if strings.HasPrefix(repo.URL, "https://") {
		return nil, errors.New("HTTPS authentication is not supported")
	}

	if repo.KeyAuth == nil || repo.KeyAuth.PrivateKey == nil {
		return nil, errors.New("no private key specified")
	}

	if repo.CommitAuthor == nil {
		return nil, errors.New("commit author must be specified")
	}

	if repo.CommitAuthor.Email == "" || repo.CommitAuthor.Name == "" {
		return nil, errors.New("commit author's name and email must be specified")
	}

	keys, err := ssh.NewPublicKeys("git", repo.KeyAuth.PrivateKey, "")
	if err != nil {
		return nil, err
	}

	return &Git{
		repository: repo,
		files:      files,
		keys:       keys,
	}, nil
}

type GitRepository struct {
	URL          string        `json:"url"`
	Ref          string        `json:"ref,omitempty"`
	KeyAuth      *GitKeyAuth   `json:"keyAuth,omitempty"`
	CommitAuthor *CommitAuthor `json:"commitAuthor"`
}

type GitKeyAuth struct {
	PrivateKey []byte
}

type CommitAuthor struct {
	// Commit author's name.
	Name string `json:"name"`

	// Commit author's email.
	Email string `json:"email"`
}

type GitFile struct {
	// Full path to the file to provision.
	Path string `json:"path"`

	// Content of the file to provision.
	Content []byte `json:"content"`
}

func generateRepoName() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 12)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
