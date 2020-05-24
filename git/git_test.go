package git

import (
	"log"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const testingDirectory = ".tests"

var testBranches = []string{"development", "staging", "feat/abc", "master"}

var repo *gogit.Repository

func TestMain(m *testing.M) {
	check(setup())
	os.Exit(m.Run())
}

func TestBranchNames(t *testing.T) {
	branches, err := branchNamesFromRepo(repo)

	assertNoError(t, err)

	sort.Strings(testBranches)
	sort.Strings(branches)

	if !reflect.DeepEqual(branches, testBranches) {
		t.Errorf("Expected %v but got %v", testBranches, branches)
	}
}

func TestCheckout(t *testing.T) {
	tests := []struct {
		name        string
		branch      string
		expectError bool
	}{
		{"existing branch", "feat/abc", false},
		{"non-existing branch", "nani", true},
	}

	for _, tt := range tests {
		err := checkoutInRepo(repo, tt.branch)

		if tt.expectError {
			if err == nil {
				t.Errorf("Expected an error but got none")
			}
			continue
		}

		assertNoError(t, err)

		ref, err := repo.Head()
		assertNoError(t, err)

		if ref.Name().Short() != tt.branch {
			t.Errorf("Expected HEAD to be at %v but is at %v", tt.branch, ref.Name().Short())
		}
	}

}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected no error but got one: %v", err)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func setup() error {
	os.RemoveAll(testingDirectory)

	err := os.Mkdir(testingDirectory, 0700)
	if err != nil {
		return err
	}

	os.Chdir(testingDirectory)

	err = os.RemoveAll(".git")
	if err != nil {
		return err
	}

	r, err := gogit.PlainInit(".", false)
	if err != nil {
		return err
	}

	file, err := os.Create("test.txt")
	if err != nil {
		return err
	}

	file.WriteString("Hello World!")
	err = file.Close()
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	w.Add(".")

	_, err = w.Commit("initial commit", &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  "Levi Ackermann",
			Email: "lackermann@email.com",
			When:  time.Now(),
		},
	})

	if err != nil {
		return err
	}

	for _, b := range testBranches {
		headRef, err := r.Head()
		if err != nil {
			return err
		}

		ref := plumbing.NewReferenceFromStrings("refs/heads/"+b, headRef.Hash().String())
		err = r.Storer.SetReference(ref)
		if err != nil {
			return err
		}
	}
	repo = r
	return nil
}
