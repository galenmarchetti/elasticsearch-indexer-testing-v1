/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package single_node_example_network

import (
	"github.com/gmarchetti/elasticsearch-indexer-testing-v1/elasticsearch_indexer/services"
	"github.com/kurtosis-tech/kurtosis-go/lib/networks"
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"github.com/palantir/stacktrace"
)

const (
	vanillaConfigId    networks.ConfigurationID = "vanilla"

	theNodeServiceId          networks.ServiceID = "the-node"
	theNodeStopTimeoutSeconds                    = 30
)

// =================================== NETWORK ===================================
type SingleNodeExampleNetwork struct{
	rawNetwork *networks.ServiceNetwork
	theNodeAdded bool
}

func NewSingleNodeExampleNetwork(rawNetwork *networks.ServiceNetwork) *SingleNodeExampleNetwork {
	return &SingleNodeExampleNetwork{
		rawNetwork: rawNetwork,
		theNodeAdded: false,
	}
}

func (network *SingleNodeExampleNetwork) AddTheNode() (services.ExampleService, error) {
	if network.theNodeAdded {
		return nil, stacktrace.NewError("The node is already added")
	}
	// TODO add example with dependencies
	availabilityChecker, err := network.rawNetwork.AddService(vanillaConfigId, theNodeServiceId, map[networks.ServiceID]bool{})
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred adding the node")
	}
	if err := availabilityChecker.WaitForStartup(); err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred waiting for the node to come up")
	}
	network.theNodeAdded = true

	serviceNode, err := network.rawNetwork.GetService(theNodeServiceId)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting the node's service information")
	}
	castedService := serviceNode.Service.(services.ExampleService)
	return castedService, nil
}

func (network *SingleNodeExampleNetwork) RemoveTheNode() error {
	if !network.theNodeAdded {
		return stacktrace.NewError("The node hasn't been added yet")
	}
	if err := network.rawNetwork.RemoveService(theNodeServiceId, theNodeStopTimeoutSeconds); err != nil {
		return stacktrace.NewError("An error occurred removing the node from the network")
	}
	network.theNodeAdded = false
	return nil
}

// =================================== NETWORK LOADER ===================================
type SingleNodeExampleNetworkLoader struct {
	serviceImage string
}

func NewSingleNodeExampleNetworkLoader(serviceImage string) *SingleNodeExampleNetworkLoader {
	return &SingleNodeExampleNetworkLoader{serviceImage: serviceImage}
}


func (loader SingleNodeExampleNetworkLoader) ConfigureNetwork(builder *networks.ServiceNetworkBuilder) error {
	builder.AddConfiguration(
		vanillaConfigId,
		loader.serviceImage,
		services.ExampleServiceInitializerCore{},
		services.ExampleAvailabilityCheckerCore{})
	return nil
}

func (loader SingleNodeExampleNetworkLoader) InitializeNetwork(network *networks.ServiceNetwork) (map[networks.ServiceID]services.ServiceAvailabilityChecker, error) {
	return map[networks.ServiceID]services.ServiceAvailabilityChecker{}, nil
}

func (loader SingleNodeExampleNetworkLoader) WrapNetwork(network *networks.ServiceNetwork) (networks.Network, error) {
	return *NewSingleNodeExampleNetwork(network), nil
}

