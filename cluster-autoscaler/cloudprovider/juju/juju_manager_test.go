package juju

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/juju/juju/rpc/params"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/juju/mocks"
)

var mockUnit1 = Unit{
	state:    cloudprovider.InstanceRunning,
	jujuName: "unit_1",
	kubeName: "unit_1_hostname",
	status: params.UnitStatus{
		AgentStatus: params.DetailedStatus{
			Status: "idle",
		},
		WorkloadStatus: params.DetailedStatus{
			Status: "active",
		},
		Machine: "machine_1",
	},
}

var mockUnit2 = Unit{
	state:    cloudprovider.InstanceRunning,
	jujuName: "unit_2",
	kubeName: "unit_2_hostname",
	status: params.UnitStatus{
		AgentStatus: params.DetailedStatus{
			Status: "error",
		},
		WorkloadStatus: params.DetailedStatus{
			Status: "blocked",
		},
		Machine: "machine_2",
	},
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

	mock_applications := map[string]params.ApplicationStatus{
		"test_application": {
			Units: map[string]params.UnitStatus{
				mockUnit1.jujuName: mockUnit1.status,
				mockUnit2.jujuName: mockUnit2.status,
			},
		},
	}

	mock_machines := map[string]params.MachineStatus{
		mockUnit1.status.Machine: {
			Hostname: mockUnit1.kubeName,
		},
		mockUnit2.status.Machine: {
			Hostname: mockUnit2.kubeName,
		},
	}

	mock_status := params.FullStatus{
		Applications: mock_applications,
		Machines:     mock_machines,
	}

	gomock.InOrder(
		mockJujuClient.EXPECT().Status(nil).Return(&mock_status, nil), // Getting initial status
		mockJujuClient.EXPECT().Status(nil).Return(&mock_status, nil), // Getting hostname of unit_1
		mockJujuClient.EXPECT().Status(nil).Return(&mock_status, nil), // Getting hostname of unit_2
	)

	m := NewManager(mockJujuClient, "test_model", "test_application")
	err := m.init()
	if err != nil {
		t.Errorf("unexpected error returned from init")
	}

	unit := *(m.units[mockUnit1.jujuName])
	opts := []cmp.Option{
		cmp.AllowUnexported(Unit{}),
	}
	if !cmp.Equal(unit, mockUnit1, opts...) {
		t.Errorf("structs are not equal")
	}

}

func TestAddUnit(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockJujuClient := mocks.NewMockJujuClient(ctl)

	m := NewManager(mockJujuClient, "test_model", "test_application")

	mock_applications_before_add := map[string]params.ApplicationStatus{
		"test_application": {
			Units: map[string]params.UnitStatus{
				mockUnit1.jujuName: mockUnit1.status,
			},
		},
	}

	mock_machines_before_add := map[string]params.MachineStatus{
		mockUnit1.status.Machine: {
			Hostname: mockUnit1.kubeName,
		},
	}

	mock_status_before_add := params.FullStatus{
		Applications: mock_applications_before_add,
		Machines:     mock_machines_before_add,
	}

	mock_applications_after_add := map[string]params.ApplicationStatus{
		"test_application": {
			Units: map[string]params.UnitStatus{
				mockUnit1.jujuName: mockUnit1.status,
				mockUnit2.jujuName: mockUnit2.status,
			},
		},
	}

	mock_machines_after_add := map[string]params.MachineStatus{
		mockUnit1.status.Machine: {
			Hostname: mockUnit1.kubeName,
		},
		mockUnit2.status.Machine: {
			Hostname: mockUnit2.kubeName,
		},
	}

	mock_status_after_add := params.FullStatus{
		Applications: mock_applications_after_add,
		Machines:     mock_machines_after_add,
	}

	gomock.InOrder(
		mockJujuClient.EXPECT().Status(nil).Return(&mock_status_before_add, nil),       // Getting previous status
		mockJujuClient.EXPECT().AddUnits(gomock.Any()).Return([]string{"unit_2"}, nil), // Add a unit
		mockJujuClient.EXPECT().Status(nil).Return(&mock_status_after_add, nil),        // Getting status after adding a unit
	)

	m.addUnits(1)
	if _, ok := m.units[mockUnit2.jujuName]; !ok {
		t.Errorf("units does not contain added unit")
	}

}
