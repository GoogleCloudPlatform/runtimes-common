package differs

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
	"github.com/golang/glog"
)

// layers of two different images.
func PipDiff(d1file, d2file string, json bool, eng bool) (string, error) {
	d1, err := utils.GetDirectory(d1file)
	if err != nil {
		glog.Errorf("Error reading directory structure from file %s: %s\n", d1file, err)
		return "", err
	}
	d2, err := utils.GetDirectory(d2file)
	if err != nil {
		glog.Errorf("Error reading directory structure from file %s: %s\n", d2file, err)
		return "", err
	}

	dirPath1 := d1.Root
	dirPath2 := d2.Root
	pack1 := getPythonPackages(dirPath1)
	pack2 := getPythonPackages(dirPath2)

	diff := utils.GetMapDiff(pack1, pack2)
	diff.Image1 = dirPath1
	diff.Image2 = dirPath2
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

// TODO: Eventually, this would make use of the shallow JSON and be diffed
// with that of another image to get only the layers that have changed.
func getImageLayers(pathToImage string) []string {
	layers := []string{}
	contents, err := ioutil.ReadDir(pathToImage)
	if err != nil {
		glog.Error(err.Error())
	}

	for _, file := range contents {
		if file.IsDir() {
			layers = append(layers, file.Name())
		}
	}
	return layers
}

func getPythonPackages(path string) map[string]utils.PackageInfo {
	packages := make(map[string]utils.PackageInfo)

	layers := getImageLayers(path)
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
