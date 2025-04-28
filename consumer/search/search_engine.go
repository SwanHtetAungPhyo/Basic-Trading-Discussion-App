package main

import (
	"bufio"
	"fmt"
	"github.com/agnivade/levenshtein"
	"os"
	"path/filepath"
	"strings"
)

type Node struct {
	Word  string
	Left  *Node
	Right *Node
}

func NewNode(word string) *Node {
	return &Node{
		Word: word,
	}
}

func (n *Node) Insert(word string) *Node {
	if n == nil {
		return &Node{Word: word}
	}
	if n.Word < word {
		n.Right = n.Right.Insert(word)
	} else if n.Word > word {
		n.Left = n.Left.Insert(word)
	}
	return n
}

func (n *Node) Search(word string) (bool, *string) {
	if n == nil {
		return false, nil
	}
	if n.Word == word {
		return true, &n.Word
	}
	distance := n.LevenshteinDistance(word)
	if distance <= 3 {
		return true, &n.Word
	}
	if n.Word < word {
		return n.Right.Search(word)
	}
	return n.Left.Search(word)
}

func (n *Node) LevenshteinDistance(word string) int {
	return levenshtein.ComputeDistance(n.Word, word)
}

func (n *Node) PrintTree(indent string) {
	if n == nil {
		return
	}
	fmt.Println(indent + n.Word)
	n.Left.PrintTree(indent + "  ")
	n.Right.PrintTree(indent + "  ")
}

func LoadFile(path string) []string {
	pwd, _ := os.Getwd()
	f, err := os.Open(filepath.Join(pwd + "/" + path))
	if err != nil {
		panic(err.Error())
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err.Error())
		}
	}(f)
	reader := bufio.NewReader(f)
	words := make([]string, 0)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		lineWords := strings.Fields(strings.Trim(string(line), ".,/{}:\"\"[]()!@#$%^&*~`"))
		words = append(words, lineWords...)
	}
	return words
}

func main() {
	var word string
	fmt.Println("Enter the word that you want to search:")
	_, err := fmt.Scanln(&word)
	if err != nil {
		return
	}

	lines := LoadFile("main.txt")
	if len(lines) == 0 {
		fmt.Println("No words found in file.")
		return
	}

	mainNode := NewNode(lines[0])
	for i := 1; i < len(lines); i++ {
		mainNode.Insert(lines[i])
	}

	// Print the tree structure for debugging
	fmt.Println("Tree structure:")
	mainNode.PrintTree("")

	// Search for the word
	found, result := mainNode.Search(word)
	if found {
		fmt.Println("Found similar word:", *result)
	} else {
		fmt.Println("No similar word found.")
	}
}
