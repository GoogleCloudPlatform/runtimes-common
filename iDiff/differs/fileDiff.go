package differs

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
)

type FileDiffer struct {
}

// FileDiff diffs two packages and compares their contents
func (d FileDiffer) Diff(image1, image2 utils.Image) (DiffResult, error) {
	img1 := image1.FSPath
	img2 := image2.FSPath

	diff, err := diffImageFiles(img1, img2)
	return &utils.DirDiffResult{Diff: diff}, err
}

func diffImageFiles(img1, img2 string) (utils.DirDiff, error) {
	var diff utils.DirDiff

	img1Contents, err := getImageContents(img1)
	if err != nil {
		return diff, fmt.Errorf("Error parsing image %s contents: %s", img1, err)
	}
	img2Contents, err := getImageContents(img2)
	if err != nil {
		return diff, fmt.Errorf("Error parsing image %s contents: %s", img2, err)
	}

	for layer1, contents1 := range img1Contents {
		sameLayer := false
		for layer2, contents2 := range img2Contents {
			if checkSameLayer(contents1, contents2) {
				delete(img2Contents, layer2)
				sameLayer = true
				break
			}
		}
		if sameLayer {
			delete(img1Contents, layer1)
		}
	}

	img1Dir := utils.Directory{
		Root:    img1,
		Content: getContentList(img1Contents),
	}
	img2Dir := utils.Directory{
		Root:    img2,
		Content: getContentList(img2Contents),
	}

	adds := utils.GetAddedEntries(img1Dir, img2Dir)
	sort.Strings(adds)
	dels := utils.GetDeletedEntries(img1Dir, img2Dir)
	sort.Strings(dels)

	diff = utils.DirDiff{
		Image1: img1,
		Image2: img2,
		Adds:   adds,
		Dels:   dels,
	}
	return diff, nil
}

func getContentList(imgMap map[string]utils.Directory) []string {
	contents := []string{}
	for layer, dir := range imgMap {
		for _, file := range dir.Content {
			contents = append(contents, filepath.Join(layer, file))
		}
	}
	return contents
}

func checkSameLayer(layer1, layer2 utils.Directory) bool {
	layerDiff := utils.DiffDirectory(layer1, layer2)
	same := true
	if len(layerDiff.Adds) != 0 || len(layerDiff.Dels) != 0 {
		same = false
	}
	if len(layerDiff.Mods) != 0 {
		if len(layerDiff.Mods) == 1 && layerDiff.Mods[0] != "/json" {
			same = false
		}
	}
	return same
}

func getImageContents(pathToImage string) (map[string]utils.Directory, error) {
	contents := map[string]utils.Directory{}
	for _, layer := range utils.GetImageLayers(pathToImage) {
		pathToLayer := filepath.Join(pathToImage, layer)
		pathToJSON := layer + ".json"
		err := utils.DirToJSON(pathToLayer, pathToJSON, true)
		if err != nil {
			return contents, fmt.Errorf("Could not convert layer %s in image %s contents to JSON: %s", layer, pathToImage, err)
		}

		layerDir, err := utils.GetDirectory(pathToJSON)
		defer os.Remove(pathToJSON)
		if err != nil {
			return contents, fmt.Errorf("Could not get Directory struct for layer %s in image %s: %s", layer, pathToImage, err)
		}
		contents[layer] = layerDir
	}
	return contents, nil
}
