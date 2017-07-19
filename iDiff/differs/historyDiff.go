package differs

import (
	"html/template"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
	"github.com/golang/glog"
)

type HistoryDiffer struct {
}

func (d HistoryDiffer) Diff(image1, image2 utils.Image, json, eng bool) (string, error) {
	return getHistoryDiff(image1, image2, json, eng)
}

/*func getHistoryList(img string, eng bool) ([]string, error) {
	validDocker, err := utils.ValidDockerVersion(eng)
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
		glog.Info("Docker version incompatible with api, shelling out to local Docker client.")
		history, err = utils.GetImageHistory(img)
		if err != nil {
			return []string{}, err
		}
	}

	strhistory := make([]string, len(history))
	for i, layer := range history {
		layerDescription := strings.TrimSpace(layer.CreatedBy)
		strhistory[i] = fmt.Sprintf("%s\n", layerDescription)
	}
	return strhistory, nil
}*/

type HistDiff struct {
	Image1 string
	Image2 string
	Adds   []string
	Dels   []string
}

func getHistoryDiff(image1, image2 utils.Image, json bool, eng bool) (string, error) {
	history1 := image1.History
	history2 := image2.History

	adds := utils.GetAdditions(history1, history2)
	dels := utils.GetDeletions(history1, history2)
	diff := HistDiff{image1.FSPath, image2.FSPath, adds, dels} //TODO: Add name to Image struct
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
		glog.Error(err)
	}
	if err := histTemplate.Execute(os.Stdout, diff); err != nil {
		glog.Error(err)
	}
	return ""
}
