package differs

import (
	"fmt"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/net/context"
)

// History compares the Docker history for each image.
func History(img1, img2 string) string {
	return get_history_diff(img1, img2)
}

func get_history_list(image string) []string {
	ctx := context.Background()
	cli, _ := client.NewEnvClient()
	history, _ := cli.ImageHistory(ctx, image)

	strhistory := make([]string, len(history))
	for i, layer := range history {
		layer_description := strings.TrimSpace(layer.CreatedBy)
		strhistory[i] = fmt.Sprintf("%s\n", layer_description)
	}
	return strhistory
}

func get_history_diff(image1 string, image2 string) string {
	history1 := get_history_list(image1)
	history2 := get_history_list(image2)

	diff := difflib.ContextDiff{
		A:        history1,
		B:        history2,
		FromFile: "IMAGE " + image1,
		ToFile:   "IMAGE " + image2,
		Eol:      "\n",
	}
	result, _ := difflib.GetContextDiffString(diff)
	return result
}
