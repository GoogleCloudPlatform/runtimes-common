# iDiff

[![Build
Status](https://travis-ci.org/GoogleCloudPlatform/runtimes-common.svg?branch=master)](https://travis-ci.org/GoogleCloudPlatform/runtimes-common)
[![codecov](https://codecov.io/gh/GoogleCloudPlatform/runtimes-common/branch/master/graph/badge.svg)](https://codecov.io/gh/GoogleCloudPlatform/runtimes-common)

## What is iDiff?

iDiff is an image differ command line tool.  iDiff can diff two images along several different criteria, currently including:
- Docker Image History
- Image file system
- apt-get installed packages
- pip installed packages
- npm installed packages

This tool can help you as a developer better understand what is changing within your images and better understand what your images contain.

## Installation

TODO: add how to install

## Quickstart

To use iDiff you need two Docker images (in the form of an ID, tarball, or URL from a repo).  Once you have those images you can run any of the following differs:

```
iDiff <img1> <img2>     [Run all differs]
iDiff <img1> <img2> -d  [History]
iDiff <img1> <img2> -f  [File System]
iDiff <img1> <img2> -p  [Pip]
iDiff <img1> <img2> -a  [Apt]
iDiff <img1> <img2> -n  [Node]
```

You can similarly run many differs at once:

```
iDiff <img1> <img2> -d -a -n [History, Apt, and Node]
```
All of the differ flags with their long versions can be seen below:

| Differ                    | Short flag | Long Flag  |
| ------------------------- |:----------:| ----------:|
| File System diff          | -f         | --file     |
| History                   | -d 	 | --history  |
| npm installed packages    | -n 	 | --node     |
| pip installed packages    | -p 	 | --pip      |
| apt-get installed packages| -a 	 | --apt      |




## Other Flags

To get a JSON version of the iDiff output add a `-j` or `--json` flag.

```iDiff <img1> <img2> -j```

To use the docker client instead of shelling out to your local docker daemon, add a `-e` or `--eng` flag.

```iDiff <img1> <img2> -e```


## Output Format

### History Diff

The history differ has the following json output structure:

```
type HistDiff struct {
	Image1 string
	Image2 string
	Adds   []string
	Dels   []string
}
```

### File System Diff

The files system differ has the following json output structure: 

```
type DirDiff struct {
	Image1 string
	Image2 string
	Adds   []string
	Dels   []string
	Mods   []string
}
```

### Package Diffs

Package differs such as pip, apt, and node inspect the packages contained within the images provided.  All packages differs currently leverage the PackageInfo struct which contains the version and size for a given package instance.

```
type PackageInfo struct {
	Version string
	Size    string
}
```

#### Single Version Diffs

The single version differs (apt, pip) have the following json output structure:

```
type PackageDiff struct {
	Image1    string
	Packages1 map[string]PackageInfo
	Image2    string
	Packages2 map[string]PackageInfo
	InfoDiff  []Info
}
```

Image1 and Image2 are the image names.  Packages1 and Packages2 map package names to PackageInfo structs which contain the version and size of the package.  InfoDiff contains a list of Info structs, each of which contains the package name (which occurred in both images but had a difference in size or version), and the PackageInfo struct for each package instance. 

#### Multi Version Diffs

The multi version differs (node) support processing images which may have multiple versions of the same package.  Below is the json output structure:

```
type MultiVersionPackageDiff struct {
	Image1    string
	Packages1 map[string]map[string]PackageInfo
	Image2    string
	Packages2 map[string]map[string]PackageInfo
	InfoDiff  []MultiVersionInfo
}
```

Image1 and Image2 are the image names.  Packages1 and Packages2 map package name to path where the package was found to PackageInfo struct (version and size of that package instance).  InfoDiff here is exanded to allow for multiple versions to be associated with a single package.

```
type MultiVersionInfo struct {
	Package string
	Info1   []PackageInfo
	Info2   []PackageInfo
}
```

## Known issues

To run iDiff on image IDs or URLs, docker must be installed.

## Example Run

// TODO: update when output format finalized


```
$ iDiff  gcr.io/google-appengine/python:latest gcr.io/google-appengine/python:2017-07-25-110644 -p -d -a -n
Packages found only in pythonlatest:

Packages found only in python2017-07-25-110644:

Version differences:
	(Package:	pythonlatest		python2017-07-25-110644)
	perl-modules:	{5.20.2-3 deb8u7 15108}	{5.20.2-3 deb8u8 15108}

	libgnutls-deb0-28:	{3.3.8-6 deb8u6 1808}	{3.3.8-6 deb8u7 1808}

	perl:	{5.20.2-3 deb8u7 17584}	{5.20.2-3 deb8u8 17584}

	perl-base:	{5.20.2-3 deb8u7 5098}	{5.20.2-3 deb8u8 5098}


Docker history lines found only in gcr.io/google-appengine/python:latest:
-/bin/sh -c #(nop) ADD file:ddbbbee34af5bc54fda5da491a14d8367a072190f63f1f44c62f1712ca14b2fc in /
-/bin/sh -c #(nop) ADD dir:b64d4dbd411116b57f00622e8e35785469c45241d96139b2c2c282b76997c4ff in /scripts
-/bin/sh -c #(nop) ADD dir:a61d73142ff699ef209e79d5ff0331ee0c1026d1a254789be57dfcb47424b9b9 in /resources

Docker history lines found only in gcr.io/google-appengine/python:2017-07-25-110644:
-/bin/sh -c #(nop) ADD file:9b537477c1fe03a9a3af141199b8848b9718bbc259e9d040e52dd78d9b1472a0 in /
-/bin/sh -c #(nop) ADD dir:b209be879a64da94090efa46c7647cf4a972d9233219a86718dea815b8b6ea62 in /scripts
-/bin/sh -c #(nop) ADD dir:de273ffbfc6ea318de85b56ed05cd1d002b5b0bfa5c721a42b7aa8d44ff60c42 in /resources

Packages found only in pythonlatest:

Packages found only in python2017-07-25-110644:

Version differences:
	(Package:	pythonlatest		python2017-07-25-110644)

Packages found only in pythonlatest:

Packages found only in python2017-07-25-110644:

Version differences:
	(Package:	pythonlatest		python2017-07-25-110644)

```


## Make your own differ

Feel free to develop your own differ leveraging the utils currently available.  PRs are welcome.

### Custom Differ Quickstart

In order to quickly make your own differ, follow these steps:

1. Add your diff identifier to the flags in [root.go](https://github.com/GoogleCloudPlatform/runtimes-common/blob/ReadMe/iDiff/cmd/root.go)
2. Determine if you can use existing differ tools.  If you can make use of existing tools, you then need to construct the structs to feed to the diff tools by getting all of the packages for each image or the analogous quality to be diffed.  To determine if you can leverage existing tools, think through these questions:
- Are you trying to diff packages?
    - Yes: Does the relevant package manager support different versions of the same package on one image?
        - Yes: Use `GetMultiVerisonMapDiff` to diff `map[string]map[string]utils.PackageInfo` objects.  See [nodeDiff.go](https://github.com/GoogleCloudPlatform/runtimes-common/blob/master/iDiff/differs/nodeDiff.go#L33) for an example.
        -  No: Use `GetMapDiff` to diff `map[string]utils.PackageInfo` objects.  See [aptDiff.go](https://github.com/GoogleCloudPlatform/runtimes-common/blob/master/iDiff/differs/aptDiff.go#L29) or [pipDiff.go](https://github.com/GoogleCloudPlatform/runtimes-common/blob/master/iDiff/differs/pipDiff.go#L23) for examples. 
    - No: Look to [History](https://github.com/GoogleCloudPlatform/runtimes-common/blob/ReadMe/iDiff/differs/historyDiff.go) and [File System](https://github.com/GoogleCloudPlatform/runtimes-common/blob/ReadMe/iDiff/differs/fileDiff.go) differs as models for diffing.

3. Write your Diff driver such that you have a struct for your differ type and a method for that differ called Diff:

```
type YourDiffer struct {}

func (d YourDiffer) Diff(image1, image2 utils.Image) (DiffResult, error) {...}
```
The arguments passed to your differ contain the path to the unpacked tar representation of the image.  That path can be accessed as such: `image1.FSPath`.  Given that path you should create the appropriate struct (determined in step 2) and then call the appropriate get-diff function (also determined in step2).

4. Create a DiffResult for your differ if you're not using existing utils or want to wrap the output.  This is where you define how your differ should output for a human readable format and as a json.  See [output_utils.go](https://github.com/GoogleCloudPlatform/runtimes-common/blob/master/iDiff/utils/output_utils.go).

5. Add your differ to the diffs map in [differs.go](https://github.com/GoogleCloudPlatform/runtimes-common/blob/master/iDiff/differs/differs.go#L22) with the corresponding Differ struct as the value.





