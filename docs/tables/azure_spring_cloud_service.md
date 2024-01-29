---
title: "Steampipe Table: azure_spring_cloud_service - Query Azure Spring Cloud Services using SQL"
description: "Allows users to query Azure Spring Cloud Services, specifically the details of the service instances, providing insights into the configuration and state of the Spring Cloud applications."
---

# Table: azure_spring_cloud_service - Query Azure Spring Cloud Services using SQL

Azure Spring Cloud is a fully managed service for Spring Boot apps that lets you focus on building and running the apps that run your business without the hassle of managing infrastructure. It provides a platform for deploying and managing Spring Boot and Spring Cloud applications in the cloud. The service is jointly built, operated, and supported by Pivotal Software and Microsoft to provide a native platform designed to be easily run and managed on Azure.

## Table Usage Guide

The `azure_spring_cloud_service` table provides insights into Azure Spring Cloud Services within Microsoft Azure. As a DevOps engineer, explore service-specific details through this table, including configurations, provisioning state, and associated metadata. Utilize it to uncover information about services, such as service versions, the network profile of the service, and the verification of service configurations.

## Examples

### Basic info
Explore the fundamental details of your Azure Spring Cloud services, such as their provisioning state, SKU details, and version. This can be used to assess the status and tier of your services, enabling effective management and optimization.

```sql+postgres
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

```sql+sqlite
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
Assess the elements within your Azure Spring Cloud Service's network profile. This query can be used to gain insights into the specific configurations and resource groups associated with your network profile, which can aid in network management and troubleshooting.

```sql+postgres
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

```sql+sqlite
select
  id,
  name,
  json_extract(network_profile, '$.AppNetworkResourceGroup') as network_profile_app_network_resource_group,
  json_extract(network_profile, '$.AppSubnetID') as network_profile_app_subnet_id,
  network_profile -> 'OutboundPublicIPs' as network_profile_outbound_public_ips,
  json_extract(network_profile, '$.ServiceCidr') as network_profile_service_cidr,
  json_extract(network_profile, '$.ServiceRuntimeNetworkResourceGroup') as network_profile_service_runtime_network_resource_group,
  json_extract(network_profile, '$.ServiceRuntimeSubnetID') as network_profile_service_runtime_subnet_id
from
  azure_spring_cloud_service;
```