package juju

import (
	"fmt"

	"github.com/juju/juju/api"
	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/jujuclient/apiconnector"
)

type Clients struct {
	applicationClient *application.Client // applicationClient is limited to application API calls
	statusClient      *api.Client         // statusClient is used to gather status information
}

func NewClientsUsingConnectorConfig(connectorConfig api.SimpleConnectorConfig) (*Clients, error) {
	connector, err := api.NewSimpleConnector(connectorConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating Juju SimpleConnector: %v", err)
	}

	return createClients(connector)
}

func NewClientsUsingClientStore(controllerName string, modelUUID string) (*Clients, error) {
	connector, err := apiconnector.New(apiconnector.Config{
		ControllerName: controllerName,
		ModelUUID:      modelUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating Juju SimpleConnector: %v", err)
	}

	return createClients(connector)
}

func createClients(connector api.Connector) (*Clients, error) {
	conn, err := connector.Connect()
	if err != nil {
		return nil, fmt.Errorf("error connecting using Juju SimpleConnector: %v", err)
	}

	c := new(Clients)
	c.applicationClient = application.NewClient(conn)
	c.statusClient = conn.Client()
	return c, nil
}
