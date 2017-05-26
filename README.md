# plurality

Simple proof of concept wrapper for runC to demonstrate ([http://runc.io/](runC)) can be an alternative to ([http://singularity.lbl.gov/](Singularity)) backed by Docker

## Description

A command line application that is a wrapper for runC to facilitate:
* creation of file system directory based containers based on Docker images
* launching of container in ([https://www.cyphar.com/blog/post/rootless-containers-with-runc](rootless ))mode

## Usage

It is just a PoC - so it will blow up on you, but a simple example might be to create and run a container based on a docker image such as:
```
# create a container based on ubuntu:16.04
./plurality create ubuntu:16.04  echo_test
# run something inside the container
./plurality run  echo_test echo "this is a test wahoo"
```

The smoke and mirrors part of this PoC is that it uses sudo and tar to unpack the exported images from the Docker image client API.  This is because root is required to create make nodes.

All containers are stored in ${HOME}/.runc .


## Install

To install, use `go get`:

```bash
# no sophisticated packaging and dependency resolution here
# install vndr from https://github.com/LK4D4/vndr
$ go get github.com/LK4D4/vndr
$ go get -d github.com/piersharding/plurality
$ cd ${GOPATH}/src/github.com/piersharding/plurality
# sort out the dependencies
$ vndr
# now build
$ go build .
```

or just simply do:
```bash
$ make
``

## Contribution

1. Fork ([https://github.com/piersharding/plurality/fork](https://github.com/piersharding/plurality/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[piersharding](https://github.com/piersharding)
