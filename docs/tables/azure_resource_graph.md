---
title: "Steampipe Table: azure_resource_graph - Query Azure Resources using KQL via SQL"
description: "Allows users to execute Azure Resource Graph KQL queries and retrieve results as SQL rows."
folder: "Resource Manager"
---

# Table: azure_resource_graph - Query Azure Resources using KQL via SQL

Azure Resource Graph is a service that allows efficient and performant resource exploration across your Azure subscriptions using Kusto Query Language (KQL). This table executes a KQL query against the Azure Resource Graph API and returns the results as SQL rows.

## Table Usage Guide

The `azure_resource_graph` table allows you to run arbitrary KQL queries against Azure Resource Graph. The `query` column is **required** — you must provide a `WHERE query = '...'` clause in every query.

Each row returned corresponds to one row from the KQL result set. The table provides typed columns for common Azure resource fields extracted from the KQL result.

**Important notes:**
- The `query` column is required. Omitting `WHERE query = '...'` will result in an error.
- Columns are **null** when the KQL query does not project the corresponding fields (e.g., aggregation queries) or when the resource does not have that property.
- The `subscription_id` column reads the `subscriptionId` field from the row when available, falling back to the session subscription ID.
- The `resource_group` column reads the `resourceGroup` field from the row when available, falling back to extraction from the `id` field.

## Columns

| Column | Type | Description |
|--------|------|-------------|
| `query` | `text` | The KQL query executed against Azure Resource Graph. **Required.** |
| `id` | `text` | The resource ID, if projected by the query. |
| `name` | `text` | The resource name, if projected by the query. |
| `type` | `text` | The resource type, if projected by the query. |
| `kind` | `text` | The kind of the resource, if available. |
| `identity` | `jsonb` | The managed identity info of the resource, if available. |
| `managed_by` | `text` | The ID of the resource that manages this resource, if available. |
| `plan` | `jsonb` | The plan info of the resource, if available. |
| `properties` | `jsonb` | The resource properties as returned by the graph query. |
| `sku` | `jsonb` | The SKU of the resource, if available. |
| `tenant_id` | `text` | The tenant ID of the resource, if available. |
| `zones` | `jsonb` | The availability zones of the resource, if available. |
| `extended_location` | `jsonb` | The extended location info of the resource, if available. |
| `tags` | `jsonb` | A map of tags for the resource. |
| `region` | `text` | The Azure region/location in which the resource is located. |
| `resource_group` | `text` | The resource group which holds this resource. |
| `title` | `text` | Title of the resource. |
| `akas` | `jsonb` | Array of globally unique identifier strings (also known as) for the resource. |
| `subscription_id` | `text` | The Azure Subscription ID in which the resource is located. |
| `cloud_environment` | `text` | The Azure Cloud Environment. |

## Examples

### List OS packages pending update

```sql+postgres
select
  name,
  properties ->> 'patchName' as patch_name,
  properties ->> 'version' as version,
  properties ->> 'kbId' as kb_id,
  properties ->> 'classifications' as classifications,
  properties ->> 'rebootBehavior' as reboot_behavior
from
  azure_resource_graph
where
  query = 'patchassessmentresources | where type == "microsoft.compute/virtualmachines/patchassessmentresults/softwarepatches"';
```

```sql+sqlite
select
  name,
  properties ->> 'patchName' as patch_name,
  properties ->> 'version' as version,
  properties ->> 'kbId' as kb_id,
  properties ->> 'classifications' as classifications,
  properties ->> 'rebootBehavior' as reboot_behavior
from
  azure_resource_graph
where
  query = 'patchassessmentresources | where type == "microsoft.compute/virtualmachines/patchassessmentresults/softwarepatches"';
```

### List updated OS packages

```sql+postgres
select
  name,
  properties ->> 'patchName' as patch_name,
  properties ->> 'version' as version,
  properties ->> 'patchInstallationState' as installation_state,
  properties ->> 'classifications' as classifications
from
  azure_resource_graph
where
  query = 'patchassessmentresources | where type == "microsoft.compute/virtualmachines/patchassessmentresults/softwarepatches"';
```

```sql+sqlite
select
  name,
  properties ->> 'patchName' as patch_name,
  properties ->> 'version' as version,
  properties ->> 'patchInstallationState' as installation_state,
  properties ->> 'classifications' as classifications
from
  azure_resource_graph
where
  query = 'patchassessmentresources | where type == "microsoft.compute/virtualmachines/patchassessmentresults/softwarepatches"';
```

### Find Azure Arc-enabled servers

```sql+postgres
select
  name,
  id,
  resource_group,
  subscription_id,
  kind,
  properties ->> 'osType' as os_type,
  properties ->> 'osVersion' as os_version,
  properties ->> 'osSku' as os_sku,
  properties ->> 'status' as status,
  properties ->> 'agentVersion' as agent_version
from
  azure_resource_graph
where
  query = 'Resources | where type == "microsoft.hybridcompute/machines" | project name, id, kind, properties';
```
