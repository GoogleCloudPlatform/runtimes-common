package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
)

type Image struct {
	Source  string
	FSPath  string
	History []string
	Layers  []string
}

type ImagePrepper struct {
	Source    string
	UseDocker bool
}

type Prepper interface {
	ImageToFS() (string, error)
}

func (p ImagePrepper) GetImage() (Image, error) {
	img := p.Source

	var prepper Prepper
	if CheckImageID(img) {
		prepper = IDPrepper{p}
	} else if CheckImageURL(img) {
		prepper = CloudPrepper{p}
	} else if CheckTar(img) {
		prepper = TarPrepper{p}
	} else {
		return Image{}, errors.New("Could not retrieve image from source")
	}

	imgPath, err := prepper.ImageToFS()
	if err != nil {
		return Image{}, err
	}

	history, err := getHistoryList(p.Source, p.UseDocker)
	if err != nil {
		return Image{}, err
	}

	return Image{
		Source:  img,
		FSPath:  imgPath,
		History: history,
	}, nil
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
	valid, err := ValidDockerVersion(p.UseDocker)
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
	valid, err := ValidDockerVersion(p.UseDocker)
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
