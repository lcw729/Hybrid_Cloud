#!/usr/bin/env python
from constructs import Construct
from cdktf import App, TerraformStack
from imports.azurerm import \
    AzurermProvider, \
    KubernetesCluster, \
    KubernetesClusterDefaultNodePool, \
    KubernetesClusterIdentity, ResourceGroupConfig, AzurermProviderFeatures

class MyStack(TerraformStack):
    def __init__(self, scope: Construct, ns: str):
        super().__init__(scope, ns)

        # define resources here
        features = AzurermProviderFeatures()
        provider = AzurermProvider(self, 'azure', features=[features])

        node_pool = KubernetesClusterDefaultNodePool(
            name='default', node_count=1, vm_size='Standard_D2_v3')
 
        resource_group = ResourceGroupConfig(name='openmcpResourceGroup', location='australiaeast')

        identity = KubernetesClusterIdentity(type='SystemAssigned')

        cluster = KubernetesCluster(
            self, 'hcp-cluster',
            name='hcp-cluster',
            default_node_pool=[node_pool],
            dns_prefix='test',
            location=resource_group.location,
            resource_group_name=resource_group.name,
            identity=[identity],
            tags={"foo": "bar"}
        )


app = App()
MyStack(app, "create-cluster-azure")

app.synth()
