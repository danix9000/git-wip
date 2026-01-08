package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func lastNonWipCommit(repository *git.Repository) (*object.Commit, error) {
	commitIter, err := repository.Log(&git.LogOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer commitIter.Close()

	wipCommits := 0

	for {
		c, err := commitIter.Next()
		if err != nil {
			return nil, err
		}

		if c.Message == "wip" {
			wipCommits += 1
		} else {
			return c, nil
		}
	}
}

func main() {
	execName := filepath.Base(os.Args[0])
	unwipCommand := execName == "git-unwip"

	dryRunFlag := flag.Bool("dry-run", false, "a dry run without generating a new commit")
	unwipFlag := flag.Bool("unwip", unwipCommand, "Unpacks wip commits")
	flag.Parse()

	repository, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}

	worktree, err := repository.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	status, err := worktree.Status()
	if err != nil {
		log.Fatal(err)
	}

	if *unwipFlag {
		lastNonWipCommit, err := lastNonWipCommit(repository)
		if err != nil {
			log.Fatal(err)
		}

		headCommitRef, err := repository.Head()
		if err != nil {
			log.Fatal(err)
		}

		if headCommitRef.Hash() == lastNonWipCommit.Hash {
			fmt.Println("No wip commits")
		} else {
			fmt.Println("Reverted HEAD to last non wip commit:", lastNonWipCommit.Hash)
			err = worktree.Reset(&git.ResetOptions{
				Commit: plumbing.NewHash(lastNonWipCommit.Hash.String()),
			})
			if err != nil {
				log.Fatal(err)
			}

		}

		return
	}

	hasTrackedChanges := false
	for _, fileStatus := range status {
		if fileStatus.Staging != git.Untracked || fileStatus.Worktree != git.Untracked {
			hasTrackedChanges = true
			break
		}
	}

	if status.IsClean() {
		fmt.Println("Nothing to commit, working tree clean")
	} else if !status.IsClean() && !hasTrackedChanges {
		fmt.Println("Nothing to commit, all changes are to untracked files")
	} else {
		if *dryRunFlag {
			fmt.Println("Dry run, not committing changes")
		} else {
			commitHash, err := worktree.Commit("wip", &git.CommitOptions{All: true})
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Added a new wip commit", commitHash)
		}
	}
}
