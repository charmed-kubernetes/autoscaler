package juju

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/rpc/params"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/juju/mocks"
)

func makeUnit(state cloudprovider.InstanceState, jujuName string, kubeName string, agentStatus string, workloadStatus string, machine string) Unit {
	return Unit{
		state:    state,
		jujuName: jujuName,
		kubeName: kubeName,
		status: params.UnitStatus{
			AgentStatus: params.DetailedStatus{
				Status: agentStatus,
			},
			WorkloadStatus: params.DetailedStatus{
				Status: workloadStatus,
			},
			Machine: machine,
		},
	}
}

func makeApplicationStatus(appName string, units map[string]*Unit) map[string]params.ApplicationStatus {
	unitStatuses := make(map[string]params.UnitStatus)
	for _, unit := range units {
		unitStatuses[unit.jujuName] = unit.status
	}

	return map[string]params.ApplicationStatus{
		appName: {
			Units: unitStatuses,
		},
	}
}

func makeMachineStatuses(units map[string]*Unit) map[string]params.MachineStatus {
	machineStatuses := make(map[string]params.MachineStatus)
	for _, unit := range units {
		machineStatuses[unit.status.Machine] = params.MachineStatus{
			Hostname: unit.kubeName,
		}
	}

	return machineStatuses
}

func makeStatus(appName string, units map[string]*Unit) params.FullStatus {
	return params.FullStatus{
		Applications: makeApplicationStatus(appName, units),
		Machines:     makeMachineStatuses(units),
	}
}

func makeManager(mockJujuClient *mocks.MockJujuClient, units []Unit) *Manager {
	m := NewManager(mockJujuClient, "test_model", "test_application")
	for _, v := range units {
		unit := v
		m.units[v.jujuName] = &unit
	}
	return m
}

func TestNewManager(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockJujuClient := mocks.NewMockJujuClient(ctl)

	m := NewManager(mockJujuClient, "test_model", "test_application")

	if m.jujuClient != mockJujuClient {
		t.Errorf("m.jujuClient = %v; want %v", m.jujuClient, mockJujuClient)
	}

	if m.model != "test_model" {
		t.Errorf("m.model = %v; want %v", m.model, "test_model")
	}

	if m.application != "test_application" {
		t.Errorf("m.application = %v; want %v", m.application, "test_application")
	}
}

func TestInit(t *testing.T) {

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockJujuClient := mocks.NewMockJujuClient(ctl)
	m := makeManager(mockJujuClient, []Unit{})

	unit1 := makeUnit(cloudprovider.InstanceRunning, "unit_1", "unit_1_hostname", "idle", "active", "machine_1")
	unit2 := makeUnit(cloudprovider.InstanceRunning, "unit_2", "unit_2_hostname", "error", "blocked", "machine_2")
	unit3 := makeUnit(cloudprovider.InstanceRunning, "unit_3", "unit_3_hostname", "idle", "active", "machine_3")
	unit4 := makeUnit(cloudprovider.InstanceDeleting, "unit_4", "unit_4_hostname", "idle", "active", "machine_4")
	units := map[string]*Unit{
		unit1.jujuName: &unit1,
		unit2.jujuName: &unit2,
		unit3.jujuName: &unit3,
		unit4.jujuName: &unit4,
	}
	ms := makeStatus("test_application", units)
	gomock.InOrder(
		mockJujuClient.EXPECT().Status(nil).Return(&ms, nil), // Getting initial status
		mockJujuClient.EXPECT().Status(nil).Return(&ms, nil), // Getting hostname of unit_1
		mockJujuClient.EXPECT().Status(nil).Return(&ms, nil), // Getting hostname of unit_2
		mockJujuClient.EXPECT().Status(nil).Return(&ms, nil), // Getting hostname of unit_3
		mockJujuClient.EXPECT().Status(nil).Return(&ms, nil), // Getting hostname of unit_4
	)

	err := m.init()
	if err != nil {
		t.Errorf("unexpected error returned from init")
	}

	if len(m.units) != 4 {
		t.Errorf("len(m.units) = %v; want %v", len(m.units), 4)
	}
}

