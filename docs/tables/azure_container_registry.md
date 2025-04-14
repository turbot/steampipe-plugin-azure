---
title: "Steampipe Table: azure_container_registry - Query Azure Container Registries using SQL"
description: "Allows users to query Azure Container Registries, providing insights into the status, SKU, network access, and other critical details."
folder: "Container Registry"
---

# Table: azure_container_registry - Query Azure Container Registries using SQL

Azure Container Registry is a managed Docker registry service provided by Microsoft Azure for storing and managing Docker images. It is integrated with Azure DevOps, Azure Kubernetes Service (AKS), Docker CLI, and other popular open-source tools. Azure Container Registry allows developers to build, store, and manage container images for Azure deployments in a central registry.

## Table Usage Guide

The `azure_container_registry` table provides insights into Azure Container Registries within Microsoft Azure. As a DevOps engineer, explore registry-specific details through this table, including the status, SKU, network access, and other critical details. Utilize it to uncover information about registries, such as those with private network access, the SKU tier, and the verification of admin user-enabled status.

## Examples

### Basic info
Explore the status and details of your Azure Container Registry instances, including their creation date and geographical location, to gain insights into the distribution and management of your resources. This can be particularly useful for auditing purposes, resource allocation, and strategizing regional deployment.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which container registries in your Azure environment are not encrypted with a customer-managed key. This can help in identifying potential security gaps and ensuring better data protection.

```sql+postgres
select
  name,
  encryption ->> 'status' as encryption_status,
  region
from
  azure_container_registry;
```

```sql+sqlite
select
  name,
  json_extract(encryption, '$.status') as encryption_status,
  region
from
  azure_container_registry;
```

### Get webhook details of registries
Webhooks in Azure Container Registry provide a way to trigger custom actions in response to events happening within the registry. These events can include the completion of Docker image pushes, or deletions in the container registry. When such an event occurs, Azure Container Registry sends an HTTP POST payload to the webhook's configured URL.

```sql+postgres
select
  name,
  w ->> 'location' as webhook_location,
  w -> 'properties' -> 'actions' as actions,
  w -> 'properties' ->> 'scope' as scope,
  w -> 'properties' ->> 'status' as status
from
  azure_container_registry,
  jsonb_array_elements(webhooks) as w;
```

```sql+sqlite
select
  name,
  json_extract(w.value, '$.location') as webhook_location,
  json_extract(w.value, '$.properties.actions') as actions,
  json_extract(w.value, '$.properties.scope') as scope,
  json_extract(w.value, '$.properties.status') as status
from
  azure_container_registry,
  json_each(webhooks) as w;
```

### List registries not configured with virtual network service endpoint
Determine the areas in which registries are not configured with a virtual network service endpoint. This is useful in identifying potential security risks where network access is allowed without restrictions.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(network_rule_set, '$.defaultAction') as network_rule_default_action,
  json_extract(network_rule_set, '$.virtualNetworkRules') as virtual_network_rules
from
  azure_container_registry
where
  network_rule_set is not null
  and json_extract(network_rule_set, '$.defaultAction') = 'Allow';
```

### List registries with admin user account enabled
Determine the areas in which administrative user accounts are activated within your Azure Container Registries. This is beneficial to ascertain potential security risks and maintain best practices for access control.

```sql+postgres
select
  name,
  admin_user_enabled,
  region
from
  azure_container_registry
where
  admin_user_enabled;
```

```sql+sqlite
select
  name,
  admin_user_enabled,
  region
from
  azure_container_registry
where
  admin_user_enabled;
```