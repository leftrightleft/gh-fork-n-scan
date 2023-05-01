package main

import (
	"flag"
	"fmt"
	"os"

	"encoding/json"
	"log"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/v2/pkg/tableprinter"
	"github.com/cli/go-gh/v2/pkg/term"
)

// struct  supporting the respons from getRepos
type Repo struct {
	Language struct {
		Name string `json:"name"`
	} `json:"primaryLanguage"`
	NameWithOwner string `json:"nameWithOwner"`
	Name          string `json:"name"`
	Owner         struct {
		Login string `json:"login"`
	} `json:"owner"`
}

type Repos []Repo

func readFlags() (string, string, string) {
	var sourceOrg, sourceRepo, destOrg string
	var showHelp bool
	flag.StringVar(&sourceOrg, "s", "", "source organization name")
	flag.StringVar(&sourceRepo, "r", "", "source repository name")
	flag.StringVar(&destOrg, "d", "", "destination organization name")
	flag.BoolVar(&showHelp, "h", false, "show help")
	flag.Parse()
	if showHelp {
		flag.Usage()
		os.Exit(0)
	}
	if destOrg == "" {
		fmt.Println("Please provide both source and destination org names")
		flag.PrintDefaults()
		os.Exit(1)
	}
	return sourceOrg, sourceRepo, destOrg
}

func GetRepos(sourceOrg string) Repos {
	var repos []Repo
	languages := []string{"python", "go", "javascript", "ruby"}
	for _, language := range languages {
		var languageRepos []Repo
		args := []string{"repo", "list", sourceOrg, "-l", language, "--no-archived", "--source", "--visibility", "public", "--json", "nameWithOwner,owner,name,primaryLanguage"}
		stdOut, stdErr, err := gh.Exec(args...)
		if err != nil {
			log.Fatal(err)
		}

		if stdErr.String() != "" {
			fmt.Println(stdErr.String())
		}

		err = json.Unmarshal(stdOut.Bytes(), &languageRepos)
		if err != nil {
			log.Fatal(err)
		}
		repos = append(repos, languageRepos...)
	}
	return repos
}

func ForkRepos(repos []Repo, destOrg string) {
	for _, repo := range repos {
		args := []string{"repo", "fork", repo.NameWithOwner, "--org", destOrg}
		_, stdErr, err := gh.Exec(args...)
		if stdErr.String() != "" {
			fmt.Println(stdErr.String())
		}

		if err != nil {
			log.Fatal(err)
		}

		log.Println("Forked", repo.NameWithOwner)

	}
}

func AssignTopicToRepos(repos []Repo, destOrg string) string {
	var topic string
	for _, repo := range repos {
		topic = repo.Owner.Login
		destRepo := fmt.Sprintf("%s/%s", destOrg, repo.Name)
		args := []string{"repo", "edit", destRepo, "--add-topic", topic}
		_, stdErr, err := gh.Exec(args...)

		if stdErr.String() != "" {
			fmt.Println(stdErr.String())
		}

		if err != nil {
			log.Fatal(err)
		}

		log.Println("Assigned topic", topic, "to", destRepo)
	}

	return topic
}

// TODO: Enable default setup for code scanning
// func EnableDefaultSetup(destOrg string, repos []Repo) {
// 	for _, repo := range repos {
// 		args := []string{"api", "--method", "patch", fmt.Sprintf("/repos/%s/%s/code-scanning/default-setup", destOrg, repo.Name), "-f", "state='configured'"}
// 		_, stdErr, err := gh.Exec(args...)
// 		if stdErr.String() != "" {
// 			log.Println(stdErr.String())
// 		}

// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		// fmt.Println(response)
// 	}
// }

// Print tabular data to a terminal or in machine-readable format for scripts.
func PrintTable(repos []Repo) int {
	terminal := term.FromEnv()
	termWidth, _, _ := terminal.Size()
	t := tableprinter.New(terminal.Out(), terminal.IsTerminalOutput(), termWidth)

	green := func(s string) string {
		return "\x1b[32m" + s + "\x1b[m"
	}

	// add a field that will render with color and will not be auto-truncated
	count := 0
	for _, repo := range repos {
		t.AddField("-", tableprinter.WithColor(green), tableprinter.WithTruncate(nil))
		t.AddField(repo.NameWithOwner, tableprinter.WithTruncate(nil))
		t.AddField(repo.Language.Name, tableprinter.WithTruncate(nil))
		t.EndRow()
		count++

	}
	if err := t.Render(); err != nil {
		log.Fatal(err)
	}

	return count
}

func main() {
	sourceOrg, sourceRepo, destOrg := readFlags()
	if sourceRepo != "" {
		fmt.Printf("Forking %s", sourceRepo)
		ForkRepos(Repos{Repo{NameWithOwner: sourceRepo}}, destOrg)
		return
	}
	repos := GetRepos(sourceOrg)
	fmt.Println("Forking the following repos")
	repoCount := PrintTable(repos)
	fmt.Printf("\nDo you want to continue (%d repos)? (y/n)", repoCount)
	var input string
	fmt.Scanln(&input)
	if input != "y" {
		fmt.Println("Exiting...")
		return
	}

	ForkRepos(repos, destOrg)
	if sourceOrg != "" {
		topic := AssignTopicToRepos(repos, destOrg)
		// var input string
		// fmt.Printf("\nEnable code scanning on %d repos? (y/n)", repoCount)
		// fmt.Scanln(&input)
		// if input != "y" {
		// 	fmt.Printf("\nDone! 🤙\n\n")
		// 	fmt.Printf("Visit the following URL to view code scanning results:\n")
		// 	fmt.Printf("https://github.com/orgs/%s/security/alerts/code-scanning?query=is%%3Aopen+topic%%3A%s\n", destOrg, topic)
		// 	EnableDefaultSetup(destOrg, repos)
		// 	return
		// }

		fmt.Printf("\nDone! 🤙\n\n")
		fmt.Printf("Visit the following URL to enable code scanning:\n")
		fmt.Printf("https://github.com/orgs/%s/security/coverage?query=topic%%3A%s\n", destOrg, topic)

	}

}
