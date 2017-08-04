package differs

import (
	"fmt"
	"os"
	"sort"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
)

type FileDiffer struct {
}

// FileDiff diffs two packages and compares their contents
func (d FileDiffer) Diff(image1, image2 utils.Image) (utils.DiffResult, error) {
	diff, err := diffImageFiles(image1, image2)
	return &utils.DirDiffResult{DiffType: "FileDiffer", Diff: diff}, err
}

func diffImageFiles(image1, image2 utils.Image) (utils.DirDiff, error) {
	img1 := image1.FSPath
	img2 := image2.FSPath

	var diff utils.DirDiff

	target1 := "j1.json"
	err := utils.DirToJSON(img1, target1, true)
	if err != nil {
		return diff, err
	}
	target2 := "j2.json"
	err = utils.DirToJSON(img2, target2, true)
	if err != nil {
		return diff, err
	}
	img1Dir, err := utils.GetDirectory(target1)
	if err != nil {
		return diff, err
	}
	img2Dir, err := utils.GetDirectory(target2)
	if err != nil {
		return diff, err
	}

	adds := utils.GetAddedEntries(img1Dir, img2Dir)
	sort.Strings(adds)
	dels := utils.GetDeletedEntries(img1Dir, img2Dir)
	sort.Strings(dels)

	diff = utils.DirDiff{
		Image1: image1.Source,
		Image2: image2.Source,
		Adds:   adds,
		Dels:   dels,
	}
	return diff, nil
}

func getContentList(imgMap map[string]utils.Directory) []string {
	contents := []string{}
	for _, dir := range imgMap {
		for _, file := range dir.Content {
			contents = append(contents, file)
		}
	}
	return contents
}

func getImageContents(pathToImage string) (map[string]utils.Directory, error) {
	contents := map[string]utils.Directory{}
	pathToJSON := pathToImage + ".json"
	err := utils.DirToJSON(pathToImage, pathToJSON, true)
	if err != nil {
		return contents, fmt.Errorf("Could not convert layer %s in image %s contents to JSON: %s", pathToImage, pathToImage, err)
	}

	layerDir, err := utils.GetDirectory(pathToJSON)
	defer os.Remove(pathToJSON)
	if err != nil {
		return contents, fmt.Errorf("Could not get Directory struct for layer %s in image %s: %s", pathToImage, pathToImage, err)
	}
	contents[pathToImage] = layerDir
	return contents, nil
}
