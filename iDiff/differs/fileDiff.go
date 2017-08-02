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
	img1 := image1.FSPath
	img2 := image2.FSPath

	diff, err := diffImageFiles(img1, img2)
	return &utils.DirDiffResult{DiffType: "File Diff", Diff: diff}, err
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
	// }
	return contents, nil
}
