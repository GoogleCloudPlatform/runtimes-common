package differs

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// AptDiff compares the packages installed by apt-get.
func AptDiff(img1, img2 string) (string, error) {
	pack1, err := getPackages(img1)
	if err != nil {
		return "", err
	}
	pack2, err := getPackages(img2)
	if err != nil {
		return "", err
	}

	diff1, diff2 := diffMaps(pack1, pack2)
	s1 := fmt.Sprintf("Image %s had the following packages which differed:\n%s\n", img1, strings.Join(diff1, "\n"))
	s2 := fmt.Sprintf("\nImage %s had the following packages which differed:\n%s", img2, strings.Join(diff2, "\n"))
	return s1 + s2, nil
}

func diffMaps(map1, map2 map[string]string) ([]string, []string) {
	diff1 := []string{}
	diff2 := []string{}
	for key1, value1 := range map1 {
		value2, ok := map2[key1]
		if !ok || value2 != value1 {
			diff1 = append(diff1, key1+":"+value1)
		} else {
			delete(map2, key1)
		}
	}
	for key2, value2 := range map2 {
		diff2 = append(diff2, key2+":"+value2)
	}
	return diff1, diff2
}

func getPackages(path string) (map[string]string, error) {
	packages := make(map[string]string)

	var layerStems []string

	layers, err := ioutil.ReadDir(path)
	if err != nil {
		return packages, err
	}
	for _, layer := range layers {
		layerStems = append(layerStems, filepath.Join(path, layer.Name(), "layer/var/lib/dpkg/status"))
	}

	for _, statusFile := range layerStems {
		if _, err := os.Stat(statusFile); err == nil {

			if file, err := os.Open(statusFile); err == nil {
				// make sure it gets closed
				defer file.Close()

				var currPackage string
				// create a new scanner and read the file line by line
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := strings.Split(scanner.Text(), ": ")
					if len(line) == 2 {
						key := line[0]
						value := line[1]
						if key == "Package" {
							currPackage = value
						}
						if key == "Version" {
							packages[currPackage] = value
						}
					}

				}

			} else {
				return packages, err
			}
		} else {
			// status file does not exist in this layer
			continue
		}

	}
	return packages, nil
}
