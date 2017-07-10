package differs

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// History compares the Docker history for each image.
func History(img1, img2 string, json bool) (string, error) {
	return getHistoryDiff(img1, img2, json)
}

func getHistoryList(img string) ([]string, error) {
	validDocker, err := validDockerVersion()
	if err != nil {
		return []string{}, err
	}
	var history []image.HistoryResponseItem
	if validDocker {
		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			return []string{}, err
		}
		history, err = cli.ImageHistory(ctx, img)
		if err != nil {
			return []string{}, err
		}
	} else {
		// TODO call local docker with exec
		return []string{}, nil
	}

	strhistory := make([]string, len(history))
	for i, layer := range history {
		layerDescription := strings.TrimSpace(layer.CreatedBy)
		strhistory[i] = fmt.Sprintf("%s\n", layerDescription)
	}
	return strhistory, nil
}

type HistDiff struct {
	Image1 string
	Image2 string
	Adds   []string
	Dels   []string
}

func getHistoryDiff(image1, image2 string, json bool) (string, error) {
	history1, err := getHistoryList(image1)
	if err != nil {
		return "", err
	}
	history2, err := getHistoryList(image2)
	if err != nil {
		return "", err
	}
	adds := utils.GetAdditions(history1, history2)
	dels := utils.GetDeletions(history1, history2)
	diff := HistDiff{image1, image2, adds, dels}
	if json {
		return utils.JSONify(diff)
	}
	result := formatDiff(diff)
	return result, nil
}

func formatDiff(diff HistDiff) string {
	const histTemp = `Docker file lines found only in {{.Image1}}:{{block "list" .Adds}}{{"\n"}}{{range .}}{{print "-" .}}{{end}}{{end}}
Docker file lines found only in {{.Image2}}:{{block "list2" .Dels}}{{"\n"}}{{range .}}{{print "-" .}}{{end}}{{end}}`

	funcs := template.FuncMap{"join": strings.Join}

	histTemplate, err := template.New("histTemp").Funcs(funcs).Parse(histTemp)
	if err != nil {
		log.Fatal(err)
	}
	if err := histTemplate.Execute(os.Stdout, diff); err != nil {
		log.Fatal(err)
	}
	return ""
}
