---
title: "Steampipe Table: azure_eventgrid_topic - Query Azure Event Grid Topics using SQL"
description: "Allows users to query Azure Event Grid Topics."
---

# Table: azure_eventgrid_topic - Query Azure Event Grid Topics using SQL

Azure Event Grid is a service within Microsoft Azure that allows you to build applications with event-based architectures. It provides a centralized way to manage and react to events from various Azure resources, such as Blob Storage, Resource Groups, and Subscriptions. Azure Event Grid helps you stay informed about the status changes and take appropriate actions when certain conditions are met.

## Table Usage Guide

The 'azure_eventgrid_topic' table provides insights into Event Grid Topics within Azure Event Grid. As a DevOps engineer, explore topic-specific details through this table, including endpoint details, provisioning state, and associated metadata. Utilize it to uncover information about topics, such as those with specific endpoint types, the provisioning state of topics, and the verification of endpoint details. The schema presents a range of attributes of the Event Grid Topic for your analysis, like the topic name, id, type, provisioning state, and associated tags.

## Examples

### Basic info
Explore which Azure Event Grid topics are currently active. This can be useful in assessing the state of your event-driven applications and ensuring they are functioning as expected.

```sql
select
  name,
  id,
  type,
  provisioning_state
from
  azure_eventgrid_topic;
```

### List domains not configured with private endpoint connections
Discover the segments that are lacking private endpoint connections within the Azure EventGrid topic. This allows for pinpointing potential security vulnerabilities in your network configuration.

```sql
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
Analyze the settings to understand which domains have local authentication disabled in your Azure EventGrid topic. This can help enhance security by identifying potential vulnerabilities and ensuring appropriate authentication measures are in place.

```sql
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