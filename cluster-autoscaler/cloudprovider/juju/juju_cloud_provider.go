/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package juju

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/juju/juju/api/connector"
	"gopkg.in/yaml.v2"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/config"
	"k8s.io/autoscaler/cluster-autoscaler/config/dynamic"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	kube_client "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	klog "k8s.io/klog/v2"
)

var _ cloudprovider.CloudProvider = (*jujuCloudProvider)(nil)

const (
	GPULabel             = "juju/gpu-node" // GPULabel is the label added to nodes with GPU resource.
	scaleToZeroSupported = true
)

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type jujuCloudConfig struct {
	User      string   `yaml:"user"`
	Password  string   `yaml:"password"`
	Endpoints []string `yaml:"endpoints"`
	CAcert    string   `yaml:"ca-cert"`
}

// jujuCloudProvider implements CloudProvider interface.
type jujuCloudProvider struct {
	resourceLimiter *cloudprovider.ResourceLimiter
	nodeGroups      []cloudprovider.NodeGroup
}

func newJujuCloudProvider(rl *cloudprovider.ResourceLimiter, nodeGroups []cloudprovider.NodeGroup) (*jujuCloudProvider, error) {
	return &jujuCloudProvider{
		resourceLimiter: rl,
		nodeGroups:      nodeGroups,
	}, nil
}

// Name returns name of the cloud provider.
func (j *jujuCloudProvider) Name() string {
	return cloudprovider.JujuProviderName
}

// NodeGroups returns all node groups configured for this cloud provider.
func (j *jujuCloudProvider) NodeGroups() []cloudprovider.NodeGroup {
	return j.nodeGroups
}

// NodeGroupForNode returns the node group for the given node, nil if the node
// should not be processed by cluster autoscaler, or non-nil error if such
// occurred. Must be implemented.
func (j *jujuCloudProvider) NodeGroupForNode(node *apiv1.Node) (cloudprovider.NodeGroup, error) {
	for _, nodeGroup := range j.nodeGroups {
		nodeGroupNodes, err := nodeGroup.Nodes()
		if err != nil {
			return nil, err
		}
		for _, nodeGroupNode := range nodeGroupNodes {
			if nodeGroupNode.Id == node.Spec.ProviderID {
				return nodeGroup, nil
			}
		}
	}
	return nil, nil
}

// Pricing returns pricing model for this cloud provider or error if not
// available. Implementation optional.
func (j *jujuCloudProvider) Pricing() (cloudprovider.PricingModel, errors.AutoscalerError) {
	return nil, cloudprovider.ErrNotImplemented
}

// GetAvailableMachineTypes get all machine types that can be requested from
// the cloud provider. Implementation optional.
func (j *jujuCloudProvider) GetAvailableMachineTypes() ([]string, error) {
	return []string{}, nil
}

// NewNodeGroup builds a theoretical node group based on the node definition
// provided. The node group is not automatically created on the cloud provider
// side. The node group is not returned by NodeGroups() until it is created.
// Implementation optional.
func (j *jujuCloudProvider) NewNodeGroup(
	machineType string,
	labels map[string]string,
	systemLabels map[string]string,
	taints []apiv1.Taint,
	extraResources map[string]resource.Quantity,
) (cloudprovider.NodeGroup, error) {
	return nil, cloudprovider.ErrNotImplemented
}

// GetResourceLimiter returns struct containing limits (max, min) for
// resources (cores, memory etc.).
func (j *jujuCloudProvider) GetResourceLimiter() (*cloudprovider.ResourceLimiter, error) {
	return j.resourceLimiter, nil
}

// GPULabel returns the label added to nodes with GPU resource.
func (j *jujuCloudProvider) GPULabel() string {
	return GPULabel
}

// GetAvailableGPUTypes return all available GPU types cloud provider supports.
func (j *jujuCloudProvider) GetAvailableGPUTypes() map[string]struct{} {
	return nil
}

// Cleanup cleans up open resources before the cloud provider is destroyed,
// i.e. go routines etc.
func (j *jujuCloudProvider) Cleanup() error {
	return nil
}

// Refresh is called before every main loop and can be used to dynamically
// update cloud provider state. In particular the list of node groups returned
// by NodeGroups() can change as a result of CloudProvider.Refresh().
func (j *jujuCloudProvider) Refresh() error {
	// Juju does not have a NodeGroup concept. Currently, a "node group" is added on startup for each
	// --nodes option passed to the autoscaler command
	// The units making up these node groups can change dynamically if a juju admin adds
	// or removes units outside of the autoscaler
	// The loop below calls the refresh function for each node group (which updates state to include any externally added or removed nodes),
	// and updates the target size to match the current size of the node group
	klog.Infof("refreshing node groups")
	for _, node := range j.nodeGroups {
		// Cast the cloudprovider.NodeGroup interface to the underlying juju NodeGroup struct
		jujuNodeGroup, ok := node.(*NodeGroup)
		if ok {
			klog.Infof("updating node group %s target", jujuNodeGroup.id)
			jujuNodeGroup.manager.refresh()
			jujuNodeGroup.target = len(jujuNodeGroup.manager.units)
		}
	}

	return nil
}

