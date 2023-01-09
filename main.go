package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	files, err := getFilesRecursively("./archive/tweets-md")
	if err != nil {
		log.Fatal(err)
	}

	contents, err := getFileContents(files)
	if err != nil {
		log.Fatal(err)
	}

	tweets, err := getTweetsByFile(contents)
	if err != nil {
		log.Fatal(err)
	}

	for file, fileTweets := range tweets {
		fmt.Println(file)
		for _, tweet := range fileTweets {
			fmt.Println(file)
			fmt.Println(tweet)
		}
	}
}

func getFilesRecursively(dir string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return files, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			additionalFiles, err := getFilesRecursively(
				fmt.Sprintf("%s/%s", dir, entry.Name()))
			if err != nil {
				return files, err
			}

			files = append(files, additionalFiles...)
		} else {
			files = append(files, fmt.Sprintf("%s/%s", dir, entry.Name()))
		}
	}

	return files, nil
}

func getFileContents(files []string) (map[string]string, error) {
	var contents = make(map[string]string)

	for _, file := range files {
		// I don't want any tweets before 2020
		if !strings.Contains(file, "202") {
			continue
		}

		content, err := ioutil.ReadFile(file)
		if err != nil {
			return contents, err
		}

		contents[file] = string(content)
	}

	return contents, nil
}

func getTweetsByFile(archiveContents map[string]string) (map[string][]string, error) {
	var tweets = make(map[string][]string)

	for key, content := range archiveContents {
		fileTweets := strings.Split(content, "----")
		tweets[key] = append(tweets[key], fileTweets...)
	}

	return tweets, nil
}
