package gobdocker

import (
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func RemoveImage(name string) error {
	c, _ := getClient()

	_, err := c.ImageRemove(
		context.Background(),
		name,
		types.ImageRemoveOptions{PruneChildren: true},
	)

	return err
}
