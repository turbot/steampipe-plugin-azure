---
title: "Steampipe Table: azure_resource_graph - Query Azure Resource Graph using SQL"
description: "Allows users to execute Azure Resource Graph KQL queries and retrieve results as SQL rows."
folder: "Resource Graph"
---

# Table: azure_resource_graph - Query Azure Resource Graph using SQL

Azure Resource Graph is a service that allows efficient and performant resource exploration across your Azure subscriptions using Kusto Query Language (KQL). This table executes a KQL query against the Azure Resource Graph API and returns the results as SQL rows.

## Table Usage Guide

The `azure_resource_graph` table allows you to run arbitrary KQL queries against Azure Resource Graph. The `query` column is **required** — you must provide a `WHERE query = '...'` clause in every query.

Each row returned corresponds to one row from the KQL result set. The table provides typed columns for common Azure resource fields extracted from the KQL result.

**Important notes:**
- The `query` column is required. Omitting `WHERE query = '...'` will result in an error.
- Columns are **null** when the KQL query does not project the corresponding fields (e.g., aggregation queries) or when the resource does not have that property.
- If you are only querying for Azure resources, consider using the azure_resource table instead; this table is better suited for aggregation, projection, and join queries via KQL.


## Examples

### Basic info
List the first 10 resources in your Azure subscriptions.

```sql+postgres
select 
    id, 
    name, 
    type, 
    region 
from 
    azure_resource_graph 
where 
    query = 'Resources | limit 10';
```

```sql+sqlite
select 
    id, 
    name, 
    type, 
    region 
from 
    azure_resource_graph 
where 
    query = 'Resources | limit 10';
```

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
  query = 'patchassessmentresources | where type == "microsoft.compute/virtualmachines/patchinstallationresults/softwarepatches"';
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
  query = 'patchassessmentresources | where type == "microsoft.compute/virtualmachines/patchinstallationresults/softwarepatches"';
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

```sql+sqlite
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
