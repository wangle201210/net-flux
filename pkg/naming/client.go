package naming

import (
	"strconv"

	"github.com/dellinger2023/net-flux/gen"
	"github.com/dellinger2023/net-flux/pkg/logger"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
)

type DiscoClient interface {
	GetGroupName() string
	RegisterInstance(instance *gen.Instance) error
	DeregisterInstance(serviceName, ip string, port uint64) error
	GetAllServices() ([]string, error)
	GetService(serviceName string, clusters []string) (*gen.Service, error)
	GetServiceInstanceByName(serviceName string) (*gen.Instance, error)
	GetServiceInstance(serviceName string, clusters []string) (*gen.Instance, error)
	GetServiceInstanceByGroup(serviceName, groupName string) (*gen.Instance, error)
	GetServiceInstancesByName(serviceName string) ([]*gen.Instance, error)
	GetServiceInstances(serviceName string, clusters []string) ([]*gen.Instance, error)

	SetConfig(dataId, content string) error
	GetConfig(dataId string) (string, error)
	DeleteConfig(dataId string) error
	ListenConfig(dataId string, onChange func(namespace, group, dataId, data string)) error
	CancelListenConfig(dataId string) error
	SearchConfig(search, dataId string) (*ConfigPage, error)
}

type nacosDiscoClient struct {
	groupName    string
	configClient *NacosConfigClient
	namingClient *NacosNamingClient
}

// GetServiceInstanceByGroup implements DiscoverClient.
func (n *nacosDiscoClient) GetServiceInstanceByGroup(serviceName string, groupName string) (*gen.Instance, error) {
	instance, err := n.namingClient.GetServiceInstance(serviceName, groupName, nil)
	if err != nil {
		return nil, err
	}
	return toInstance(instance)
}

// CancelListenConfig implements DiscoverClient.
func (n *nacosDiscoClient) CancelListenConfig(dataId string) error {
	return n.configClient.CancelListenConfig(dataId, n.groupName)
}

// DeleteConfig implements DiscoverClient.
func (n *nacosDiscoClient) DeleteConfig(dataId string) error {
	return n.configClient.DeleteConfig(dataId, n.groupName)
}

// GetConfig implements DiscoverClient.
func (n *nacosDiscoClient) GetConfig(dataId string) (string, error) {
	return n.configClient.GetConfig(dataId, n.groupName)
}

// ListenConfig implements DiscoverClient.
func (n *nacosDiscoClient) ListenConfig(
	dataId string,
	onChange func(namespace string, group string, dataId string, data string)) error {
	return n.configClient.ListenConfig(dataId, n.groupName, onChange)
}

// SearchConfig implements DiscoverClient.
func (n *nacosDiscoClient) SearchConfig(search string, dataId string) (*ConfigPage, error) {
	configs, err := n.configClient.SearchConfig(search, n.groupName, dataId)
	if err != nil {
		return nil, err
	}

	items := make([]ConfigItem, 0)
	for _, config := range configs.PageItems {
		items = append(items, ConfigItem{
			Id:      string(config.Id),
			DataId:  config.DataId,
			Group:   config.Group,
			Content: config.Content,
		})
	}
	return &ConfigPage{
		TotalCount:     configs.TotalCount,
		PageNumber:     configs.PageNumber,
		PagesAvailable: configs.PagesAvailable,
		PageItems:      items,
	}, nil
}

// SetConfig implements DiscoverClient.
func (n *nacosDiscoClient) SetConfig(dataId string, content string) error {
	return n.configClient.SetConfig(dataId, n.groupName, content)
}

// DeregisterInstance implements DiscoveryClient.
func (n *nacosDiscoClient) DeregisterInstance(
	serviceName string, ip string, port uint64) error {
	return n.namingClient.DeregisterInstance(serviceName, n.groupName, ip, port)
}

// GetAllServices implements DiscoveryClient.
func (n *nacosDiscoClient) GetAllServices() ([]string, error) {
	return n.namingClient.GetAllServices(n.groupName)
}

// GetService implements DiscoveryClient.
func (n *nacosDiscoClient) GetService(
	serviceName string, clusters []string) (*gen.Service, error) {
	service, err := n.namingClient.GetService(serviceName, n.groupName, clusters)
	if err != nil {
		return nil, err
	}

	hosts := make([]*gen.Instance, 0)
	for _, host := range service.Hosts {
		inst, err := toInstance(&host)
		if err != nil {
			logger.Errorf("toInstance error: %v", err)
		}
		hosts = append(hosts, inst)
	}

	return &gen.Service{
		Instances: hosts,
		Cluster:   service.Clusters,
		Name:      service.Name,
		GroupName: service.GroupName,
		Valid:     service.Valid,
	}, nil
}

