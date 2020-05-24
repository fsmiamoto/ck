package main

import (
	"log"
	"os"

	"github.com/fsmiamoto/ck/git"
	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	path := "."

	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	if err := run(path); err != nil {
		log.Fatal(err)
	}
}

func run(path string) error {
	repo, err := git.OpenRepository(path)

	branches, err := repo.Branches()
	if err != nil {
		log.Fatal(err)
	}

	index, err := fuzzyfinder.Find(branches, func(i int) string {
		return branches[i]
	}, fuzzyfinder.WithPromptString("Checkout: "))

	if err != nil {
		log.Fatal(err)
	}

	return repo.Checkout(branches[index])
}
