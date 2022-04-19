package juju

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/juju/juju/api/client/application"
	"github.com/juju/juju/rpc/params"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/juju/mocks"
)

// Note: If you need to generate new mocks , run go generate ./... in the cloudprovider/juju directory
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

func makeManager(mockJujuClient *mocks.MockJujuClient, units map[string]*Unit) (*Manager, error) {
	ms := makeStatus("test_application", units)
	mockJujuClient.EXPECT().Status(nil).Return(&ms, nil)
	m, err := NewManager(mockJujuClient, "test_model", "test_application")

	// The constructor initializes the units it gets from status in the instance running state, lets fix those so that the managers units
	// equal the units passed in as an argument
	m.units = units
	if err != nil {
		return nil, err
	}

	return m, nil
}

func TestNewManager(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockJujuClient := mocks.NewMockJujuClient(ctl)
	unit1 := makeUnit(cloudprovider.InstanceRunning, "unit_1", "unit_1_hostname", "idle", "active", "machine_1")
	unit2 := makeUnit(cloudprovider.InstanceCreating, "unit_2", "unit_2_hostname", "error", "blocked", "machine_2")
	units := map[string]*Unit{
		unit1.jujuName: &unit1,
		unit2.jujuName: &unit2,
	}

	ms := makeStatus("test_application", units)
	mockJujuClient.EXPECT().Status(nil).Return(&ms, nil)

	m, err := NewManager(mockJujuClient, "test_model", "test_application")
	if err != nil {
		t.Fatalf("error creating manager")
	}

	if m.jujuClient != mockJujuClient {
		t.Errorf("m.jujuClient = %v; want %v", m.jujuClient, mockJujuClient)
	}

	if m.model != "test_model" {
		t.Errorf("m.model = %v; want %v", m.model, "test_model")
	}

	if m.application != "test_application" {
		t.Errorf("m.application = %v; want %v", m.application, "test_application")
	}

	for unitName, unit := range m.units {
		if unit.state != units[unitName].state {
			t.Errorf("%v state = %v; want %v", unitName, unit.state, units[unitName].state)
		}

		if unit.jujuName != units[unitName].jujuName {
			t.Errorf("%v jujuName = %v; want %v", unitName, unit.jujuName, units[unitName].jujuName)
		}

		if unit.kubeName != units[unitName].kubeName {
			t.Errorf("%v kubeName = %v; want %v", unitName, unit.kubeName, units[unitName].kubeName)
		}

		if !cmp.Equal(unit.status, units[unitName].status) {
			t.Errorf("%v status = %v; want %v", unitName, unit.status, units[unitName].status)
		}
	}

	// Test the error path for status()
	mockJujuClient.EXPECT().Status(nil).Return(nil, errors.New("status error"))
	_, err = NewManager(mockJujuClient, "test_model", "test_application")
	if err.Error() != "status error" {
		t.Errorf("expected status error but did not get it")
	}

}

