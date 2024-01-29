---
title: "Steampipe Table: azure_eventgrid_topic - Query Azure Event Grid Topics using SQL"
description: "Allows users to query Azure Event Grid Topics, providing insights into the event routing service which helps in efficiently and reliably routing events from any source, to any destination, at any scale."
---

# Table: azure_eventgrid_topic - Query Azure Event Grid Topics using SQL

Azure Event Grid is a service within Microsoft Azure that enables the development of event-based applications and simplifies the creation of serverless workflows. It is a fully managed intelligent event routing service that uses a publish-subscribe model for uniform event consumption. Event Grid efficiently and reliably routes events from any source, to any destination, at any scale.

## Table Usage Guide

The `azure_eventgrid_topic` table provides insights into Azure Event Grid Topics within Microsoft Azure. As a developer or system administrator, explore topic-specific details through this table, including event routing details, message retention policy, and associated metadata. Utilize it to uncover information about topics, such as those with specific event types, the routing policies, and the verification of event schemas.

## Examples

### Basic info
Gain insights into the status and details of your Azure EventGrid topics. This query is useful in monitoring the provisioning state and type of each topic, helping ensure smooth operation of your event-driven applications.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state
from
  azure_eventgrid_topic;
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state
from
  azure_eventgrid_topic;
```

### List domains not configured with private endpoint connections
Determine the areas in which domains are not set up with private endpoint connections. This can help in identifying potential security risks and ensuring all domains are properly configured.

```sql+postgres
select
  name,
  id,
  type,
  private_endpoint_connections
from
  azure_eventgrid_topic
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
  azure_eventgrid_topic
where
  private_endpoint_connections is null;
```

### List domains with local authentication disabled
Explore which domains have local authentication disabled to ensure high security. This is useful for identifying potential weak spots in your system's security configuration.

```sql+postgres
select
  name,
  id,
  type,
  disable_local_auth
from
  azure_eventgrid_topic
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
  azure_eventgrid_topic
where
  disable_local_auth = 1;
```