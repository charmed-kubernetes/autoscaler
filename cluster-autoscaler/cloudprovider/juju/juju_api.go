package juju

import (
	"fmt"

	"github.com/juju/juju/api"
	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/api/connector"
	"github.com/juju/juju/rpc/params"
)

type JujuAPI struct {
	applicationClient *application.Client // applicationClient is limited to application API calls
	statusClient      *api.Client         // statusClient is used to gather status information
}

func NewJujuAPi(connector *connector.SimpleConnector) (*JujuAPI, error) {
	conn, err := connector.Connect()
	if err != nil {
		return nil, fmt.Errorf("error connecting using Juju SimpleConnector: %v", err)
	}

	jujuAPI := new(JujuAPI)
	jujuAPI.applicationClient = application.NewClient(conn)
	jujuAPI.statusClient = api.NewClient(conn)
	return jujuAPI, nil
}

func (jujuAPI *JujuAPI) AddUnits(args application.AddUnitsParams) ([]string, error) {
	return jujuAPI.applicationClient.AddUnits(args)
}

func (jujuAPI *JujuAPI) DestroyUnits(args application.DestroyUnitsParams) ([]params.DestroyUnitResult, error) {
	return jujuAPI.applicationClient.DestroyUnits(args)
}

func (jujuAPI *JujuAPI) Status(patterns []string) (*params.FullStatus, error) {
	return jujuAPI.statusClient.Status(patterns)
}