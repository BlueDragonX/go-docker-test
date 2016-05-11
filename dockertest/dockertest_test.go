package dockertest_test

import (
	dockertest "."
	"github.com/fsouza/go-dockerclient"
	"testing"
)

var versions = []string{
	"1.8",
	"1.9",
	"1.10",
}

func TestClient(t *testing.T) {
	d, err := dockertest.New()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	client, err := d.Client()
	if err != nil {
		t.Error(err)
	} else if err := client.Ping(); err != nil {
		t.Error(err)
	}
}

func TestURL(t *testing.T) {
	d, err := dockertest.New()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	client, err := docker.NewClient(d.URL())
	if err != nil {
		t.Error(err)
	} else if err := client.Ping(); err != nil {
		t.Error(err)
	}
}

func TestVersions(t *testing.T) {
	for _, version := range versions {
		func() {
			t.Logf("test version %s", version)
			d, err := dockertest.NewVersion(version)
			if err != nil {
				t.Error(err)
				return
			}
			defer d.Close()
			client, err := d.Client()
			if err != nil {
				t.Error(err)
			} else if err := client.Ping(); err != nil {
				t.Error(err)
			}
		}()
	}
}
