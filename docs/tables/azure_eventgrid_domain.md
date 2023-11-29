---
title: "Steampipe Table: azure_eventgrid_domain - Query Azure Event Grid Domains using SQL"
description: "Allows users to query Azure Event Grid Domains"
---

# Table: azure_eventgrid_domain - Query Azure Event Grid Domains using SQL

Azure Event Grid Domain is an Azure service that simplifies event routing and delivery from source to destination. It is a management and organization layer for event publishing, allowing you to route events from many sources to many destinations. Azure Event Grid Domains provide a single service for managing routing of events from various sources, all with the same security and authentication model.

## Table Usage Guide

The 'azure_eventgrid_domain' table provides insights into Event Grid Domains within Azure Event Grid. As a DevOps engineer, explore domain-specific details through this table, including endpoint, provision state, and associated metadata. Utilize it to uncover information about domains, such as those with specific input schema, the provisioning state, and the endpoint. The schema presents a range of attributes of the Event Grid Domain for your analysis, like the domain name, resource group, and associated tags.

## Examples

### Basic info
Explore the status and types of your EventGrid Domains in Azure. This can help you manage and organize your resources effectively.

```sql
select
  name,
  id,
  type,
  provisioning_state
from
  azure_eventgrid_domain;
```

### List domains not configured with private endpoint connections
Uncover the details of domains lacking private endpoint connections within the Azure EventGrid. This query is useful for identifying potential security vulnerabilities and ensuring proper configuration for secure data transmission.

```sql
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
This query helps identify domains where local authentication has been disabled, providing a quick way to review security settings and ensure proper access control measures are in place. This can be particularly useful in large-scale environments where manual review would be time-consuming.

```sql
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