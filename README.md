# git-wip

Quickly create and remove "wip" commits.

`git-wip` stages all tracked changes and creates a commit with the message "wip". When you're ready to continue working, `git-unwip` removes consecutive wip commits from HEAD and leaves your changes unstaged.

## Install

```
brew install danix9000/tap/git-wip
```

## Usage

```
# Create a wip commit with all tracked changes
git wip

# Preview without committing
git wip --dry-run

# Remove wip commits and restore changes
git unwip
```

`git-unwip` is also available via `git wip --unwip`.
