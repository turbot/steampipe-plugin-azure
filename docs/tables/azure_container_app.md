---
title: "Steampipe Table: azure_container_app - Query Azure Container Apps using SQL"
description: "Allows users to query Azure Container Apps, providing insights into the running status, ingress, scaling and revision configuration of each app."
folder: "Container App"
---

# Table: azure_container_app - Query Azure Container Apps using SQL

Azure Container Apps is a serverless platform for running containerized applications and microservices without managing the underlying Kubernetes infrastructure. Each container app runs inside a managed environment, scales on HTTP traffic, events or KEDA rules, and can be reached through an internal or external ingress endpoint.

## Table Usage Guide

The `azure_container_app` table provides insights into container apps within Microsoft Azure. As a DevOps or platform engineer, explore app-specific details through this table, including the environment the app runs in, its running and provisioning state, its ingress endpoint and its replica bounds. Use it to find apps exposed to the internet, apps that can scale to zero, and apps pinned to a particular workload profile.

## Examples

### Basic info
Explore the container apps in your subscription with the environment they run in and their current running status, to get an overview of your serverless container estate.

```sql+postgres
select
  name,
  id,
  running_status,
  provisioning_state,
  environment_id,
  region
from
  azure_container_app;
```

```sql+sqlite
select
  name,
  id,
  running_status,
  provisioning_state,
  environment_id,
  region
from
  azure_container_app;
```

### List container apps with an external ingress endpoint
Identify container apps that expose a public HTTP endpoint, so their reachability from the internet can be reviewed.

```sql+postgres
select
  name,
  ingress_fqdn,
  ingress_target_port,
  ingress_transport,
  region
from
  azure_container_app
where
  ingress_external;
```

```sql+sqlite
select
  name,
  ingress_fqdn,
  ingress_target_port,
  ingress_transport,
  region
from
  azure_container_app
where
  ingress_external = 1;
```

### List container apps that can scale to zero
Determine which apps are allowed to have no running replica, since those trade cold starts for a lower bill.

```sql+postgres
select
  name,
  min_replicas,
  max_replicas,
  region
from
  azure_container_app
where
  min_replicas = 0;
```

```sql+sqlite
select
  name,
  min_replicas,
  max_replicas,
  region
from
  azure_container_app
where
  min_replicas = 0;
```

### Get the container images each app runs
Inspect the image and compute request of every container in an app's template, to audit what is actually deployed.

```sql+postgres
select
  a.name,
  c ->> 'name' as container_name,
  c ->> 'image' as image,
  c -> 'resources' ->> 'cpu' as cpu,
  c -> 'resources' ->> 'memory' as memory
from
  azure_container_app as a,
  jsonb_array_elements(a.template -> 'containers') as c;
```

```sql+sqlite
select
  a.name,
  json_extract(c.value, '$.name') as container_name,
  json_extract(c.value, '$.image') as image,
  json_extract(c.value, '$.resources.cpu') as cpu,
  json_extract(c.value, '$.resources.memory') as memory
from
  azure_container_app as a,
  json_each(json_extract(a.template, '$.containers')) as c;
```

### List container apps running multiple active revisions
Find apps in multiple-revision mode, where traffic can be split across revisions and a stale revision may still be serving.

```sql+postgres
select
  name,
  active_revisions_mode,
  latest_revision_name,
  latest_ready_revision_name,
  region
from
  azure_container_app
where
  active_revisions_mode = 'Multiple';
```

```sql+sqlite
select
  name,
  active_revisions_mode,
  latest_revision_name,
  latest_ready_revision_name,
  region
from
  azure_container_app
where
  active_revisions_mode = 'Multiple';
```

### Get the environment each container app runs in
Join each app to its managed environment to see the default domain and workload profiles behind it.

```sql+postgres
select
  a.name as app_name,
  e.name as environment_name,
  e.default_domain,
  a.workload_profile_name,
  a.region
from
  azure_container_app as a
  left join azure_container_app_environment as e on lower(e.id) = lower(a.environment_id);
```

```sql+sqlite
select
  a.name as app_name,
  e.name as environment_name,
  e.default_domain,
  a.workload_profile_name,
  a.region
from
  azure_container_app as a
  left join azure_container_app_environment as e on lower(e.id) = lower(a.environment_id);
```

### List container apps pulling from a private registry
Determine which apps authenticate to a private container registry, and how they authenticate to it.

```sql+postgres
select
  a.name,
  r ->> 'server' as registry_server,
  r ->> 'identity' as registry_identity,
  r ->> 'username' as registry_username
from
  azure_container_app as a,
  jsonb_array_elements(a.configuration -> 'registries') as r;
```

```sql+sqlite
select
  a.name,
  json_extract(r.value, '$.server') as registry_server,
  json_extract(r.value, '$.identity') as registry_identity,
  json_extract(r.value, '$.username') as registry_username
from
  azure_container_app as a,
  json_each(json_extract(a.configuration, '$.registries')) as r;
```
