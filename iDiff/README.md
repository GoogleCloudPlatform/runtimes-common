# iDiff

[![Build
Status](https://travis-ci.org/GoogleCloudPlatform/runtimes-common.svg?branch=master)](https://travis-ci.org/GoogleCloudPlatform/runtimes-common)
[![codecov](https://codecov.io/gh/GoogleCloudPlatform/runtimes-common/branch/master/graph/badge.svg)](https://codecov.io/gh/GoogleCloudPlatform/runtimes-common)

## What is iDiff

iDiff is the tool you've been needing but are too busy to make.

iDiff is an image differ command line tool.  iDiff can diff two images on many different levels, currently including:
- Docker Image History
- Image file system
- apt-get installed packages

Additional differs in development:
- pip installed packages

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

## Make your own differ

Feel free to develop your own differ leveraging the utils currently available.  PRs are welcome.


