package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	_ "embed"

	"github.com/pilcrowonpaper/go-json"
	"golang.org/x/net/html"
)

//go:embed topics.json
var topicIdsJSON string

func main() {
	topicIdsJSONArray, err := json.ParseArray(topicIdsJSON)
	if err != nil {
		fmt.Printf("Failed to json parse topics.json: %s\n", err.Error())
		os.Exit(0)
	}
	topicIds := []string{}
	for i := range topicIdsJSONArray.Length {
		topicId, err := topicIdsJSONArray.GetString(i)
		if err != nil {
			fmt.Printf("Failed to get string from json array: %s\n", err.Error())
			os.Exit(0)
		}
		topicIds = append(topicIds, topicId)
	}

	err = os.RemoveAll("dist")
	if err != nil {
		fmt.Printf("Failed to remove 'dist' directory and its content: %s\n", err.Error())
		os.Exit(0)
	}

	err = os.MkdirAll("dist", os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to create directory 'dist': %s\n", err.Error())
		os.Exit(0)
	}

	topicsDirectoryEntries, err := os.ReadDir("topics")
	if err != nil {
		fmt.Printf("Failed to read 'topics' directory: %s\n", err.Error())
		os.Exit(0)
	}

	topics := []topicStruct{}
	for _, topicsDirectoryEntry := range topicsDirectoryEntries {
		if topicsDirectoryEntry.IsDir() {
			continue
		}
		filename := topicsDirectoryEntry.Name()
		extension := path.Ext(filename)
		if extension != ".html" {
			continue
		}

		topicId := strings.TrimSuffix(filename, ".html")

		topicHTMLFilePath := path.Join("topics", filename)
		topicHTMLBytes, err := os.ReadFile(topicHTMLFilePath)
		if err != nil {
			fmt.Printf("Failed to read %s file: %s\n", topicHTMLFilePath, err.Error())
			os.Exit(0)
		}

		topicHTMLNode, err := html.Parse(bytes.NewReader(topicHTMLBytes))
		if err != nil {
			fmt.Printf("Failed to html parse %s file: %s\n", topicHTMLFilePath, err.Error())
			os.Exit(0)
		}

		topicTitle, ok := getH1TitleFromHTMLNode(topicHTMLNode)
		if !ok {
			fmt.Printf("Title not defined in %s file\n", topicHTMLFilePath)
			os.Exit(0)
		}

		topicPathId := strings.ReplaceAll(topicId, "_", "-")

		topicPageHTML := createTopicPageHTML(topicTitle, topicPathId, string(topicHTMLBytes))
		destinationFilename := fmt.Sprintf("%s.html", topicPathId)
		destinationFilePath := path.Join("dist", destinationFilename)
		err = os.WriteFile(destinationFilePath, []byte(topicPageHTML), os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to write file '%s': %s\n", destinationFilePath, err.Error())
			os.Exit(0)
		}

		topic := topicStruct{
			id:     topicId,
			title:  topicTitle,
			pathId: topicPathId,
		}

		topics = append(topics, topic)
	}

	sortedTopics := []topicStruct{}
	for _, topicId := range topicIds {
		found := false
		for _, topic := range topics {
			if topic.id == topicId {
				sortedTopics = append(sortedTopics, topic)
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("invalid topic id '%s'\n", topicId)
			os.Exit(0)
		}
	}

	homePageHTML := createHomePageHTML(sortedTopics)
	err = os.WriteFile("dist/index.html", []byte(homePageHTML), os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to write file 'dist/index.html': %s\n", err.Error())
		os.Exit(0)
	}
}

func getH1TitleFromHTMLNode(htmlNode *html.Node) (string, bool) {
	if htmlNode.Type == html.ElementNode && htmlNode.Data == "h1" && htmlNode.FirstChild != nil && htmlNode.FirstChild.Type == html.TextNode {
		return htmlNode.FirstChild.Data, true
	}
	childsNodes := htmlNode.ChildNodes()
	for childsNode := range childsNodes {
		title, ok := getH1TitleFromHTMLNode(childsNode)
		if ok {
			return title, true
		}
	}
	return "", false
}
