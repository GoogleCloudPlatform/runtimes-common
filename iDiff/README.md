# iDiff

[![Build
Status](https://travis-ci.org/GoogleCloudPlatform/runtimes-common.svg?branch=master)](https://travis-ci.org/GoogleCloudPlatform/runtimes-common)
[![codecov](https://codecov.io/gh/GoogleCloudPlatform/runtimes-common/branch/master/graph/badge.svg)](https://codecov.io/gh/GoogleCloudPlatform/runtimes-common)

## What is iDiff?

iDiff is an image differ command line tool.  iDiff can diff two images on many different levels, currently including:
- Docker Image History
- Image file system
- apt-get installed packages

Additional differs in development:
- pip installed packages
- npm installed packages

This tool can help you as a developer better understand what is changing within your images and better understand what your images contain.

## Installation

Not sure how to install this...

## Quickstart

To use iDiff you need two Docker images.  Once you have those images you can do any of the following:

```
go run main.go iDiff <img1> <img2> hist
go run main.go iDiff <img1> <img2> dir
go run main.go iDiff <img1> <img2> apt
```

## Piping output

To get a JSON version of the iDiff output add a `-j` or `-json` flag.

```go run main.go iDiff <img1> <img2> <differ> -j```

## Known issues

Docker with API Version 1.29 is currently required.  If either of those make you unable to use the tool, in development is a way to run iDiff on tarballs in order to bypass any Docker dependency. 

## Make your own differ

Feel free to develop your own differ leveraging the utils currently available.  PRs are welcome.




