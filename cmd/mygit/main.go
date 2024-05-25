package main

import (
	"fmt"
	"github.com/codecrafters-io/git-starter-go/internal/git_commands"
	"github.com/codecrafters-io/git-starter-go/internal/utils"
	"log"
	"os"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: mygit <command> [<args>...]\n")
	}

	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Fatalf("Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			log.Fatalf("Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")
	case "cat-file":
		if len(os.Args) < 4 {
			log.Fatalf("usage: cat-file -p <blob_sha>")
		}

		err := git_commands.CatFile(os.Args[3], os.Args[2], utils.AppUtils{})

		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unknown command %s\n", command)
	}
}
