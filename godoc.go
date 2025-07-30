package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
)

func main() {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	fmt.Println("Pulling nginx image...")
	reader, err := cli.ImagePull(ctx, "nginx:latest", image.PullOptions{})
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)

	// Absolute path to nginx.conf
	nginxConfPath, err := filepath.Abs("nginx.conf")
	if err != nil {
		panic(err)
	}

	// Absolute path to index.html
	indexConfPath, err := filepath.Abs("index.html")
	if err != nil {
		panic(err)
	}

	// Define port binding
	portSet := nat.PortSet{
		"80/tcp": struct{}{},
	}
	portMap := nat.PortMap{
		"80/tcp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "8002",
			},
		},
	}

	// Create container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        "nginx:latest",
		ExposedPorts: portSet,
	}, &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/etc/nginx/nginx.conf", nginxConfPath),
			fmt.Sprintf("%s:/var/www/html/index.html", indexConfPath),
		},
		PortBindings: portMap,
	}, &network.NetworkingConfig{}, nil, "go-nginx")
	if err != nil {
		panic(err)
	}

	// Start container
	fmt.Println("Starting container...")
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println("NGINX running at http://localhost:8002")

	// Show logs
	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		panic(err)
	}
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
