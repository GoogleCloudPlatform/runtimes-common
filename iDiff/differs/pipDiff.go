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
func (d PipDiffer) Diff(image1, image2 utils.Image) (utils.DiffResult, error) {
	diff, err := multiVersionDiff(image1, image2, d)
	return diff, err
}

func getPythonVersion(pathToLayer string) ([]string, error) {
	matches := []string{}
	libPath := filepath.Join(pathToLayer, "usr/local/lib")
	libContents, err := ioutil.ReadDir(libPath)
	if err != nil {
		return matches, err
	}

	for _, file := range libContents {
		pattern := regexp.MustCompile("^python[0-9]+\\.[0-9]+$")
		match := pattern.FindString(file.Name())
		if match != "" {
			matches = append(matches, match)
		}
	}
	return matches, nil
}

func (d PipDiffer) getPackages(path string) (map[string]map[string]utils.PackageInfo, error) {
	packages := make(map[string]map[string]utils.PackageInfo)

	pythonVersions, err := getPythonVersion(path)
	if err != nil {
		// layer doesn't have a Python version installed
		return packages, nil
	}
	for _, pyVersion := range pythonVersions {
		packagesPath := filepath.Join(path, "usr/local/lib", pyVersion, "site-packages")
		contents, err := ioutil.ReadDir(packagesPath)
		if err != nil {
			// layer's Python folder doesn't have a site-packages folder
			return packages, nil
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
				currPackage := utils.PackageInfo{Version: version, Size: size}
				addToMap(packages, packageName, pyVersion, currPackage)
			}

			// if not package, check if Python file
			pythonFile := regexp.MustCompile(".+\\.py$")
			fileMatch := pythonFile.FindString(fileName)
			if fileMatch != "" {
				currPackage := utils.PackageInfo{}
				addToMap(packages, fileName, pyVersion, currPackage)
			}
		}
	}

	return packages, nil
}

func addToMap(packages map[string]map[string]utils.PackageInfo, pack string, pyVersion string, packInfo utils.PackageInfo) {
	if _, ok := packages[pack]; !ok {
		// package not yet seen
		infoMap := make(map[string]utils.PackageInfo)
		infoMap[pyVersion] = packInfo
		packages[pack] = infoMap
		return
	}
	packages[pack][pyVersion] = packInfo
}
