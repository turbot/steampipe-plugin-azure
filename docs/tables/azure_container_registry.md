---
title: "Steampipe Table: azure_container_registry - Query Azure Container Registries using SQL"
description: "Allows users to query Azure Container Registries for detailed information about their configuration, status, and associated metadata."
---

# Table: azure_container_registry - Query Azure Container Registries using SQL

Azure Container Registry is a managed Docker registry service provided by Microsoft Azure for storing and managing private Docker container images and related artifacts. It allows you to build, store, and manage container images and artifacts in a private registry for all types of container deployments. This service also integrates well with existing container development and deployment pipelines.

## Table Usage Guide

The 'azure_container_registry' table provides insights into Container Registries within Microsoft Azure. As a DevOps engineer, explore registry-specific details through this table, including SKU, login server, creation date, and associated metadata. Utilize it to uncover information about registries, such as those with admin user enabled, the network rule set, and the encryption status. The schema presents a range of attributes of the Container Registry for your analysis, like the registry name, resource group, region, and associated tags.

## Examples

### Basic info
Explore the status and details of your Azure Container Registry. This query can help you assess the creation date, region, and the tier of your registry, providing insights into your resource usage and allocation.

```sql
select
  name,
  id,
  provisioning_state,
  creation_date,
  sku_tier,
  region
from
  azure_container_registry;
```

### List registries not encrypted with a customer-managed key
Explore which Azure container registries are not encrypted with a customer-managed key. This is useful for identifying potential security vulnerabilities in your Azure cloud environment.

```sql
select
  name,
  encryption ->> 'status' as encryption_status,
  region
from
  azure_container_registry;
```

### List registries not configured with virtual network service endpoint
Analyze the settings to understand which Azure Container Registries are not configured with a virtual network service endpoint. This is useful to pinpoint potential security gaps where data might be exposed to untrusted networks.

```sql
select
  name,
  network_rule_set ->> 'defaultAction' as network_rule_default_action,
  network_rule_set ->> 'virtualNetworkRules' as virtual_network_rules
from
  azure_container_registry
where
  network_rule_set is not null
  and network_rule_set ->> 'defaultAction' = 'Allow';
```

### List registries with admin user account enabled
Discover the segments where the admin user account is enabled in the Azure container registry. This is useful for identifying potential security risks and ensuring the proper configuration of user permissions.

```sql
select
  name,
  admin_user_enabled,
  region
from
  azure_container_registry
where
  admin_user_enabled;
```