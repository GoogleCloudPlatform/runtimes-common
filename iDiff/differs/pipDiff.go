package differs

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
)

type PipDiffer struct {
}

// PipDiff compares pip-installed Python packages between layers of two different images.
func (d PipDiffer) Diff(image1, image2 utils.Image, json, eng bool) (string, error) {
	img1 := image1.FSPath
	img2 := image2.FSPath

	pack1 := getPythonPackages(img1)
	pack2 := getPythonPackages(img2)

	diff := utils.GetMapDiff(pack1, pack2)
	diff.Image1 = img1
	diff.Image2 = img2

	if json {
		return utils.JSONify(diff)
	}
	utils.Output(diff)
	return "", nil
}

func getPythonVersion(pathToLayer string) (string, bool) {
	libPath := filepath.Join(pathToLayer, "/layer/usr/local/lib")
	libContents, err := ioutil.ReadDir(libPath)
	if err != nil {
		return "", false
	}

	for _, file := range libContents {
		pattern := regexp.MustCompile("^python[0-9]+\\.[0-9]+$")
		match := pattern.FindString(file.Name())
		if match != "" {
			return match, true
		}
	}
	return "", false
}

func getPythonPackages(path string) map[string]utils.PackageInfo {
	packages := make(map[string]utils.PackageInfo)

	// TODO: Eventually, this would make use of the shallow JSON and be diffed
	// with that of another image to get only the layers that have changed.
	layers := utils.GetImageLayers(path)
	for _, layer := range layers {
		pathToLayer := filepath.Join(path, layer)
		pythonVersion, exists := getPythonVersion(pathToLayer)
		if !exists {
			// layer doesn't have a Python folder installed
			continue
		}
		packagesPath := filepath.Join(pathToLayer, "layer/usr/local/lib", pythonVersion, "site-packages")
		contents, err := ioutil.ReadDir(packagesPath)
		if err != nil {
			// layer's Python folder doesn't have a site-packages folder
			continue
		}

		for _, c := range contents {
			fileName := c.Name()

			// check if package
			packageDir := regexp.MustCompile("^([a-z|A-Z]+)-(([0-9]+?\\.){3})dist-info$")
			packageMatch := packageDir.FindStringSubmatch(fileName)
			if len(packageMatch) != 0 {
				packageName := packageMatch[1]
				version := packageMatch[2][:len(packageMatch[2])-1]
				size := strconv.FormatInt(c.Size(), 10)
				packages[packageName] = utils.PackageInfo{version, size}

				continue
			}

			// if not package, check if Python file
			pythonFile := regexp.MustCompile(".+\\.py$")
			fileMatch := pythonFile.FindString(fileName)

			if fileMatch != "" {
				packages[fileName] = utils.PackageInfo{}
			}
		}
	}

	return packages
}
