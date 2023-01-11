package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// TweetType classifies a type of tweet
type TweetType int

const (
	// Post An original tweet
	Post TweetType = iota + 1 // 1

	// Reply A reply to a tweet
	Reply // 2

	// Retweet a repost of a tweeet
	Retweet // 3
)

// Tweet is a post on twitter
type Tweet struct {
	Type TweetType
	Body string
}

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

	totalTweets := 0
	potentialQuoteTweets := 0
	for file, fileTweets := range tweets {
		fmt.Println(file)
		for _, tweet := range fileTweets {
			trimmedTweet := strings.TrimSpace(tweet)
			fmt.Println(file)
			fmt.Println(trimmedTweet)
			totalTweets = totalTweets + 1

			if strings.Contains(trimmedTweet, "twitter\\.com") {
				potentialQuoteTweets = potentialQuoteTweets + 1
			}
		}
	}

	postTweets, err := getPostTweets(tweets)
	if err != nil {
		log.Fatal(err)
	}

	replyTweets, err := getReplyTweets(tweets)
	if err != nil {
		log.Fatal(err)
	}

	retweets, err := getRetweets(tweets)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total Tweets: %d\n", totalTweets)
	fmt.Printf("Post Tweets: %d\n", len(postTweets))
	fmt.Printf("Reply Tweets: %d\n", len(replyTweets))
	fmt.Printf("Retweets: %d\n", len(retweets))
	fmt.Printf("Potential Quote Tweets: %d\n", potentialQuoteTweets)
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

func getPostTweets(tweets map[string][]string) ([]Tweet, error) {
	var postTweets []Tweet

	for _, fileTweets := range tweets {
		for _, tweet := range fileTweets {
			trimmedTweet := strings.TrimSpace(tweet)
			if strings.HasPrefix(trimmedTweet, ">") && !strings.HasPrefix(trimmedTweet, "> RT") {
				postTweets = append(postTweets, Tweet{
					Type: Post,
					Body: tweet,
				})
			}
		}
	}

	return postTweets, nil
}

func getReplyTweets(tweets map[string][]string) ([]Tweet, error) {
	var replyTweets []Tweet

	for _, fileTweets := range tweets {
		for _, tweet := range fileTweets {
			trimmedTweet := strings.TrimSpace(tweet)
			if strings.HasPrefix(trimmedTweet, "Replying to") {
				replyTweets = append(replyTweets, Tweet{
					Type: Reply,
					Body: trimmedTweet,
				})
			}
		}
	}

	return replyTweets, nil
}

func getRetweets(tweets map[string][]string) ([]Tweet, error) {
	var retweets []Tweet

	for _, fileTweets := range tweets {
		for _, tweet := range fileTweets {
			trimmedTweet := strings.TrimSpace(tweet)
			if strings.HasPrefix(trimmedTweet, "> RT") {
				retweets = append(retweets, Tweet{
					Type: Retweet,
					Body: tweet,
				})
			}
		}
	}

	return retweets, nil
}
