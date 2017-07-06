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

Additional differs in development:
- npm installed packages

This tool can help you as a developer better understand what is changing within your images and better understand what your images contain.

## Installation

Not sure how to install this...

## Quickstart

To use iDiff you need two Docker images (in the form of an ID, tarball, or URL from a repo).  Once you have those images you can run any of the following differs:

```
go run main.go iDiff history <img1> <img2>
go run main.go iDiff file <img1> <img2>
go run main.go iDiff apt <img1> <img2>
go run main.go iDiff pip <img1> <img2>
```

## Piping output

To get a JSON version of the iDiff output add a `-j` or `-json` flag.

```go run main.go iDiff <differ> <img1> <img2> -j```

## Known issues

To run iDiff on image IDs or URLs, Docker Engine on the client and server side must have compatible versions, such as by modifying the Docker API version used at runtime on the client side to match that of the server or modifying the server or client side Docker engine to match.  Currently in development is a way to run iDiff that bypasses Docker client/server model dependencies. 

## Make your own differ

Feel free to develop your own differ leveraging the utils currently available.  PRs are welcome.




