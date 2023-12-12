---
title: "Steampipe Table: azure_eventgrid_domain - Query Azure Event Grid Domains using SQL"
description: "Allows users to query Azure Event Grid Domains, specifically providing details about the domain name, resource group, location, input schema, metric resource id, and other related data."
---

# Table: azure_eventgrid_domain - Query Azure Event Grid Domains using SQL

Azure Event Grid Domain is a management tool within Microsoft Azure that allows you to route events from your apps and services to specific handlers. It provides a centralized way to manage and route events that occur within your applications, including virtual machines, databases, web applications, and more. Azure Event Grid Domain helps you stay informed about the events occurring in your Azure resources and take appropriate actions when certain conditions are met.

## Table Usage Guide

The `azure_eventgrid_domain` table provides insights into Event Grid Domains within Microsoft Azure. As a DevOps engineer, explore domain-specific details through this table, including domain name, resource group, location, input schema, and metric resource id. Utilize it to uncover information about the events routing, such as the domain's endpoint, the input schema of the domain, and the provisioning state of the domain.

## Examples

### Basic info
Discover the segments that have been provisioned within your Azure EventGrid domain. This query is useful for gaining insights into the current state of your domain, including identifying the type and status of each segment.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state
from
  azure_eventgrid_domain;
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state
from
  azure_eventgrid_domain;
```

### List domains not configured with private endpoint connections
Identify instances where Azure EventGrid domains are not configured with private endpoint connections. This can be useful for pinpointing potential security gaps in your network infrastructure.

```sql+postgres
select
  name,
  id,
  type,
  private_endpoint_connections
from
  azure_eventgrid_domain
where
  private_endpoint_connections is null;
```

```sql+sqlite
select
  name,
  id,
  type,
  private_endpoint_connections
from
  azure_eventgrid_domain
where
  private_endpoint_connections is null;
```

### List domains with local authentication disabled
Identify instances where domains have local authentication disabled within Azure's event grid. This can be useful to assess potential security risks and ensure compliance with security policies.

```sql+postgres
select
  name,
  id,
  type,
  disable_local_auth
from
  azure_eventgrid_domain
where
  disable_local_auth;
```

```sql+sqlite
select
  name,
  id,
  type,
  disable_local_auth
from
  azure_eventgrid_domain
where
  disable_local_auth = 1;
```