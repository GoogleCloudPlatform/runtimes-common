package differs

import (
	"reflect"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
)

type MultiVersionPackageDiffer interface {
	getPackages(path string) (map[string]map[string]utils.PackageInfo, error)
}

type SingleVersionPackageDiffer interface {
	getPackages(path string) (map[string]utils.PackageInfo, error)
}

func multiVersionDiff(image1, image2 utils.Image, differ MultiVersionPackageDiffer) (utils.DiffResult, error) {
	img1FS := image1.FSPath
	img2FS := image2.FSPath

	pack1, err := differ.getPackages(img1FS)
	if err != nil {
		return &utils.MultiVersionPackageDiffResult{}, err
	}
	pack2, err := differ.getPackages(img2FS)
	if err != nil {
		return &utils.MultiVersionPackageDiffResult{}, err
	}

	diff := utils.GetMultiVersionMapDiff(pack1, pack2, image1.Source, image2.Source)
	diff.DiffType = reflect.TypeOf(differ).Name()
	return &diff, nil
}

func singleVersionDiff(image1, image2 utils.Image, differ SingleVersionPackageDiffer) (utils.DiffResult, error) {
	img1FS := image1.FSPath
	img2FS := image2.FSPath

	pack1, err := differ.getPackages(img1FS)
	if err != nil {
		return &utils.PackageDiffResult{}, err
	}
	pack2, err := differ.getPackages(img2FS)
	if err != nil {
		return &utils.PackageDiffResult{}, err
	}

	diff := utils.GetMapDiff(pack1, pack2, image1.Source, image2.Source)
	diff.DiffType = reflect.TypeOf(differ).Name()
	return &diff, nil
}