// GetServiceInstance implements DiscoveryClient.
func (n *nacosDiscoClient) GetServiceInstance(
	serviceName string, clusters []string) (*gen.Instance, error) {
	instance, err := n.namingClient.GetServiceInstance(serviceName, n.groupName, clusters)
	if err != nil {
		return nil, err
	}

	return toInstance(instance)
}

// GetServiceInstanceByName implements DiscoveryClient.
func (n *nacosDiscoClient) GetServiceInstanceByName(serviceName string) (*gen.Instance, error) {
	instance, err := n.namingClient.GetServiceInstanceByName(serviceName)
	if err != nil {
		return nil, err
	}

	return toInstance(instance)
}

// GetServiceInstances implements DiscoveryClient.
func (n *nacosDiscoClient) GetServiceInstances(
	serviceName string, clusters []string) ([]*gen.Instance, error) {
	items, err := n.namingClient.GetServiceInstances(serviceName, n.groupName, clusters)
	if err != nil {
		return nil, err
	}

	instances := make([]*gen.Instance, 0)
	for _, instance := range items {
		inst, err := toInstance(&instance)
		if err != nil {
			logger.Errorf("toInstance error: %v", err)
		}
		instances = append(instances, inst)
	}

	return instances, nil
}

// GetServiceInstancesByName implements DiscoveryClient.
func (n *nacosDiscoClient) GetServiceInstancesByName(serviceName string) ([]*gen.Instance, error) {
	items, err := n.namingClient.GetServiceInstancesByName(serviceName)
	if err != nil {
		return nil, err
	}
	instances := make([]*gen.Instance, 0)
	for _, instance := range items {
		inst, err := toInstance(&instance)
		if err != nil {
			logger.Errorf("toInstance error: %v", err)
		}
		instances = append(instances, inst)
	}

	return instances, nil
}

// RegisterInstance implements DiscoveryClient.
func (n *nacosDiscoClient) RegisterInstance(instance *gen.Instance) error {
	inst := fromInstance(instance)
	return n.namingClient.RegisterInstance(inst.ServiceName, strconv.Itoa(int(instance.Node)),
		inst.Ip, inst.Port, inst.Metadata)
}

// GetGroupName implements DiscoverClient.
func (n *nacosDiscoClient) GetGroupName() string {
	return n.groupName
}

func NewNacosDiscoverClient(cfg DiscoSetting) (DiscoClient, error) {
	configClient, err := NewConfigClient(cfg)
	if err != nil {
		return nil, err
	}
	namingClient, err := NewNamingClient(cfg)
	if err != nil {
		return nil, err
	}
	return &nacosDiscoClient{groupName: cfg.GroupName, configClient: configClient, namingClient: namingClient}, err
}

func toInstance(instance *model.Instance) (*gen.Instance, error) {
	innerIP := instance.Metadata["inner_ip"]
	innerPort, err := strconv.Atoi(instance.Metadata["inner_port"])
	if err != nil {
		return nil, err
	}
	publicIP := instance.Metadata["public_ip"]
	publicPort, err := strconv.Atoi(instance.Metadata["public_port"])
	if err != nil {
		return nil, err
	}

	node, err := strconv.Atoi(instance.Metadata["node"])
	if err != nil {
		return nil, err
	}

	return &gen.Instance{
		InstanceId:  instance.InstanceId,
		PrivateIp:   instance.Ip,
		PrivatePort: int32(instance.Port),
		InnerIp:     innerIP,
		InnerPort:   int32(innerPort),
		PublicIp:    publicIP,
		PublicPort:  int32(publicPort),
		Weight:      float32(instance.Weight),
		Healthy:     instance.Healthy,
		Enable:      instance.Enable,
		Ephemeral:   instance.Ephemeral,
		Node:        int32(node),
		Extra:       instance.Metadata,
	}, nil
}

func fromInstance(instance *gen.Instance) *model.Instance {
	metadata := make(map[string]string)
	metadata["inner_ip"] = instance.InnerIp
	metadata["inner_port"] = strconv.Itoa(int(instance.InnerPort))
	metadata["public_ip"] = instance.PublicIp
	metadata["public_port"] = strconv.Itoa(int(instance.PublicPort))
	metadata["node"] = strconv.Itoa(int(instance.Node))
	return &model.Instance{
		InstanceId:  instance.InstanceId,
		ServiceName: instance.InstanceName,
		Ip:          instance.PrivateIp,
		Port:        uint64(instance.PrivatePort),
		Metadata:    metadata,
	}
}
