package differs

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
	"github.com/golang/glog"
)

type PipDiffer struct {
}

// PipDiff compares pip-installed Python packages between layers of two different images.
func (d PipDiffer) Diff(image1, image2 utils.Image) (utils.DiffResult, error) {
	diff, err := singleVersionDiff(image1, image2, d)
	return diff, err
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

func (d PipDiffer) getPackages(path string) (map[string]utils.PackageInfo, error) {
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

		for i := 0; i < len(contents); i++ {
			c := contents[i]
			fileName := c.Name()

			// check if package
			packageDir := regexp.MustCompile("^([a-z|A-Z]+)-(([0-9]+?\\.){3})dist-info$")
			packageMatch := packageDir.FindStringSubmatch(fileName)
			if len(packageMatch) != 0 {
				packageName := packageMatch[1]
				version := packageMatch[2][:len(packageMatch[2])-1]

				// Retrieves size for actual package/script corresponding to each dist-info metadata directory
				// by taking the file entry alphabetically before it (for a package) or after it (for a script)
				var size string
				if contents[i-1].Name() == packageName {
					packagePath := filepath.Join(packagesPath, packageName)
					intSize, err := utils.GetDirectorySize(packagePath)
					if err != nil {
						glog.Errorf("Could not obtain size for package %s", packagePath)
						size = ""
					} else {
						size = strconv.FormatInt(intSize, 10)
					}
				} else if contents[i+1].Name() == packageName+".py" {
					size = strconv.FormatInt(contents[i+1].Size(), 10)

				} else {
					glog.Errorf("Could not find Python package %s for corresponding metadata info", packageName)
					continue
				}

				packages[packageName] = utils.PackageInfo{Version: version, Size: size}
			}
		}
	}

	return packages, nil
}
