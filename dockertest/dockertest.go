package dockertest

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"time"
)

// The default version of the Docker-in-Docker container to run.
const DEFAULT_VERSION = "1.10"

// Daemon runs a Docker-in-Docker container fo testing.
type Docker struct {
	id     string
	port   int
	client *docker.Client
}

// NewVersion creates a new, running Docker-in-Docker container running the
// specified version of Docker.
func NewVersion(version string) (*Docker, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create test docker client")
	}
	daemon := &Docker{
		client: client,
	}

	if version == "" {
		version = DEFAULT_VERSION
	}
	tag := fmt.Sprintf("%s-dind", version)
	return daemon, daemon.start(tag)
}

// New runs the default verison of the Docker-in-Docker container.
func New() (*Docker, error) {
	return NewVersion("")
}

// URL returns the URL where the daemon is reachable by clients.
func (d *Docker) URL() string {
	return fmt.Sprintf("tcp://127.0.0.1:%d/", d.port)
}

// Client returns a client connected to the daemon.
func (d *Docker) Client() (*docker.Client, error) {
	return docker.NewClient(d.URL())
}

// Close stops the daemon and removes it from the parent Docker.
func (d *Docker) Close() error {
	return d.remove(true)
}

// pull the Docker-in-Docker image.
func (d *Docker) pull(tag string) error {
	opts := docker.PullImageOptions{
		Repository: "docker",
		Tag: tag,
		OutputStream: os.Stderr,
	}
	auth := docker.AuthConfiguration{}
	return d.client.PullImage(opts, auth)
}

// start the Docker-in-Docker container.
func (d *Docker) start(tag string) error {
	if d.id != "" {
		return errors.Errorf("daemon already running at %s", d.id)
	}
	if err := d.pull(tag); err != nil {
		return errors.Wrapf(err, "unable to pull image docker:%s", tag)
	}
	image := fmt.Sprintf("docker:%s", tag)
	fmt.Printf("image: %s\n", image)
	container, err := d.client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: image,
		},
	})
	if err != nil {
		return errors.Wrapf(err, "unable to create %s", image)
	}
	d.id = container.ID

	err = d.client.StartContainer(d.id, &docker.HostConfig{
		Privileged: true,
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port("2375/tcp"): {
				{
					HostIP:   "127.0.0.1",
					HostPort: "0",
				},
			},
		},
	})
	if err != nil {
		d.remove(false)
		return errors.Wrapf(err, "unable to start %s", image)
	}

	container, err = d.client.InspectContainer(d.id)
	if err != nil {
		d.remove(false)
		return errors.Wrapf(err, "unable to inspect %s", d.id)
	}

	portStr := container.NetworkSettings.Ports["2375/tcp"][0].HostPort
	port, err := strconv.Atoi(portStr)
	if err != nil {
		d.remove(false)
		return errors.Wrapf(err, "unable to convert port %s to integer", portStr)
	}
	d.port = port
	return d.wait()
}

// remove the container.
func (d *Docker) remove(force bool) error {
	return d.client.RemoveContainer(docker.RemoveContainerOptions{
		ID:    d.id,
		Force: force,
	})
}

// wait for the daemon to become available.
func (d *Docker) wait() error {
	client, err := docker.NewClient(d.URL())
	if err != nil {
		return err
	}
	for i := 0; i < 4; i++ {
		err = client.Ping()
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	return err
}
