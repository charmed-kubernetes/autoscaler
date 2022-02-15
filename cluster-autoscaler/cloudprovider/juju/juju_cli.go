package juju

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	klog "k8s.io/klog/v2"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Unit struct {
	state      cloudprovider.InstanceState
	jujuName   string
	kubeName   string
	workload   string
	agent      string
	registered bool
}

type Manager struct {
	model       string
	application string
	units       map[string]*Unit
}

func (m *Manager) init() error {
	var status []byte
	var hostname string

	status, _ = exec.Command("juju", "status", "-m", m.model, m.application).Output()
	for _, line := range strings.Split(string(status), "\n") {
		if strings.Contains(line, m.application+"/") {
			info := strings.Fields(line)
			unitName := strings.Replace(info[0], "*", "", -1)
			nodeExec, _ := exec.Command("juju", "exec", "-m", m.model, "-u", unitName, "hostname").Output()
			hostname = strings.Fields(string(nodeExec))[0]
			exec.Command("kubectl", "patch", "node", hostname, "-p", `{"spec":{"providerID":"`+hostname+`"}}`).Output()
			m.units[unitName] = &Unit{
				state:      cloudprovider.InstanceRunning,
				jujuName:   unitName,
				kubeName:   hostname,
				registered: true,
			}
		}
	}

	return nil
}

func (m *Manager) addUnits(delta int) error {
	juju, _ := exec.LookPath("juju")

	prevStatus := m.getStatus()

	cmd := exec.Cmd{
		Path:   juju,
		Args:   []string{juju, "add-unit", "-m", m.model, "-n", strconv.Itoa(delta), m.application},
		Stderr: os.Stdout,
	}
	cmd.Run()

	for key, _ := range m.getStatus() {
		if _, ok := prevStatus[key]; !ok {
			m.units[key] = &Unit{
				state:    cloudprovider.InstanceCreating,
				jujuName: key,
			}
		}
	}

	return nil
}

func (m *Manager) removeUnit(name string) error {
	juju, _ := exec.LookPath("juju")
	unit := m.getUnit(name)
	unit.state = cloudprovider.InstanceDeleting

	cmd := exec.Cmd{
		Path:   juju,
		Args:   []string{juju, "run-action", "-m", m.model, unit.jujuName, "pause", "--wait"},
		Stderr: os.Stdout,
	}
	cmd.Run()

	cmd = exec.Cmd{
		Path:   juju,
		Args:   []string{juju, "remove-unit", "-m", m.model, unit.jujuName},
		Stderr: os.Stdout,
	}
	cmd.Run()

	return nil
}

func (m *Manager) refresh() error {
	for key, val := range m.getStatus() {
		if _, ok := m.units[key]; ok {
			m.units[key].agent = val[0]
			m.units[key].workload = val[1]
		}
	}

	for _, unit := range m.units {
		if unit.state == cloudprovider.InstanceCreating {
			if unit.kubeName == "" {
				nodeExec, _ := exec.Command("juju", "exec", "-m", m.model, "-u", unit.jujuName, "hostname").Output()
				if len(strings.Fields(string(nodeExec))) > 0 {
					unit.kubeName = strings.Fields(string(nodeExec))[0]
				}
			}

			if unit.workload == "active" && !unit.registered {
				output, _ := exec.Command("kubectl", "patch", "node", unit.kubeName, "-p", `{"spec":{"providerID":"`+unit.kubeName+`"}}`).Output()
				if string(output) == "node/"+unit.kubeName+" patched" {
					unit.registered = true
					unit.state = cloudprovider.InstanceRunning
					klog.Warningf(unit.kubeName + " registered.")
				}
			}
		} else if unit.state == cloudprovider.InstanceDeleting {
			delete(m.units, unit.jujuName)
		}
	}

	return nil
}

func (m *Manager) getUnit(name string) *Unit {
	for _, unit := range m.units {
		if unit.kubeName == name {
			return unit
		}
	}
	return nil
}

func (m *Manager) getStatus() map[string][]string {
	var status []byte
	units := make(map[string][]string)

	status, _ = exec.Command("juju", "status", "-m", m.model, m.application).Output()
	for _, line := range strings.Split(string(status), "\n") {
		if strings.Contains(line, m.application+"/") {
			info := strings.Fields(line)
			unitName := strings.Replace(info[0], "*", "", -1)
			if info[1] == "terminated" {
				continue
			} else {
				units[unitName] = info[0:]
			}
		}
	}
	return units
}
