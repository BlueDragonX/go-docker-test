Go Docker Test Server
=====================
[![Godocs](https://img.shields.io/badge/go-docs-green.svg?style=flat)](https://godoc.org/github.com/BlueDragonX/go-docker-test/dockertest)

This library automates the creation of an empty Docker test environment. It
leverages the [Docker-in-Docker][1] containers published by Docker.

Usage
-----
You can launch an environment as follows (error handling omitted for brevity):

    d, _ := dockertest.New()
    defer d.Close()
    client, _ := d.Client()
    client.Ping()

The above will launch the latest Docker-in-Docker version available. To launch a particular version:

    d, _ := dockertest.NewVersion("1.8")
    defer d.Close()
    fmt.Printf("Daemon available at %s\n", d.URL())

[1]: https://hub.docker.com/_/docker/ "Docker-in-Docker"
