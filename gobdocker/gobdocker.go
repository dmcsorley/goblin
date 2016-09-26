// import github.com/dmcsorley/goblin/gobdocker
package gobdocker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"log"
)

type ExitedBuild struct {
	Id          string
	ContainerId string
	Name        string
	Time        string
	Exit        string
}

func getClient() (*client.Client, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	return client.NewClient(client.DefaultDockerHost, client.DefaultVersion, nil, defaultHeaders)
}

func ListenForBuildExits(callback func(*ExitedBuild)) {
	cli, err := getClient()
	if err != nil {
		panic(err)
	}

	f := filters.NewArgs()
	f.Add("label", "goblin.build")
	f.Add("event", "die")

	eventchan, errchan := cli.Events(
		context.Background(),
		types.EventsOptions{
			Filters: f,
		},
	)

	for {
	select {
	case event := <- eventchan:
		callback(&ExitedBuild{
			Id:          event.Actor.Attributes["goblin.id"],
			ContainerId: event.Actor.ID,
			Name:        event.Actor.Attributes["goblin.name"],
			Time:        event.Actor.Attributes["goblin.time"],
			Exit:        event.Actor.Attributes["exitCode"],
		})
	case err := <- errchan:
		log.Printf("%v\n", err)
		return
	}
	}
}

func CreateVolume(name string) (string, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, _ := client.NewClient(client.DefaultDockerHost, client.DefaultVersion, nil, defaultHeaders)

	volume, err := cli.VolumeCreate(
		context.Background(),
		types.VolumeCreateRequest{
			Name: name,
		},
	)

	if err != nil {
		return "", err
	}

	return volume.Name, nil
}

func RemoveVolume(name string) {
	cli, _ := getClient()
	err := cli.VolumeRemove(context.Background(), name, false)
	if err != nil {
		log.Println(err)
	}
}
