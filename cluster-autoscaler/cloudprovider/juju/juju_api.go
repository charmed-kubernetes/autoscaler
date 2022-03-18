package juju

import (
	"fmt"

	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/rpc/params"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
)

type Unit struct {
	state    cloudprovider.InstanceState
	jujuName string
	kubeName string
	workload string
	agent    string
}

type Manager struct {
	clients     *Clients
	controller  string
	model       string
	application string
	units       map[string]*Unit
}

func (m *Manager) init() error {
	var err error
	m.clients, err = NewClientsUsingClientStore(m.controller, m.model)
	if err != nil {
		return err
	}

	fullStatus, err := m.getStatus()
	if err != nil {
		return fmt.Errorf("error getting status: %v", err)
	}

	app := fullStatus.Applications[m.application]
	for unitName := range app.Units {
		hostname, err := m.getHostnameForUnitNamed(unitName)
		if err != nil {
			return fmt.Errorf("error getting hostname for unit %v: %v", unitName, err)
		}
		m.units[unitName] = &Unit{
			state:    cloudprovider.InstanceRunning,
			jujuName: unitName,
			kubeName: hostname,
		}
	}

	return nil
}

func (m *Manager) addUnits(delta int) error {
	prevStatus, err := m.getStatus()
	if err != nil {
		return err
	}

	_, err = m.clients.applicationClient.AddUnits(application.AddUnitsParams{
		ApplicationName: m.application,
		NumUnits:        delta,
	})
	if err != nil {
		return err
	}

	jujuStatus, err := m.getStatus()
	if err != nil {
		return err
	}

	for unitName, _ := range jujuStatus.Applications[m.application].Units {
		if _, ok := prevStatus.Applications[m.application].Units[unitName]; !ok {
			m.units[unitName] = &Unit{
				state:    cloudprovider.InstanceCreating,
				jujuName: unitName,
			}
		}
	}

	return nil
}

func (m *Manager) removeUnit(name string) error {
	unit := m.getUnitByHostname(name)
	unit.state = cloudprovider.InstanceDeleting

	units := []string{unit.jujuName}
	args := application.DestroyUnitsParams{
		Units:          units,
		DestroyStorage: false,
		Force:          false,
	}

	_, err := m.clients.applicationClient.DestroyUnits(args)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) refresh() error {
	fullStatus, err := m.getStatus()
	if err != nil {
		return err
	}

	for unitName, unitStatus := range fullStatus.Applications[m.application].Units {
		if _, ok := m.units[unitName]; ok {
			m.units[unitName].agent = unitStatus.AgentStatus.Status
			m.units[unitName].workload = unitStatus.WorkloadStatus.Status
		}
	}

	for unitName, unit := range m.units {
		if unit.state == cloudprovider.InstanceCreating {
			if unit.kubeName == "" {
				hostname, err := m.getHostnameForUnitNamed(unitName)
				if err != nil {
					return fmt.Errorf("error getting hostname for unit %v: %v", unit, err)
				}
				unit.kubeName = hostname
			}

			if unit.workload == "active" {
				unit.state = cloudprovider.InstanceRunning
			}
		} else if unit.state == cloudprovider.InstanceDeleting {
			delete(m.units, unitName)
		}
	}

	return nil
}

func (m *Manager) getUnitByHostname(hostname string) *Unit {
	for _, unit := range m.units {
		if unit.kubeName == hostname {
			return unit
		}
	}
	return nil
}

func (m *Manager) getHostnameForUnitNamed(unitName string) (string, error) {
	fullStatus, err := m.clients.statusClient.Status(nil)
	if err != nil {
		return "", err
	} else {
		app := fullStatus.Applications[m.application]
		unitStatus := app.Units[unitName]
		return fullStatus.Machines[unitStatus.Machine].Hostname, nil
	}
}

func (m *Manager) getStatus() (*params.FullStatus, error) {
	status, err := m.clients.statusClient.Status(nil)
	if err != nil {
		return nil, fmt.Errorf("error getting status: %v", err)
	}

	return status, nil
}
