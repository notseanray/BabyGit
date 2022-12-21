package main

import (
	"fmt"
	"os"

	// "os/exec"
	// "sort"
	// "strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func pull(path string) {
	// We instantiate a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(path)
	CheckIfError(err)

	// Get the working directory for the repository
	w, err := r.Worktree()
	CheckIfError(err)

	// Pull the latest changes from the origin remote and merge into the current branch
	Info("git pull origin")
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	CheckIfError(err)

	// Print the latest commit that was just pulled
	ref, err := r.Head()
	CheckIfError(err)
	commit, err := r.CommitObject(ref.Hash())
	CheckIfError(err)

	fmt.Println(commit)
}

type Commit struct {
    author string
    authorEmail string
    hash string
    date string
    message string
}

func get_commits(path string) []Commit {
	r, err := git.PlainOpen(path)
	CheckIfError(err)

	// Length of the HEAD history
	Info("git rev-list HEAD --count")

	// ... retrieving the HEAD reference
	ref, err := r.Head()
	CheckIfError(err)

	// ... retrieves the commit history
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	CheckIfError(err)

	// ... just iterates over the commits
    var commits []Commit
	err = cIter.ForEach(func(c *object.Commit) error {
        commits = append(commits, Commit {
            author: c.Author.Name,
            authorEmail: c.Author.Email,
            hash: c.Hash.String(),
            date: "s",
            message: strings.Trim(c.Message, ""),
        })
		return nil
	})
	CheckIfError(err)
    return commits
}

var RepoPath = ""

func MenuSelect(uri fyne.ListableURI, e error) {
    RepoPath = uri.Path()
}

func main() {
    commits := get_commits("/home/sean/Desktop/stuff/crisp")
	a := app.New()
	w := a.NewWindow("BabyGit")
	var buttons []fyne.CanvasObject
	var rebaseStatus = "h"
	var layout []fyne.CanvasObject
	var rebaseMenu []fyne.CanvasObject
	pullButton := widget.NewButton("Pull", func() {
		fmt.Println("test")
		rebaseStatus = "yup"
	})
	layout = append(layout, pullButton)
	rebaseLabel := widget.NewLabel(rebaseStatus)
	rebaseButtonLocal := widget.NewButton("Rebase, set Remote -> Local", func() {
		rebaseLabel.SetText(rebaseStatus)
	})
	rebaseButtonRemote := widget.NewButton("Rebase, set Local -> Remote", func() {
		rebaseLabel.SetText(rebaseStatus)
	})
	rebaseMenu = append(rebaseMenu, rebaseButtonRemote)
	rebaseMenu = append(rebaseMenu, rebaseButtonLocal)

	screenEntry := widget.NewButton("Push", func() {
		rebaseButtonLocal.Disable()
		rebaseButtonRemote.Disable()
		fmt.Println("test")
	})
	layout = append(layout, screenEntry)
	subContent := container.NewGridWithColumns(len(layout), layout...)
	screenCard := widget.NewCard(
		"Repository",
		"Last Sync: ", subContent)
	rebasePanel := container.NewGridWithColumns(len(rebaseMenu), rebaseMenu...)
    fsMenu := dialog.NewFolderOpen(MenuSelect, w)
	fsSelect := widget.NewButton("Choose Respository", func() {
        fsMenu.Show()
	})
    buttons = append(buttons, fsSelect)
	rebaseCard := widget.NewCard(
		"Rebase Settings",
		"Status: ", rebasePanel)
	buttons = append(buttons, screenCard)
	buttons = append(buttons, rebaseCard)
    var commitCards []fyne.CanvasObject
    for _, item := range(commits) {
        // widget.NewCard
        itemCard := widget.NewCard(
            "Author: " + item.author + " <" + item.authorEmail + ">" ,
            "Message: " + item.message, widget.NewLabel(item.hash))
        commitCards = append(commitCards, itemCard)
    }
    commitScroll := container.NewGridWithRows(len(commitCards), commitCards...)
    buttons = append(buttons, container.NewVScroll(commitScroll))
	content := container.NewGridWithRows(4, buttons...)
	w.SetContent(content)
	w.ShowAndRun()
}
