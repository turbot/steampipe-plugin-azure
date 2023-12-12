---
title: "Steampipe Table: azure_service_fabric_cluster - Query Azure Service Fabric Clusters using SQL"
description: "Allows users to query Azure Service Fabric Clusters, providing insights into the structure, health, and configuration of each cluster."
---

# Table: azure_service_fabric_cluster - Query Azure Service Fabric Clusters using SQL

Azure Service Fabric is a distributed systems platform that makes it easy to package, deploy, and manage scalable and reliable microservices and containers. It also provides comprehensive runtime and lifecycle management capabilities to applications that are composed of these microservices or containers. This makes it an ideal tool for developers and administrators looking to manage complex microservices architectures.

## Table Usage Guide

The `azure_service_fabric_cluster` table provides insights into Service Fabric Clusters within Azure. Developers and administrators can explore cluster-specific details through this table, including the cluster's health, configuration, and node types. Utilize it to uncover information about clusters, such as their reliability tier, upgrade mode, and the version of Service Fabric they're running.

## Examples

### Basic info
Explore the status and configuration of your Azure Service Fabric clusters to understand their current operational state and setup. This is crucial for managing your clusters effectively and ensuring they are configured according to your organization's standards.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that contain key Azure Active Directory details for specific clusters. This is useful for understanding the configuration and security measures of your Azure Service Fabric clusters.

```sql+postgres
select
  name,
  id,
  azure_active_directory ->> 'clientApplication' as client_application,
  azure_active_directory ->> 'clusterApplication' as cluster_application,
  azure_active_directory ->> 'tenantId' as tenant_id
from
  azure_service_fabric_cluster;
```

```sql+sqlite
select
  name,
  id,
  json_extract(azure_active_directory, '$.clientApplication') as client_application,
  json_extract(azure_active_directory, '$.clusterApplication') as cluster_application,
  json_extract(azure_active_directory, '$.tenantId') as tenant_id
from
  azure_service_fabric_cluster;
```

### List certificate details for clusters
Determine the security status of your clusters by examining the details of their associated certificates. This is useful for ensuring the integrity and validity of your clusters' security certificates.

```sql+postgres
select
  name,
  id,
  certificate ->> 'thumbprint' as thumbprint,
  certificate ->> 'thumbprintSecondary' as thumbprint_secondary,
  certificate ->> 'x509StoreName' as x509_store_name
from
  azure_service_fabric_cluster;
```

```sql+sqlite
select
  name,
  id,
  json_extract(certificate, '$.thumbprint') as thumbprint,
  json_extract(certificate, '$.thumbprintSecondary') as thumbprint_secondary,
  json_extract(certificate, '$.x509StoreName') as x509_store_name
from
  azure_service_fabric_cluster;
```

### List fabric setting details for clusters
Determine the configuration details for your clusters in Azure Service Fabric. This can help you understand and manage the settings parameters for each cluster, ensuring optimal performance and security.

```sql+postgres
select
  name,
  id,
  settings ->> 'name' as settings_name,
  jsonb_pretty(settings -> 'parameters') as settings_parameters
from
  azure_service_fabric_cluster,
  jsonb_array_elements(fabric_settings) as settings;
```

```sql+sqlite
select
  name,
  c.id,
  json_extract(settings.value, '$.name') as settings_name,
  json_extract(settings.value, '$.parameters') as settings_parameters
from
  azure_service_fabric_cluster as c,
  json_each(fabric_settings) as settings;
```

### List node type details for clusters
Explore the characteristics of different nodes within your Azure Service Fabric Clusters. This query helps you understand the configuration and capabilities of each node, which can be beneficial for managing resources and optimizing performance.

```sql+postgres
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

```sql+sqlite
select
  name,
  c.id,
  json_extract(types.value, '$.clientConnectionEndpointPort') as type_client_connection_endpoint_port,
  json_extract(types.value, '$.durabilityLevel') as type_durability_level,
  json_extract(types.value, '$.httpGatewayEndpointPort') as type_http_gateway_endpoint_port,
  json_extract(types.value, '$.isPrimary') as type_is_primary,
  json_extract(types.value, '$.name') as type_name,
  json_extract(types.value, '$.vmInstanceCount') as type_vm_instance_count,
  types.value as settings_application_ports,
  types.value as settings_ephemeral_ports
from
  azure_service_fabric_cluster as c,
  json_each(node_types) as types;
```