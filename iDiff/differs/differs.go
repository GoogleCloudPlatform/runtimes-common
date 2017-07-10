package differs

import (
	"errors"
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
	"github.com/docker/docker/client"
)

var diffs = map[string]func(string, string, bool) (string, error){
	"hist": History,
	"dir":  Package,
	"apt":  AptDiff,
}

func Diff(arg1, arg2, differ string, json bool) (string, error) {
	if f, exists := diffs[differ]; exists {
		validDocker, err := validDockerVersion()
		if err != nil {
			return "", err
		}
		if differ == "hist" {
			return f(arg1, arg2, json)
		}
		return specificDiffer(f, arg1, arg2, json, validDocker)
	}
	return "", errors.New("Unknown differ")
}

func validDockerVersion() (bool, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return false, fmt.Errorf("Docker client error: %s", err)
	}
	version := cli.ClientVersion()
	if version == "1.31" {
		return true, nil
	}
	return false, nil
}

func imageToDirCmd(image string) error {
	cmdName := "docker"
	cmdArgs := []string{"save", image, ">", image + ".tar"}

	cmd := exec.command(cmdName, cmdArgs...)
	err := cmd.Start()
	if err != nil {
		return err
	}

}

func prepareDir(image string, validDocker bool) (string, string, error) {
	if validDocker {
		return utils.ImageToDir(image)
	}
	return "", "", imageToDirCmd(image) // TODO add exec calls for local docker client
	// return "", "", nil
}

func specificDiffer(f func(string, string, bool) (string, error), img1, img2 string, json, validDocker bool) (string, error) {
	jsonPath1, dirPath1, err := prepareDir(img1, validDocker)
	if err != nil {
		return "", err
	}
	jsonPath2, dirPath2, err := prepareDir(img2, validDocker)
	if err != nil {
		return "", err
	}
	diff, err := f(jsonPath1, jsonPath2, json)
	if err != nil {
		return "", err
	}

	errStr := remove(dirPath1, true, "")
	errStr = remove(dirPath2, true, errStr)
	errStr = remove(jsonPath1, false, errStr)
	errStr = remove(jsonPath2, false, errStr)

	if errStr != "" {
		return diff, errors.New(errStr)
	}

	return diff, err
}

func remove(path string, dir bool, errStr string) string {
	var err error
	if dir {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}
	if err != nil {
		errStr += "\nUnable to remove " + path
	}
	return errStr
}
