---
title: "Steampipe Table: azure_maintenance_configuration - Query Azure Maintenance Configurations using SQL"
description: "Enables users to query Azure Maintenance Configurations to manage and control maintenance operations across Azure resources."
---

# Table: azure_maintenance_configuration - Query Azure Maintenance Configurations using SQL

Azure Maintenance Configurations provide a mechanism to manage and control maintenance operations for Azure resources. These configurations offer insights into the scheduled maintenance windows, allowing for better planning and minimal disruption to services.

## Table Usage Guide

The `azure_maintenance_configuration` table allows system administrators, DevOps engineers, and Azure resource managers to delve into maintenance configurations within their Azure environments. It offers a detailed view of maintenance schedules, scope, and metadata associated with each configuration. This table can be particularly useful for planning maintenance windows and understanding the impact on resources.

## Examples

### Basic info
Retrieve essential information about maintenance configurations, including IDs, names, and the scope of maintenance. This query is useful for getting an overview of maintenance activities and their reach within your Azure environment.

```sql+postgres
select
  id,
  name,
  created_at,
  created_by,
  created_by_type,
  last_modified_at,
  last_modified_by,
  last_modified_by_type,
  visibility
from
  azure_maintenance_configuration;
```

```sql+sqlite
select
  id,
  name,
  created_at,
  created_by,
  created_by_type,
  last_modified_at,
  last_modified_by,
  last_modified_by_type,
  visibility
from
  azure_maintenance_configuration;
```

### List configurations that are publicly visible
Understanding which maintenance configurations are public can help in assessing security and compliance postures. Organizations often need to ensure that only the appropriate configurations are public to prevent unintentional exposure of sensitive operational details.

```sql+postgres
select
  id,
  name,
  created_at,
  created_by,
  created_by_type,
  visibility
from
  azure_maintenance_configuration
where
  visibility = 'VisibilityPublic';
```

```sql+sqlite
select
  id,
  name,
  created_at,
  created_by,
  created_by_type,
  visibility
from
  azure_maintenance_configuration
where
  visibility = 'VisibilityPublic';
```

### Maintenance window specifics
Understand specific maintenance windows, including their schedules and any custom properties defined. This query helps in precise planning and adjustments to minimize impact.

```sql+postgres
select
  name,
  window ->> 'StartDateTime' as start_date_time,
  window ->> 'ExpirationDateTime' as expiration_date_time,
  window ->> 'Duration' as duration,
  window ->> 'TimeZone' as time_zone,
  window ->> 'RecurEvery' as recur_every
from
  azure_maintenance_configuration;
```

```sql+sqlite
select
  name,
  json_extract(window, '$.StartDateTime') as start_date_time,
  json_extract(window, '$.ExpirationDateTime') as expiration_date_time,
  json_extract(window, '$.Duration') as duration,
  json_extract(window, '$.TimeZone') as time_zone,
  json_extract(window, '$.RecurEvery') as recur_every
from
  azure_maintenance_configuration;
```