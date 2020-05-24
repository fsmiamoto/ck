package git

import (
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Repository struct {
	Path string
	r    *gogit.Repository
}

func OpenRepository(path string) (*Repository, error) {
	repo, err := gogit.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Path: path,
		r:    repo,
	}, nil
}

func (repo *Repository) Checkout(branch string) error {
	return checkoutInRepo(repo.r, branch)
}

func (repo *Repository) Branches() ([]string, error) {
	return branchNamesFromRepo(repo.r)
}

func branchNamesFromRepo(repo *gogit.Repository) ([]string, error) {
	refs, err := repo.Branches()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0)

	refs.ForEach(func(ref *plumbing.Reference) error {
		names = append(names, ref.Name().Short())
		return nil
	})

	return names, nil
}

func checkoutInRepo(repo *gogit.Repository, branch string) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	return worktree.Checkout(&gogit.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + branch),
	})
}
