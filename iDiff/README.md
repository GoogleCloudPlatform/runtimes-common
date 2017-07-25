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

Not sure how to install this...

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

## Piping output

To get a JSON version of the iDiff output add a `-j` or `-json` flag.

```iDiff <img1> <img2> -j```

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
 

## Make your own differ

Feel free to develop your own differ leveraging the utils currently available.  PRs are welcome.




