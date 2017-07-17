package differs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
)

// Package diffs two packages and compares their contents
func Package(img1, img2 string, json bool, eng bool) (string, error) {
	diff, err := diffImageFiles(img1, img2, eng)
	if err != nil {
		return "", err
	}

	output, err := getDiffOutput(diff, json)
	if err != nil {
		return "", err
	}
	return output, nil
}

func getDiffOutput(dirDiff utils.DirDiff, json bool) (string, error) {
	if json {
		return utils.JSONify(dirDiff)
	}

	var buffer bytes.Buffer

	s := fmt.Sprintf("These entries have been added to %s\n", dirDiff.Image1)
	buffer.WriteString(s)
	if len(dirDiff.Adds) == 0 {
		buffer.WriteString("\tNo files have been added\n")
	} else {
		for _, f := range dirDiff.Adds {
			s = fmt.Sprintf("\t%s\n", f)
			buffer.WriteString(s)
		}
	}

	s = fmt.Sprintf("These entries have been deleted from %s\n", dirDiff.Image1)
	buffer.WriteString(s)
	if len(dirDiff.Dels) == 0 {
		buffer.WriteString("\tNo files have been deleted\n")
	} else {
		for _, f := range dirDiff.Dels {
			s = fmt.Sprintf("\t%s\n", f)
			buffer.WriteString(s)
		}
	}
	s = fmt.Sprintf("These entries have been changed between %s and %s\n", dirDiff.Image1, dirDiff.Image2)
	buffer.WriteString(s)
	if len(dirDiff.Mods) == 0 {
		buffer.WriteString("\tNo files have been modified\n")
	} else {
		for _, f := range dirDiff.Mods {
			s = fmt.Sprintf("\t%s\n", f)
			buffer.WriteString(s)
		}
	}

	return buffer.String(), nil
}

func diffImageFiles(img1, img2 string, eng bool) (utils.DirDiff, error) {
	var diff utils.DirDiff
	img1FS, err := utils.ImageToFS(img1, eng)
	if err != nil {
		return diff, fmt.Errorf("Error retrieving image %s file system: %s", img1, err)
	}
	img2FS, err := utils.ImageToFS(img2, eng)
	if err != nil {
		return diff, fmt.Errorf("Error retrieving image %s file system: %s", img2, err)
	}

	img1Contents, err := getImageContents(img1FS)
	if err != nil {
		return diff, fmt.Errorf("Error parsing image %s contents: %s", img1, err)
	}
	img2Contents, err := getImageContents(img2FS)
	if err != nil {
		return diff, fmt.Errorf("Error parsing image %s contents: %s", img2, err)
	}
	defer os.RemoveAll(img1FS)
	defer os.RemoveAll(img2FS)

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
	diff = utils.DiffDirectory(img1Dir, img2Dir)
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
