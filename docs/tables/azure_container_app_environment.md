---
title: "Steampipe Table: azure_container_app_environment - Query Azure Container App managed environments using SQL"
description: "Allows users to query Azure Container App managed environments, providing insights into the networking, workload profiles and log configuration each environment provides to the apps it hosts."
folder: "Container App"
---

# Table: azure_container_app_environment - Query Azure Container App managed environments using SQL

An Azure Container Apps managed environment is the secure boundary around one or more container apps. Apps in the same environment share a virtual network, write logs to the same destination, and draw compute from the same set of workload profiles. The environment owns the default domain apps are published under and the static outbound IP they egress from.

## Table Usage Guide

The `azure_container_app_environment` table provides insights into Container App managed environments within Microsoft Azure. As a platform engineer, explore environment-specific details through this table, including the vnet and ingress configuration, the workload profiles available to apps, zone redundancy, and where app logs are shipped. Use it to review the network exposure of an environment and to see which compute profiles the apps inside it can be pinned to.

## Examples

### Basic info
Explore the managed environments in your subscription with their provisioning state and default domain, to get an overview of where container apps can be deployed.

```sql+postgres
select
  name,
  id,
  provisioning_state,
  default_domain,
  static_ip,
  region
from
  azure_container_app_environment;
```

```sql+sqlite
select
  name,
  id,
  provisioning_state,
  default_domain,
  static_ip,
  region
from
  azure_container_app_environment;
```

### List environments that are not zone-redundant
Identify environments without zone redundancy, since the apps they host cannot survive the loss of an availability zone.

```sql+postgres
select
  name,
  zone_redundant,
  region
from
  azure_container_app_environment
where
  not coalesce(zone_redundant, false);
```

```sql+sqlite
select
  name,
  zone_redundant,
  region
from
  azure_container_app_environment
where
  coalesce(zone_redundant, 0) = 0;
```

### List environments reachable from public networks
Determine which environments allow public traffic, so their exposure can be reviewed against the intended network design.

```sql+postgres
select
  name,
  public_network_access,
  vnet_configuration ->> 'internal' as vnet_internal,
  region
from
  azure_container_app_environment
where
  public_network_access = 'Enabled';
```

```sql+sqlite
select
  name,
  public_network_access,
  json_extract(vnet_configuration, '$.internal') as vnet_internal,
  region
from
  azure_container_app_environment
where
  public_network_access = 'Enabled';
```

### Get the workload profiles available in each environment
Inspect the compute profiles apps can be pinned to, including the consumption profile and any dedicated ones with their node counts.

```sql+postgres
select
  e.name,
  p ->> 'name' as profile_name,
  p ->> 'workloadProfileType' as profile_type,
  p ->> 'minimumCount' as minimum_count,
  p ->> 'maximumCount' as maximum_count
from
  azure_container_app_environment as e,
  jsonb_array_elements(e.workload_profiles) as p;
```

```sql+sqlite
select
  e.name,
  json_extract(p.value, '$.name') as profile_name,
  json_extract(p.value, '$.workloadProfileType') as profile_type,
  json_extract(p.value, '$.minimumCount') as minimum_count,
  json_extract(p.value, '$.maximumCount') as maximum_count
from
  azure_container_app_environment as e,
  json_each(e.workload_profiles) as p;
```

### Get the subnet each environment is injected into
Find the vnet subnet backing an environment, to tie it back to the network stack that owns the address space.

```sql+postgres
select
  name,
  vnet_configuration ->> 'infrastructureSubnetId' as infrastructure_subnet_id,
  vnet_configuration ->> 'internal' as internal,
  infrastructure_resource_group,
  region
from
  azure_container_app_environment;
```

```sql+sqlite
select
  name,
  json_extract(vnet_configuration, '$.infrastructureSubnetId') as infrastructure_subnet_id,
  json_extract(vnet_configuration, '$.internal') as internal,
  infrastructure_resource_group,
  region
from
  azure_container_app_environment;
```

### Get the log destination configured for each environment
Determine where the apps in each environment ship their logs, to confirm they land in the expected workspace.

```sql+postgres
select
  name,
  app_logs_configuration ->> 'destination' as log_destination,
  app_logs_configuration -> 'logAnalyticsConfiguration' ->> 'customerId' as log_analytics_workspace_id,
  region
from
  azure_container_app_environment;
```

```sql+sqlite
select
  name,
  json_extract(app_logs_configuration, '$.destination') as log_destination,
  json_extract(app_logs_configuration, '$.logAnalyticsConfiguration.customerId') as log_analytics_workspace_id,
  region
from
  azure_container_app_environment;
```

### Count the container apps hosted in each environment
Join the apps back to their environment to see how densely each environment is used.

```sql+postgres
select
  e.name as environment_name,
  e.region,
  count(a.id) as container_app_count
from
  azure_container_app_environment as e
  left join azure_container_app as a on lower(a.environment_id) = lower(e.id)
group by
  e.name,
  e.region
order by
  container_app_count desc;
```

```sql+sqlite
select
  e.name as environment_name,
  e.region,
  count(a.id) as container_app_count
from
  azure_container_app_environment as e
  left join azure_container_app as a on lower(a.environment_id) = lower(e.id)
group by
  e.name,
  e.region
order by
  container_app_count desc;
```