func TestAddUnit(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockAddedUnit := makeUnit(cloudprovider.InstanceRunning, "added_unit", "added_unit_hostname", "idle", "active", "added_machine")
	mockJujuClient := mocks.NewMockJujuClient(ctl)

	unit1 := makeUnit(cloudprovider.InstanceRunning, "unit_1", "unit_1_hostname", "idle", "active", "machine_1")
	units := map[string]*Unit{
		unit1.jujuName: &unit1,
	}
	m, err := makeManager(mockJujuClient, units)
	if err != nil {
		t.Fatalf("error creating manager")
	}

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

	// State should now be InstanceCreating for the added unit
	if m.units[mockAddedUnit.jujuName].state != cloudprovider.InstanceCreating {
		t.Errorf("state = %v; want %v", m.units[mockAddedUnit.jujuName].state, cloudprovider.InstanceCreating)
	}

	// Test error path when getting previous status
	mockJujuClient.EXPECT().Status(nil).Return(nil, errors.New("previous status error"))
	err = m.addUnits(1)
	if err.Error() != "previous status error" {
		t.Errorf("expected previous status error but did not get it")
	}

	// Test error path when calling AddUnits
	gomock.InOrder(
		mockJujuClient.EXPECT().Status(nil).Return(&mockStatusBeforeAdd, nil),                    // Getting previous status
		mockJujuClient.EXPECT().AddUnits(gomock.Any()).Return(nil, errors.New("AddUnits error")), // Add a unit
	)
	err = m.addUnits(1)
	if err.Error() != "AddUnits error" {
		t.Errorf("expected AddUnits error but did not get it")
	}

	// Test error path when getting current status
	gomock.InOrder(
		mockJujuClient.EXPECT().Status(nil).Return(&mockStatusBeforeAdd, nil),               // Getting previous status
		mockJujuClient.EXPECT().AddUnits(gomock.Any()).Return([]string{"added_unit"}, nil),  // Add a unit
		mockJujuClient.EXPECT().Status(nil).Return(nil, errors.New("current status error")), // Getting status after adding a unit
	)
	err = m.addUnits(1)
	if err.Error() != "current status error" {
		t.Errorf("expected current status error but did not get it")
	}
}

func TestRemoveUnit(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockJujuClient := mocks.NewMockJujuClient(ctl)
	unit1 := makeUnit(cloudprovider.InstanceRunning, "unit_1", "unit_1_hostname", "idle", "active", "machine_1")
	units := map[string]*Unit{
		unit1.jujuName: &unit1,
	}
	m, err := makeManager(mockJujuClient, units)
	if err != nil {
		t.Fatalf("error creating manager")
	}

	args := application.DestroyUnitsParams{
		Units:          []string{"unit_1"},
		DestroyStorage: false,
		Force:          false,
	}

	mockJujuClient.EXPECT().DestroyUnits(args).Return(nil, nil)

	err = m.removeUnit("unit_1_hostname")
	if err != nil {
		t.Errorf("error removing unit: %s", err.Error())
	}

	if m.getUnitByHostname("unit_1_hostname").state != cloudprovider.InstanceDeleting {
		t.Errorf("state = %v; want %v", m.getUnitByHostname("unit_1_hostname").state, cloudprovider.InstanceDeleting)
	}

	// Test case when hostname is not found
	err = m.removeUnit("the_host_does_not_exist")
	if err.Error() != "unit with hostname the_host_does_not_exist not found" {
		t.Errorf("error = %v, want %v", "unit with hostname the_host_does_not_exist not found", err.Error())
	}

	// Test error path when calling DestroyUnits
	mockJujuClient.EXPECT().DestroyUnits(gomock.Any()).Return(nil, errors.New("some error"))
	err = m.removeUnit("unit_1_hostname")
	if err.Error() != "some error" {
		t.Errorf("expected some error but did not get it")
	}
}

