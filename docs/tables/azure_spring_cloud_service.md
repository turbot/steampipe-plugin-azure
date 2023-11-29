---
title: "Steampipe Table: azure_spring_cloud_service - Query Azure Spring Cloud Services using SQL"
description: "Allows users to query Azure Spring Cloud Services, providing data on various aspects of the service such as the service's ID, name, type, and location, as well as detailed information on the service's properties, tags, and encryption settings."
---

# Table: azure_spring_cloud_service - Query Azure Spring Cloud Services using SQL

Azure Spring Cloud is a service that lets developers build, deploy, and scale Spring Boot applications on Azure. It provides a fully managed service for Spring Boot apps, allowing developers to focus on building their applications without the worry of managing infrastructure. Azure Spring Cloud is designed to be simple, safe, and scalable, providing a robust platform for enterprise-grade applications.

## Table Usage Guide

The 'azure_spring_cloud_service' table provides insights into Azure Spring Cloud Services. As a DevOps engineer, explore service-specific details through this table, including service properties, tags, and encryption settings. Utilize it to uncover information about services, such as the service's ID, name, type, and location. The schema presents a range of attributes of the Azure Spring Cloud Service for your analysis, like the service's provisioning state, active deployment name, and network profile.

## Examples

### Basic info
Explore the various features of your Azure Spring Cloud services, such as their current provisioning state, type, and version. This can help you manage and optimize your resources effectively.

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
This example helps you explore the details of your network profile in Azure Spring Cloud Service. It's particularly useful when you need to understand your network configuration for troubleshooting or optimizing your cloud services.

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