func TestAddUnit(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockAddedUnit := makeUnit(cloudprovider.InstanceRunning, "added_unit", "added_unit_hostname", "idle", "active", "added_machine")
	mockJujuClient := mocks.NewMockJujuClient(ctl)

	unit1 := makeUnit(cloudprovider.InstanceRunning, "unit_1", "unit_1_hostname", "idle", "active", "machine_1")
	units := []Unit{unit1}
	m := makeManager(mockJujuClient, units)

	mockStatusBeforeAdd := makeStatus(m.application, m.units)

	// We also need to mock a status that reflects the old status, plus the addition of the newly added machine
	unitsPlusOne := map[string]*Unit{
		unit1.jujuName:         &unit1,
		mockAddedUnit.jujuName: &mockAddedUnit,
	}
	mockStatusAfterAdd := makeStatus(m.application, unitsPlusOne)

	gomock.InOrder(
		mockJujuClient.EXPECT().Status(nil).Return(&mockStatusBeforeAdd, nil),              // Getting previous status
		mockJujuClient.EXPECT().AddUnits(gomock.Any()).Return([]string{"added_unit"}, nil), // Add a unit
		mockJujuClient.EXPECT().Status(nil).Return(&mockStatusAfterAdd, nil),               // Getting status after adding a unit
	)

	if _, ok := m.units[mockAddedUnit.jujuName]; ok {
		t.Errorf("units contains added unit before addUnits was callled")
	}

	m.addUnits(1)

	if _, ok := m.units[mockAddedUnit.jujuName]; !ok {
		t.Errorf("units does not contain added unit")
	}

	// state should now be creating
	if m.units[mockAddedUnit.jujuName].state != cloudprovider.InstanceCreating {
		t.Errorf("state = %v; want %v", m.units[mockAddedUnit.jujuName].state, cloudprovider.InstanceCreating)
	}

}

func TestRemoveUnit(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockJujuClient := mocks.NewMockJujuClient(ctl)
	unit1 := makeUnit(cloudprovider.InstanceRunning, "unit_1", "unit_1_hostname", "idle", "active", "machine_1")
	units := []Unit{unit1}
	m := makeManager(mockJujuClient, units)

	args := application.DestroyUnitsParams{
		Units:          []string{"unit_1"},
		DestroyStorage: false,
		Force:          false,
	}

	mockJujuClient.EXPECT().DestroyUnits(args).Return(nil, nil)

	err := m.removeUnit("unit_1_hostname")
	if err != nil {
		t.Errorf("error removing unit: %s", err.Error())
	}

	if m.getUnitByHostname("unit_1_hostname").state != cloudprovider.InstanceDeleting {
		t.Errorf("state = %v; want %v", m.getUnitByHostname("unit_1_hostname").state, cloudprovider.InstanceDeleting)
	}

	// Test case when hostname is not found
	err = m.removeUnit("the_host_does_not_exist")
	if err == nil {
		t.Errorf("expected error but did not get one")
	}

	// Test case when DestroyUnits returns an error
	args = application.DestroyUnitsParams{
		Units:          []string{"unit_1"},
		DestroyStorage: false,
		Force:          false,
	}
	mockJujuClient.EXPECT().DestroyUnits(args).Return(nil, errors.New("some error"))
	err = m.removeUnit("unit_1_hostname")
	if err.Error() != "some error" {
		t.Errorf("expected some error but did not get it")
	}
}

func TestRefresh(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockJujuClient := mocks.NewMockJujuClient(ctl)
	unit1 := makeUnit(cloudprovider.InstanceRunning, "unit_1", "unit_1_hostname", "idle", "active", "machine_1")
	unit2 := makeUnit(cloudprovider.InstanceRunning, "unit_2", "unit_2_hostname", "error", "blocked", "machine_2")
	unit3 := makeUnit(cloudprovider.InstanceCreating, "unit_3", "unit_3_hostname", "idle", "active", "machine_3")
	unit4 := makeUnit(cloudprovider.InstanceDeleting, "unit_4", "unit_4_hostname", "idle", "active", "machine_4")
	units := []Unit{unit1, unit2, unit3, unit4}
	m := makeManager(mockJujuClient, units)

	ms := makeStatus(m.application, m.units)
	gomock.InOrder(
		mockJujuClient.EXPECT().Status(nil).Return(&ms, nil), // Getting previous status
	)

	// mockUnit3 state should be creating before the call
	if m.units[unit3.jujuName].state != cloudprovider.InstanceCreating {
		t.Errorf("state = %v; want %v", m.units[unit3.jujuName].state, cloudprovider.InstanceCreating)
	}

	err := m.refresh()
	if err != nil {
		t.Errorf("error refreshing: %s", err.Error())
	}

	// mockUnit3 state should now be running since it was previously creating (and active)
	if m.units[unit3.jujuName].state != cloudprovider.InstanceRunning {
		t.Errorf("state = %v; want %v", m.units[unit3.jujuName].state, cloudprovider.InstanceRunning)
	}

	// mockUnit4 should be deleted
	if _, ok := m.units[unit4.jujuName]; ok {
		t.Errorf("units contain unit that should have been removed")
	}

}
