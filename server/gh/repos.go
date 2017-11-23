package gh

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// GetRepository gets a repository.
func (c *Client) GetRepository() (*github.Repository, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	repo, _, err := c.Repositories.Get(context.Background(), c.owner, c.repo)
	if err != nil {
		logrus.Errorf("failed to get repository c.repo %s: %v", c.repo, err)
		return nil, err
	}
	logrus.Debugf("succeed in getting repository %s", c.repo)
	return repo, nil
}
