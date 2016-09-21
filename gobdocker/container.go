package gobdocker

import (
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
	"log"
)

func RemoveContainer(name string) {
	c, _ := getClient()
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
	c, _ := getClient()
	return c.ContainerWait(context.Background(), name)
}
