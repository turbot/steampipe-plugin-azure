---
title: "Steampipe Table: azure_lb - Query Azure Load Balancers using SQL"
description: "Allows users to query Azure Load Balancers, providing detailed information about their configuration, location, and operational status."
---

# Table: azure_lb - Query Azure Load Balancers using SQL

Azure Load Balancer is a high-performance, ultra low-latency Layer 4 load balancing service for all UDP and TCP protocols. It enables you to build highly scalable and highly available applications by providing automatic routing of network traffic to virtual machines. This service also provides health probes to detect the failure of an application on a virtual machine.

## Table Usage Guide

The `azure_lb` table provides insights into Load Balancers within Azure. As a Network Administrator, explore Load Balancer-specific details through this table, including configuration, location, and operational status. Utilize it to uncover information about Load Balancers, such as their health status, associated resources, and traffic routing rules.

## Examples

### Basic info
Explore which Azure Load Balancer resources are currently being provisioned in different regions. This is useful for managing and optimizing geographically distributed resources.

```sql+postgres
select
  id,
  name,
  type,
  provisioning_state,
  etag,
  region
from
  azure_lb;
```

```sql+sqlite
select
  id,
  name,
  type,
  provisioning_state,
  etag,
  region
from
  azure_lb;
```

### List failed load balancers
Determine the areas in which load balancers have failed within your Azure environment. This can aid in troubleshooting and improving the reliability of your network infrastructure.

```sql+postgres
select
  id,
  name,
  type,
  provisioning_state
from
  azure_lb
where
  provisioning_state = 'Failed';
```

```sql+sqlite
select
  id,
  name,
  type,
  provisioning_state
from
  azure_lb
where
  provisioning_state = 'Failed';
```