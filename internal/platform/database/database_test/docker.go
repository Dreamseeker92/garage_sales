package database_test

import (
	docker "github.com/fsouza/go-dockerclient"
	"sync"
	"testing"
)

var (
	once     sync.Once
	instance *docker.Client
)

// getClient implements Singleton pattern to share a single instance of a *docker.Client across the functions.
func getClient() *docker.Client {
	once.Do(func() {
		instance, _ = docker.NewClientFromEnv()
	})

	return instance
}

func startContainer(t *testing.T) *docker.Container {
	t.Helper()
	client := getClient()

	port := docker.Port("5432/tcp")
	bindings := []docker.PortBinding{{HostIP: "0.0.0.0", HostPort: "5432"}}

	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Name: "postgres_test",
		Config: &docker.Config{
			Image: "postgres:13-alpine",
			Env:   []string{"POSTGRES_USER=garage", "POSTGRES_PASSWORD=garage", "POSTGRES_DB=garage_test"},
		},
		HostConfig: &docker.HostConfig{
			PortBindings: map[docker.Port][]docker.PortBinding{port: bindings},
		},
	})

	if err != nil {
		t.Fatalf("Creating db container: %v", err)
	}

	if err := client.StartContainer(container.ID, &docker.HostConfig{AutoRemove: true}); err != nil {
		stopContainer(t, container)
	}

	return inspectContainer(t, container)
}

// inspectContainer is a tool to reinitialize a container instance with all parameters set.
func inspectContainer(t *testing.T, container *docker.Container) *docker.Container {
	client := getClient()

	cont, err := client.InspectContainerWithOptions(docker.InspectContainerOptions{ID: container.ID})
	if err != nil {
		t.Fatalf("Error reinitializing container with ID %s: %v", container.ID, err)
	}

	return cont
}

// stopContainer stops and removes a container.
func stopContainer(t *testing.T, container *docker.Container) {
	t.Helper()

	client := getClient()
	removeConfig := docker.RemoveContainerOptions{ID: container.ID, Force: true}
	if err := client.RemoveContainer(removeConfig); err != nil {
		t.Fatalf("Stopping container: %v", err)
	}
}
