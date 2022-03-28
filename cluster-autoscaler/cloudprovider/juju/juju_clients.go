package juju

import (
	"fmt"

	"github.com/juju/juju/api"
	"github.com/juju/juju/api/client/application"
<<<<<<< HEAD
	"github.com/juju/juju/api/connector"
=======
	"github.com/juju/juju/jujuclient/apiconnector"
>>>>>>> 1cb7c9a8c04b7de79c2dd46f84bd5239eed4ee16
)

type Clients struct {
	applicationClient *application.Client // applicationClient is limited to application API calls
	statusClient      *api.Client         // statusClient is used to gather status information
}

<<<<<<< HEAD
func NewClients(connectorConfig connector.SimpleConfig) (*Clients, error) {
	connector, err := connector.NewSimple(connectorConfig)
=======
func NewClientsUsingConnectorConfig(connectorConfig api.SimpleConnectorConfig) (*Clients, error) {
	connector, err := api.NewSimpleConnector(connectorConfig)
>>>>>>> 1cb7c9a8c04b7de79c2dd46f84bd5239eed4ee16
	if err != nil {
		return nil, fmt.Errorf("error creating Juju SimpleConnector: %v", err)
	}

<<<<<<< HEAD
=======
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
>>>>>>> 1cb7c9a8c04b7de79c2dd46f84bd5239eed4ee16
	conn, err := connector.Connect()
	if err != nil {
		return nil, fmt.Errorf("error connecting using Juju SimpleConnector: %v", err)
	}

	c := new(Clients)
	c.applicationClient = application.NewClient(conn)
<<<<<<< HEAD
	c.statusClient = api.NewClient(conn)
=======
	c.statusClient = conn.Client()
>>>>>>> 1cb7c9a8c04b7de79c2dd46f84bd5239eed4ee16
	return c, nil
}
