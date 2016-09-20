package gobdocker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func RemoveImage(name string) error {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	c, _ := client.NewClient(client.DefaultDockerHost, client.DefaultVersion, nil, defaultHeaders)

	_, err := c.ImageRemove(
		context.Background(),
		name,
		types.ImageRemoveOptions{PruneChildren: true},
	)

	return err
}
