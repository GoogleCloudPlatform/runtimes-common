package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/golang/glog"
)

var sourceToPrepMap = map[string]Prepper{
	"ID":  IDPrepper{},
	"URL": CloudPrepper{},
	"tar": TarPrepper{},
}

var sourceCheckMap = map[string]func(string) bool{
	"ID":  CheckImageID,
	"URL": CheckImageURL,
	"tar": CheckTar,
}

type Image struct {
	Source  string
	FSPath  string
	History []string
	Layers  []string
}

type ImagePrepper struct {
	Source string
}

type Prepper interface {
	ImageToFS() (string, error)
}

func (p ImagePrepper) GetImage() (Image, error) {
	img := p.Source

	var prepper Prepper
	for source, check := range sourceCheckMap {
		if check(img) {
			typePrepper := reflect.TypeOf(sourceToPrepMap[source])
			prepper = reflect.New(typePrepper).Interface().(Prepper)
			reflect.ValueOf(prepper).Elem().Field(0).Set(reflect.ValueOf(p))
			break
		}
	}
	if prepper == nil {
		return Image{}, errors.New("Could not retrieve image from source")
	}

	imgPath, err := prepper.ImageToFS()
	if err != nil {
		return Image{}, err
	}

	history, err := getHistory(imgPath)
	if err != nil {
		return Image{}, err
	}

	return Image{
		Source:  img,
		FSPath:  imgPath,
		History: history,
	}, nil
}

type histJSON struct {
	History []histLayer `json:"history"`
}

type histLayer struct {
	Created    string `json:"created"`
	CreatedBy  string `json:"created_by"`
	EmptyLayer bool   `json:"empty_layer"`
}

func getHistory(imgPath string) ([]string, error) {
	contents, err := BuildLayerTargets(imgPath, "")
	if err != nil {
		return []string{}, err
	}
	for _, item := range contents {
		if filepath.Ext(item) == ".json" {
			file, err := ioutil.ReadFile(item)
			if err != nil {
				return []string{}, err
			}
			var histJ histJSON
			json.Unmarshal(file, &histJ)
			historyList := []string{}
			for _, layer := range histJ.History {
				historyList = append(historyList, layer.CreatedBy)
			}
			return historyList, nil
		}
	}
	return []string{}, nil
}

func getImageFromTar(tarPath string) (string, error) {
	err := ExtractTar(tarPath)
	if err != nil {
		return "", err
	}
	path := strings.TrimSuffix(tarPath, filepath.Ext(tarPath))
	return path, nil
}

// CloudPrepper prepares images sourced from a Cloud registry
type CloudPrepper struct {
	ImagePrepper
}

func (p CloudPrepper) ImageToFS() (string, error) {
	// check client compatibility with Docker API
	valid, err := ValidDockerVersion()
	if err != nil {
		return "", err
	}
	var tarPath string
	if !valid {
		glog.Info("Docker version incompatible with api, shelling out to local Docker client.")
		imageID, imageName, err := pullImageCmd(p.Source)
		if err != nil {
			return "", err
		}
		tarPath, err = imageToTarCmd(imageID, imageName)
	} else {
		imageID, imageName, err := pullImageFromRepo(p.Source)
		if err != nil {
			return "", err
		}
		tarPath, err = saveImageToTar(imageID, imageName)
	}
	if err != nil {
		return "", err
	}

	defer os.Remove(tarPath)
	return getImageFromTar(tarPath)
}

type IDPrepper struct {
	ImagePrepper
}

func (p IDPrepper) ImageToFS() (string, error) {
	// check client compatibility with Docker API
	valid, err := ValidDockerVersion()
	if err != nil {
		return "", err
	}
	var tarPath string
	if !valid {
		glog.Info("Docker version incompatible with api, shelling out to local Docker client.")
		tarPath, err = imageToTarCmd(p.Source, p.Source)
	} else {
		tarPath, err = saveImageToTar(p.Source, p.Source)
	}
	if err != nil {
		return "", err
	}

	defer os.Remove(tarPath)
	return getImageFromTar(tarPath)
}

type TarPrepper struct {
	ImagePrepper
}

func (p TarPrepper) ImageToFS() (string, error) {
	return getImageFromTar(p.Source)
}