func TestRefresh(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockJujuClient := mocks.NewMockJujuClient(ctl)
	// Unit 1 will have a missing hostname at first
	unit1 := makeUnit(cloudprovider.InstanceRunning, "unit_1", "", "idle", "active", "machine_1")
	unit2 := makeUnit(cloudprovider.InstanceRunning, "unit_2", "unit_2_hostname", "error", "blocked", "machine_2")
	unit3 := makeUnit(cloudprovider.InstanceCreating, "unit_3", "unit_3_hostname", "idle", "active", "machine_3")
	unit4 := makeUnit(cloudprovider.InstanceDeleting, "unit_4", "unit_4_hostname", "idle", "active", "machine_4")
	units := map[string]*Unit{
		unit1.jujuName: &unit1,
		unit2.jujuName: &unit2,
		unit3.jujuName: &unit3,
		unit4.jujuName: &unit4,
	}
	m, err := makeManager(mockJujuClient, units)
	if err != nil {
		t.Fatalf("error creating manager")
	}

	// Unit 1 now has a hostname
	// Unit 2 has been deleted outside of the autoscaler
	// Unit 5 has been added outside of the autoscaler
	// Unit 6 has been added outside of the autoscaler
	unit1s := makeUnit(cloudprovider.InstanceRunning, "unit_1", "unit_1_hostname", "idle", "active", "machine_1")
	unit3s := makeUnit(cloudprovider.InstanceCreating, "unit_3", "unit_3_hostname", "idle", "active", "machine_3")
	unit4s := makeUnit(cloudprovider.InstanceDeleting, "unit_4", "unit_4_hostname", "idle", "active", "machine_4")
	unit5s := makeUnit(cloudprovider.InstanceCreating, "unit_5", "unit_5_hostname", "idle", "active", "machine_5")
	unit6s := makeUnit(cloudprovider.InstanceCreating, "unit_6", "unit_6_hostname", "executing", "waiting", "machine_6")
	statusUnits := map[string]*Unit{
		unit1s.jujuName: &unit1s,
		unit3s.jujuName: &unit3s,
		unit4s.jujuName: &unit4s,
		unit5s.jujuName: &unit5s,
		unit6s.jujuName: &unit6s,
	}
	ms := makeStatus(m.application, statusUnits)
	gomock.InOrder(
		mockJujuClient.EXPECT().Status(nil).Return(&ms, nil), // Getting previous status
	)

	// unit1 kubeName should be empty before the call
	if m.units[unit1.jujuName].kubeName != "" {
		t.Errorf("before calling refresh: kubeName = %v; want %v", m.units[unit1.jujuName].kubeName, "")
	}

	// unit3 state should be InstanceCreating before the call
	if m.units[unit3.jujuName].state != cloudprovider.InstanceCreating {
		t.Errorf("before calling refresh: state = %v; want %v", m.units[unit3.jujuName].state, cloudprovider.InstanceCreating)
	}

	// unit5s and unit6s should not be managed before the call
	_, ok5 := m.units[unit5s.jujuName]
	_, ok6 := m.units[unit6s.jujuName]
	if ok5 && ok6 {
		t.Errorf("externally added unit exists in manager units before calling refresh")
	}

	// unit2 should be managed before the call
	if _, ok := m.units[unit2.jujuName]; !ok {
		t.Errorf("expected unit does not exist in manager units before calling refresh")
	}

	err = m.refresh()
	if err != nil {
		t.Errorf("error refreshing: %s", err.Error())
	}

	// unit1 kubeName should be unit_1_hostname after the call
	if m.units[unit1.jujuName].kubeName != "unit_1_hostname" {
		t.Errorf("after calling refresh: kubeName = %v; want %v", m.units[unit1.jujuName].kubeName, "unit_1_hostname")
	}

	// unit3 state should now be running since it was previously creating (and active)
	if m.units[unit3.jujuName].state != cloudprovider.InstanceRunning {
		t.Errorf("after calling refresh: state = %v; want %v", m.units[unit3.jujuName].state, cloudprovider.InstanceRunning)
	}

	// unit4 should be deleted
	if _, ok := m.units[unit4.jujuName]; ok {
		t.Errorf("units contain unit that should have been removed")
	}

	// unit2 should be gone as well since it was removed externally
	if _, ok := m.units[unit2.jujuName]; ok {
		t.Errorf("units contain unit that was removed externally")
	}

	// unit5 should now exist in the manager units
	if _, ok := m.units[unit5s.jujuName]; !ok {
		t.Errorf("externally added unit does not exist in manager units")
	}

	// unit6 should still not be in the manager units since it was not active/idle
	if _, ok := m.units[unit6s.jujuName]; ok {
		t.Errorf("non-active/non-idle externally added unit exists in manager units")
	}

	// Test error path when getting status
	// Test the error path for status()
	mockJujuClient.EXPECT().Status(nil).Return(nil, errors.New("status error"))
	err = m.refresh()
	if err.Error() != "status error" {
		t.Errorf("expected status error but did not get it")
	}

}
