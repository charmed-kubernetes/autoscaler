package juju

//go:generate mockgen -destination=./mocks/mock_juju_client.go --build_flags=--mod=mod -package=mocks . JujuClient
//go:generate mockgen -destination=./mocks/mock_kube_client.go --build_flags=--mod=mod -package=mocks k8s.io/client-go/kubernetes Interface
//go:generate mockgen -destination=./mocks/mock_core_v1.go --build_flags=--mod=mod -package=mocks k8s.io/client-go/kubernetes/typed/core/v1 CoreV1Interface,NodeInterface

import (
	ctx "context"
	"fmt"

	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/rpc/params"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	kube_client "k8s.io/client-go/kubernetes"
	klog "k8s.io/klog/v2"
)

type Unit struct {
	state      cloudprovider.InstanceState
	jujuName   string
	hostname   string
	status     params.UnitStatus
	providerID string
}

type JujuClient interface {
	AddUnits(args application.AddUnitsParams) ([]string, error)
	DestroyUnits(args application.DestroyUnitsParams) ([]params.DestroyUnitResult, error)
	Status(patterns []string) (*params.FullStatus, error)
}

type Manager struct {
	jujuClient  JujuClient
	kubeClient  kube_client.Interface
	model       string
	application string
	units       map[string]*Unit
}

func NewManager(jujuClient JujuClient, kubeClient kube_client.Interface, model string, application string) (*Manager, error) {
	klog.Infof("creating manager")
	m := new(Manager)
	m.jujuClient = jujuClient
	m.kubeClient = kubeClient
	m.model = model
	m.application = application
	m.units = make(map[string]*Unit)

	fullStatus, err := m.jujuClient.Status(nil)
	if err != nil {
		klog.Error("error getting status from juju client", err.Error())
		return nil, err
	}

	app := fullStatus.Applications[m.application]
	for unitName, unitStatus := range app.Units {
		unitState := cloudprovider.InstanceCreating
		if unitStatus.WorkloadStatus.Status == "active" && unitStatus.AgentStatus.Status == "idle" {
			unitState = cloudprovider.InstanceRunning
		}
		m.units[unitName] = &Unit{
			state:    unitState,
			jujuName: unitName,
			hostname: fullStatus.Machines[unitStatus.Machine].Hostname,
			status:   unitStatus,
		}
	}

	return m, nil
}

func (m *Manager) addUnits(delta int) error {
	prevStatus, err := m.jujuClient.Status(nil)
	if err != nil {
		return err
	}

	_, err = m.jujuClient.AddUnits(application.AddUnitsParams{
		ApplicationName: m.application,
		NumUnits:        delta,
	})
	if err != nil {
		return err
	}

	currentStatus, err := m.jujuClient.Status(nil)
	if err != nil {
		return err
	}

	for unitName, unitStatus := range currentStatus.Applications[m.application].Units {
		if _, ok := prevStatus.Applications[m.application].Units[unitName]; !ok {
			m.units[unitName] = &Unit{
				state:    cloudprovider.InstanceCreating,
				jujuName: unitName,
				status:   unitStatus,
			}
			klog.Infof("added unit %s to managed units", unitName)
		}
	}

	return nil
}

func (m *Manager) removeUnit(hostname string) error {
	unit := m.getUnitByHostname(hostname)
	if unit == nil {
		return fmt.Errorf("unit with hostname %s not found", hostname)
	}
	unit.state = cloudprovider.InstanceDeleting

	units := []string{unit.jujuName}
	args := application.DestroyUnitsParams{
		Units:          units,
		DestroyStorage: false,
		Force:          false,
	}

	_, err := m.jujuClient.DestroyUnits(args)
	if err != nil {
		return err
	}

	klog.Infof("unit %s state changed to InstanceDeleting", unit.jujuName)
	return nil
}

func (m *Manager) refresh() error {
	fullStatus, err := m.jujuClient.Status(nil)
	if err != nil {
		return err
	}

	// Loop over the units in the status and update the manager to match
	// This could mean updating the state of units currently managed by the manager
	// or incorporating a totally new unit that was added by the cluster-admin manually
	for unitName, unitStatus := range fullStatus.Applications[m.application].Units {

		unitHostname := fullStatus.Machines[unitStatus.Machine].Hostname
		unitProviderID := ""
		if unitHostname != "" {
			node, err := m.kubeClient.CoreV1().Nodes().Get(ctx.TODO(), unitHostname, v1.GetOptions{})
			if err != nil {
				klog.Errorf("error getting provider ID for unit %v with hostname %v: %v", unitName, unitHostname, err.Error())
			} else {
				unitProviderID = node.Spec.ProviderID
			}
		}

		// Check if we aren't already managing this unit
		if _, ok := m.units[unitName]; !ok {
			// Check if the unit is active and idle
			// This is necessary since when a unit gets deleted it does not happen immediately
			// We want to make sure we only add externally added units, not recently deleted units that are still showing up in status
			if unitStatus.WorkloadStatus.Status == "active" && unitStatus.AgentStatus.Status == "idle" {
				// The unit was added manually. Need to add it to the units list as a new unit
				m.units[unitName] = &Unit{
					state:      cloudprovider.InstanceRunning,
					jujuName:   unitName,
					hostname:   unitHostname,
					status:     unitStatus,
					providerID: unitProviderID,
				}
				klog.Infof("detected unmanaged unit %s", unitName)
				klog.Infof("added unit %s to managed units", unitName)
			}
		} else {
			// Update the status, hostname, and providerID of each unit
			m.units[unitName].status = unitStatus
			m.units[unitName].hostname = unitHostname
			m.units[unitName].providerID = unitProviderID
		}
	}

	// Based on the state, decide if we need to delete any units, or update any freshly created units to running
	for unitName, unit := range m.units {
		// Check if any unit we are managing does not exist in the list of units we got from status
		if _, ok := fullStatus.Applications[m.application].Units[unitName]; !ok {
			// A unit we were managing does not exist in the list of units we got from Juju status.
			// Change the state to InstanceDeleting so it gets removed below
			unit.state = cloudprovider.InstanceDeleting
			klog.Infof("detected managed unit %s that has been removed", unit.jujuName)
		}

		if unit.state == cloudprovider.InstanceCreating {
			if unit.status.WorkloadStatus.Status == "active" && unit.status.AgentStatus.Status == "idle" {
				unit.state = cloudprovider.InstanceRunning
			}
		} else if unit.state == cloudprovider.InstanceDeleting {
			delete(m.units, unitName)
			klog.Infof("removed unit %s from managed units", unit.jujuName)
		}
	}

	return nil
}

func (m *Manager) getUnitByHostname(hostname string) *Unit {
	for _, unit := range m.units {
		if unit.hostname == hostname {
			return unit
		}
	}
	return nil
}
