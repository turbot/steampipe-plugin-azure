# Table: azure_service_fabric_cluster

An Azure Service Fabric cluster is a network-connected set of virtual or physical machines into which your microservices are deployed and managed. It rebalances the partition replicas and instances across the increase or decreased number of nodes to make better use of the hardware on each node. It allows for the creation of clusters on any VMs or computers running Windows Server or Linux.

## Examples

### Basic info

```sql
select
  name,
  id,
  provisioning_state, 
  type,
  cluster_code_version,
  management_endpoint,
  upgrade_mode,
  vm_image
from
  azure_service_fabric_cluster;
```

### List azure active directory details for clusters

```sql
select
  name,
  id,
  azure_active_directory ->> 'clientApplication' as client_application,
  azure_active_directory ->> 'clusterApplication' as cluster_application,
  azure_active_directory ->> 'tenantId' as tenant_id
from
  azure_service_fabric_cluster;
```

### List certificate details for clusters

```sql
select
  name,
  id,
  certificate ->> 'thumbprint' as thumbprint,
  certificate ->> 'thumbprintSecondary' as thumbprint_secondary,
  certificate ->> 'x509StoreName' as x509_store_name
from
  azure_service_fabric_cluster;
```

### List fabric setting details for clusters

```sql
select
  name,
  id,
  settings ->> 'name' as settings_name,
  jsonb_pretty(settings -> 'parameters') as settings_parameters
from
  azure_service_fabric_cluster,
  jsonb_array_elements(fabric_settings) as settings;
```

### List node type details for clusters

```sql
select
  name,
  id,
  types ->> 'clientConnectionEndpointPort' as type_client_connection_endpoint_port,
  types ->> 'durabilityLevel' as type_durability_level,
  types -> 'httpGatewayEndpointPort' as type_http_gateway_endpoint_port,
  types -> 'isPrimary' as type_is_primary,
  types ->> 'name' as type_name,
  types -> 'vmInstanceCount' as type_vm_instance_count,
  jsonb_pretty(types -> 'applicationPorts') as settings_application_ports,
  jsonb_pretty(types -> 'ephemeralPorts') as settings_ephemeral_ports
from
  azure_service_fabric_cluster,
  jsonb_array_elements(node_types) as types;
```
