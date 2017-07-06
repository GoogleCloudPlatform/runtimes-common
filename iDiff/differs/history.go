package differs

import (
	"fmt"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/net/context"
)

// History compares the Docker history for each image.
func History(img1, img2 string, json bool) (string, error) {
	return getHistoryDiff(img1, img2, json)
}

func getHistoryList(image string) ([]string, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return []string{}, err
	}
	history, err := cli.ImageHistory(ctx, image)
	if err != nil {
		return []string{}, err
	}

	strhistory := make([]string, len(history))
	for i, layer := range history {
		layer_description := strings.TrimSpace(layer.CreatedBy)
		strhistory[i] = fmt.Sprintf("%s\n", layer_description)
	}
	return strhistory, nil
}

func getHistoryDiff(image1 string, image2 string, json bool) (string, error) {
	history1, err := getHistoryList(image1)
	if err != nil {
		return "", err
	}
	history2, err := getHistoryList(image2)
	if err != nil {
		return "", err
	}
	diff := difflib.ContextDiff{
		A:        history1,
		B:        history2,
		FromFile: "IMAGE " + image1,
		ToFile:   "IMAGE " + image2,
		Eol:      "\n",
	}
	result, _ := difflib.GetContextDiffString(diff)
	if json {
		return jsonDiff(result), nil
	}
	return result, nil
}

func jsonDiff(diff string) string {
	// TODO write JSON parser
	return "Hello World"
}
