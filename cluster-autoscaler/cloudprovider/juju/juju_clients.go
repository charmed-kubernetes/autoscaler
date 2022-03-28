package juju

import (
	"fmt"

	"github.com/juju/juju/api"
	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/api/connector"
)

type Clients struct {
	applicationClient *application.Client // applicationClient is limited to application API calls
	statusClient      *api.Client         // statusClient is used to gather status information
}

func NewClients(connectorConfig connector.SimpleConfig) (*Clients, error) {
	connector, err := connector.NewSimple(connectorConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating Juju SimpleConnector: %v", err)
	}

	conn, err := connector.Connect()
	if err != nil {
		return nil, fmt.Errorf("error connecting using Juju SimpleConnector: %v", err)
	}

	c := new(Clients)
	c.applicationClient = application.NewClient(conn)
	c.statusClient = api.NewClient(conn)
	return c, nil
}
