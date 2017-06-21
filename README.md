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
```


## To Test

After following above "Install" instructions a simple single executable container can be tested with the following:

```bash
# no sophisticated packaging and dependency resolution here
# install vndr from https://github.com/LK4D4/vndr
$ make clean
$ make test
```

You will get output similar to:
```bash
sudo rm -rf /home/piers/.runc/container/echo_test
rm -f plurality example/echo-server
docker rm -f echo_test 2>/dev/null || true
docker rmi echo-server 2>/dev/null || true
Untagged: echo-server:latest
Deleted: sha256:47b9872bcf5da8c2b34a4c7eaeffb1274c10956ae06e9abc29cfc2eec954aa2c
Deleted: sha256:5c695b2dd5040af41115fa5dce5ba6578c91a988ab794316f4e73d7ca8ca6983
Deleted: sha256:d6132d897daaadf76da4365c6d3a2cb2ccb3b59057147b2ec216f674177751a6
go build -i -o plurality .
cd example && CGO_ENABLED=0 go build  -a -ldflags '-s' -o echo-server server.go
cd example && docker build -t echo-server .
Sending build context to Docker daemon  1.716MB
Step 1/3 : FROM scratch
 ---> 
Step 2/3 : ADD ./echo-server /echo-server
 ---> 7ff2953abdf4
Removing intermediate container a627ea13319f
Step 3/3 : CMD /echo-server
 ---> Running in bfc795914c28
 ---> 90b67d76f415
Removing intermediate container bfc795914c28
Successfully built 90b67d76f415
Successfully tagged echo-server:latest
./plurality --debug --nosudo create --nopull echo-server:latest echo_test
2017/06/22 09:11:40 09:11:40.262 CmdCreate ▶ DEBU 001 Creating container  echo_test  from  echo-server:latest
2017/06/22 09:11:40 09:11:40.262 fileExits ▶ DEBU 002 checking path  /home/piers/.runc/container/echo_test
2017/06/22 09:11:40 09:11:40.262 CmdCreate ▶ DEBU 003 runC is available at  /usr/local/sbin/runc
2017/06/22 09:11:40 09:11:40.262 CmdCreate ▶ DEBU 004 sudo is available at  /usr/bin/sudo
2017/06/22 09:11:40 09:11:40.262 CmdCreate ▶ DEBU 005 tar is available at  /bin/tar
2017/06/22 09:11:40 09:11:40.400 CmdCreate ▶ INFO 006 Found image:  sha256:90b   [echo-server:latest]
2017/06/22 09:11:40 09:11:40.400 CmdCreate ▶ INFO 007 Creating container:  tmp-echo_test
2017/06/22 09:11:40 09:11:40.484 CmdCreate ▶ INFO 008 Running container:  tmp-echo_test
2017/06/22 09:11:40 09:11:40.806 CmdCreate ▶ INFO 009 Logs for container:  tmp-echo_test
Listening on :3333
2017/06/22 09:11:40 09:11:40.845 CmdCreate ▶ INFO 00a Exporting container:  tmp-echo_test
2017/06/22 09:11:40 09:11:40.847 CmdCreate ▶ DEBU 00b Writing container:  tmp-echo_test  to:  /tmp/tmp-echo_test.tar
2017/06/22 09:11:41 09:11:41.177 CmdCreate ▶ DEBU 00c We are NOT using sudo to untar container:  /tmp/tmp-echo_test.tar
2017/06/22 09:11:41 09:11:41.187 CmdCreate ▶ DEBU 00d runC spec:  []
./plurality run echo_test /echo-server
```

You should now have an echo server running on port 3333.  This can be tested using netcat with the following:

```bash
$ echo "Hello rootless container!" | nc localhost 3333
Hi, I received your message! It was 26 bytes long and that's what it said: "Hello rootless container!" ! Honestly I have no clue about what to do with your messages, so Bye Bye!
```


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
