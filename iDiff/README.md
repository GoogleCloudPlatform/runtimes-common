# iDiff

[![Build
Status](https://travis-ci.org/GoogleCloudPlatform/runtimes-common.svg?branch=master)](https://travis-ci.org/GoogleCloudPlatform/runtimes-common)

## What is iDiff?

iDiff is an image differ command line tool.  iDiff can diff two images along several different criteria, currently including:
- Docker Image History
- Image file system
- apt-get installed packages
- pip installed packages
- npm installed packages

This tool can help you as a developer better understand what is changing within your images and better understand what your images contain.

## Installation

### macOS
```shell
curl -LO iDiff https://storage.googleapis.com/idiff/v0.1.0/iDiff-darwin-amd64
```

### Linux
```shell
curl -LO https://storage.googleapis.com/idiff/v0.1.0/iDiff-linux-amd64 && chmod +x iDiff-linux-amd64 && sudo mv iDiff-linux-amd64 /usr/local/bin/
```

### Windows
Download the [iDiff-windows-amd64.exe](https://storage.googleapis.com/idiff/v0.1.0/iDiff-windows-amd64.exe) file, rename it to `iDiff.exe` and add it to your path


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

```
$ iDiff gcr.io/google-appengine/python:latest gcr.io/google-appengine/python:2017-07-25-110644 -p -a -n

-----Apt Diff-----

Packages found only in gcr.io/google-appengine/python:latest: None

Packages found only in gcr.io/google-appengine/python:2017-07-25-110644: None

Version differences:
        (Package:        gcr.io/google-appengine/python:latest        gcr.io/google-appengine/python:2017-07-25-110644)
        perl:            {5.20.2-3 deb8u7 17584}                      {5.20.2-3 deb8u8 17584}
        
        libgnutls-deb0-28:        {3.3.8-6 deb8u6 1808}        {3.3.8-6 deb8u7 1808}
        
        perl-base:        {5.20.2-3 deb8u7 5098}        {5.20.2-3 deb8u8 5098}
        
        perl-modules:        {5.20.2-3 deb8u7 15108}        {5.20.2-3 deb8u8 15108}
        
-----Node Diff-----

Packages found only in gcr.io/google-appengine/python:latest: None

Packages found only in gcr.io/google-appengine/python:2017-07-25-110644: None

Version differences: None

-----Pip Diff-----

Packages found only in gcr.io/google-appengine/python:latest: None

Packages found only in gcr.io/google-appengine/python:2017-07-25-110644: None

Version differences: None

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
The arguments passed to your differ contain the path to the unpacked tar representation of the image.  That path can be accessed as such: `image1.FSPath`.  

If using existing package differ tools, you should create the appropriate structs to diff (determined in step 2 - either `map[string]map[string]utils.PackageInfo` or `map[string]utils.PackageInfo`) and then call the appropriate get diff function (also determined in step2 - either `GetMultiVerisonMapDiff` or `GetMapDiff`).

Otherwise, create your own differ which should yield information to fill a DiffResult in the next step.

4. Create a DiffResult for your differ if you're not using existing utils or want to wrap the output.  This is where you define how your differ should output for a human readable format and as a struct which can then be written to a `.json` file.  See [output_utils.go](https://github.com/GoogleCloudPlatform/runtimes-common/blob/master/iDiff/utils/output_utils.go).

5. Add your differ to the diffs map in [differs.go](https://github.com/GoogleCloudPlatform/runtimes-common/blob/master/iDiff/differs/differs.go#L22) with the corresponding Differ struct as the value.





