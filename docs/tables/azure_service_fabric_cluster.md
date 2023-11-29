---
title: "Steampipe Table: azure_service_fabric_cluster - Query Azure Service Fabric Clusters using SQL"
description: "Allows users to query Azure Service Fabric Clusters."
---

# Table: azure_service_fabric_cluster - Query Azure Service Fabric Clusters using SQL

Azure Service Fabric is a distributed systems platform that makes it easy to package, deploy, and manage scalable and reliable microservices and containers. It also provides comprehensive runtime and lifecycle management capabilities to applications that are composed of these microservices or containers. This platform simplifies the delivery of cloud services and provides developers with a comprehensive, agnostic and intrinsically secure approach to building, scaling and updating cloud applications.

## Table Usage Guide

The 'azure_service_fabric_cluster' table provides insights into Service Fabric Clusters within Azure Service Fabric. As a DevOps engineer, explore cluster-specific details through this table, including cluster code versions, reliability levels, upgrade modes, and associated metadata. Utilize it to uncover information about clusters, such as those with specific reliability levels, the upgrade modes of the clusters, and the verification of cluster health policies. The schema presents a range of attributes of the Service Fabric Cluster for your analysis, like the cluster ID, creation date, upgrade mode, and associated tags.

## Examples

### Basic info
Explore which Azure Service Fabric Clusters are being used by reviewing their basic information. This helps in managing resources and understanding their provisioning states and upgrade modes.

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
Explore the Azure Active Directory details associated with your clusters to understand the client and cluster applications. This can be beneficial for assessing the elements within your Azure Service Fabric Cluster, such as the tenant ID.

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
Discover the segments that have specific certificate details for clusters. This can be useful in identifying potential security vulnerabilities or ensuring compliance with organizational policies.

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
Analyze the settings to understand the configuration details for specific clusters within the Azure Service Fabric. This can help in managing and troubleshooting your service fabric clusters effectively.

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
Assess the configuration of cluster nodes to better understand their connection points, durability levels, and port settings. This information can be useful for optimizing resource allocation and enhancing network security.

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