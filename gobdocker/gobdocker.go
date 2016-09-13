// import github.com/dmcsorley/goblin/gobdocker
package gobdocker

import (
	"bufio"
	"encoding/json"
	"github.com/dmcsorley/goblin/goblog"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"log"
)

func ListenForBuildExits() {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient(client.DefaultDockerHost, client.DefaultVersion, nil, defaultHeaders)
	if err != nil {
		panic(err)
	}

	f := filters.NewArgs()
	f.Add("label", "goblin.build")
	f.Add("event", "die")

	eventstream, err := cli.Events(
		context.Background(),
		types.EventsOptions{
			Filters: f,
		},
	)
	if err != nil {
		panic(err)
	}

	defer eventstream.Close()

	s := bufio.NewScanner(eventstream)
	for s.Scan() {
		event := events.Message{}
		err := json.Unmarshal([]byte(s.Text()), &event)
		if err != nil {
			log.Printf("%v\n", err)
		} else {
			goblog.Log(
				event.Actor.Attributes["goblin.id"],
				"EXITED " + event.Actor.Attributes["exitCode"],
			)
		}
	}
}
