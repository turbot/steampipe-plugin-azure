# Table: azure_spring_cloud_service

Azure Spring Cloud is a platform as a service (PaaS) for Spring developers. It manages the lifecycle of your Spring Boot applications with comprehensive monitoring and diagnostics, configuration management, service discovery, CI/CD integration, blue-green deployments and more.

## Examples

### Basic info

```sql
select
  id,
  name,
  type,
  provisioning_state,
  sku_name,
  sku_tier,
  version
from
  azure_spring_cloud_service;
```

### List network profile details

```sql
select
  id,
  name,
  network_profile ->> 'AppNetworkResourceGroup' as network_profile_app_network_resource_group,
  network_profile ->> 'AppSubnetID' as network_profile_app_subnet_id,
  jsonb_pretty(network_profile -> 'OutboundPublicIPs') as network_profile_outbound_public_ips,
  network_profile ->> 'ServiceCidr' as network_profile_service_cidr,
  network_profile ->> 'ServiceRuntimeNetworkResourceGroup' as network_profile_service_runtime_network_resource_group,
  network_profile ->> 'ServiceRuntimeSubnetID' as network_profile_service_runtime_subnet_id
from
  azure_spring_cloud_service;
```