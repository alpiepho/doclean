package main

// go run doclean.go somedir1 somedir2
//
// This is a Golang utility to compare two directory trees and remove
// any files in the second tree that are not in the first.
// This is useful for cleaning a software source tree after a build
// process runs and possibly leaves build artifacts.

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var treePath1 string
var treePath2 string
var list1 []string
var list2 []string
var list3 []string // keep list

func usage() {
	fmt.Println("usage: go run doclean.go <clean dir> <dirty dir> [<keep file>]")
	fmt.Println("	<clean dir> - known clean directory/tree")
	fmt.Println("	<dirty dir> - dirty directory/tree to be cleaned vs frist tree")
	fmt.Println("	[<keep file>] - optional file containing list of dirty files to keep")

}

func getKeepList(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file %s\n", filename)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		list3 = append(list3, line)
		temp := filepath.Dir(line)
		list3 = append(list3, temp)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s\n", filename)
		return
	}
}

// save all files/directories found
func visit1(path string, d fs.DirEntry, err error) error {
	list1 = append(list1, path)
	return nil
}

// save only files/directories not found in tree1
func visit2(path string, d fs.DirEntry, err error) error {
	temp2 := strings.Replace(path, treePath2, "", 1)
	// DEBUG
	// fmt.Println(treePath2)
	// fmt.Println(temp2)
	// fmt.Println("")
	found := false
	for _, n := range list1 {
		temp1 := strings.Replace(n, treePath1, "", 1)
		if temp1 == temp2 {
			found = true
			break
		}
	}

	// test if path matches keep list, mark as found
	for _, k := range list3 {
		if strings.Contains(path, k) {
			found = true
			break
		}
	}

	if !found {
		list2 = append(list2, path)
	}
	return nil
}

func main() {
	flag.Parse()
	flag.CommandLine.Args()
	if len(flag.CommandLine.Args()) < 1 {
		usage()
		return
	}
	treePath1 = flag.Arg(0)
	treePath2 = flag.Arg(1)
	// for windows
	treePath1 = strings.Replace(treePath1, "/", "\\", -1)
	treePath2 = strings.Replace(treePath2, "/", "\\", -1)

	if len(flag.CommandLine.Args()) == 3 {
		getKeepList(flag.Arg(2))
	}

	// walk first tree to save list of files/directories
	filepath.WalkDir(treePath1, visit1)
	// walk second tree find new files/directories
	filepath.WalkDir(treePath2, visit2)

	// reverse to remove files before directories
	sort.Sort(sort.Reverse(sort.StringSlice(list2)))
	for _, n := range list2 {
		fmt.Println("removing " + n)
		// DEBUG
		os.Remove(n)
	}
	fmt.Println(len(list2))
	fmt.Println("Please verify with: 'meld %s %s'", flag.Arg(0), flag.Arg(0))
}
