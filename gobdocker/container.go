package gobdocker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"log"
)

func RemoveContainer(name string) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	c, _ := client.NewClient(client.DefaultDockerHost, client.DefaultVersion, nil, defaultHeaders)
	err := c.ContainerRemove(
		context.Background(),
		name,
		types.ContainerRemoveOptions{
			Force:         true,
			RemoveVolumes: true,
		},
	)
	if err != nil {
		log.Println(err)
	}
}

func WaitContainer(name string) (int, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	c, _ := client.NewClient(client.DefaultDockerHost, client.DefaultVersion, nil, defaultHeaders)
	return c.ContainerWait(context.Background(), name)
}