// BuildJuju builds the Juju cloud provider.
func BuildJuju(
	opts config.AutoscalingOptions,
	do cloudprovider.NodeGroupDiscoveryOptions,
	rl *cloudprovider.ResourceLimiter,
) cloudprovider.CloudProvider {
	flag.Parse()
	kubeClient := createKubeClient(getKubeConfig())

	var configRC io.ReadCloser
	if opts.CloudConfig != "" {
		var err error
		configRC, err = os.Open(opts.CloudConfig)
		if err != nil {
			klog.Fatalf("Couldn't open cloud provider configuration %s: %#v", opts.CloudConfig, err)
		}
		defer configRC.Close()
	}

	jujuConfig, err := readCloudConfigYaml(configRC)
	if err != nil {
		klog.Fatalf("Couldn't read cloud provider configuration yaml file %s", err)
	}

	ngs := []cloudprovider.NodeGroup{}
	for _, nodeGroupSpecString := range do.NodeGroupSpecs {
		nodeGroupSpec, err := dynamic.SpecFromString(nodeGroupSpecString, scaleToZeroSupported)
		if err != nil {
			klog.Fatalf("failed to parse node group spec: %v", err)
			continue
		}
		model, application, err := parseNodeGroupName(nodeGroupSpec.Name)
		if err != nil {
			klog.Fatalf("failed to parse node group name: %v", err)
			continue
		}

		connector, err := connector.NewSimple(connector.SimpleConfig{
			ControllerAddresses: jujuConfig.Endpoints,
			CACert:              jujuConfig.CAcert,
			ModelUUID:           model,
			Username:            jujuConfig.User,
			Password:            jujuConfig.Password,
		})

		if err != nil {
			klog.Fatalf("failed to create simple connector %v", err)
			continue
		}

		jujuAPI, err := NewJujuAPi(connector)
		if err != nil {
			klog.Fatalf("failed to create JujuClient %v", err)
			continue
		}

		man, err := NewManager(jujuAPI, kubeClient, model, application)
		if err != nil {
			klog.Fatalf("error creating manager: %v", err)
			continue
		}

		jujuID := fmt.Sprintf("juju-%s-%s", model, application)
		ng := &NodeGroup{
			id:      jujuID,
			minSize: nodeGroupSpec.MinSize,
			maxSize: nodeGroupSpec.MaxSize,
			target:  len(man.units),
			manager: man,
		}
		ngs = append(ngs, ng)
	}

	provider, err := newJujuCloudProvider(rl, ngs)
	if err != nil {
		klog.Fatalf("Failed to create Juju cloud provider: %v", err)
	}

	return provider
}

func parseNodeGroupName(name string) (string, string, error) {
	s := strings.Split(name, ":")
	if len(s) != 2 {
		return "", "", fmt.Errorf("failed to parse node group name: %s, expected <model>:<application>", name)
	}
	model := s[0]
	application := s[1]
	return model, application, nil
}

func readCloudConfigYaml(configRC io.ReadCloser) (jujuCloudConfig, error) {
	t := jujuCloudConfig{}
	b, err := io.ReadAll(configRC)
	if err != nil {
		return t, err
	}

	err = yaml.Unmarshal(b, &t)
	return t, err
}

func getKubeConfig() *rest.Config {
	kubeConfigFile := flag.Lookup("kubeconfig").Value.(flag.Getter).Get().(string)
	if kubeConfigFile != "" {
		klog.V(1).Infof("Using kubeconfig file: %s", kubeConfigFile)
		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", kubeConfigFile)
		if err != nil {
			klog.Fatalf("Failed to build config: %v", err)
		}
		return config
	}
	kubernetes := flag.Lookup("kubernetes").Value.(flag.Getter).Get().(string)
	url, err := url.Parse(kubernetes)
	if err != nil {
		klog.Fatalf("Failed to parse Kubernetes url: %v", err)
	}

	kubeConfig, err := config.GetKubeClientConfig(url)
	if err != nil {
		klog.Fatalf("Failed to build Kubernetes client configuration: %v", err)
	}

	return kubeConfig
}

func createKubeClient(kubeConfig *rest.Config) kube_client.Interface {
	return kube_client.NewForConfigOrDie(kubeConfig)
}
