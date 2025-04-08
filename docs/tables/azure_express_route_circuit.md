---
title: "Steampipe Table: azure_express_route_circuit - Query Azure Express Route Circuits using SQL"
description: "Allows users to query Azure Express Route Circuits, providing detailed information about each circuit's configuration, status, and performance."
folder: "Networking"
---

# Table: azure_express_route_circuit - Query Azure Express Route Circuits using SQL

Azure Express Route Circuits is a dedicated connectivity option in Microsoft Azure that enables users to create private connections between Azure datacenters and infrastructure on their premises or in a colocation environment. It provides a more reliable, faster, and lower-latency network connection than typical internet-based connections. Express Route Circuits are highly beneficial for bandwidth-heavy tasks, data migration, and secure connectivity needs.

## Table Usage Guide

The `azure_express_route_circuit` table provides insights into Express Route Circuits within Microsoft Azure. As a Network Administrator, explore circuit-specific details through this table, including peering information, service provider details, and bandwidth. Utilize it to monitor the performance and status of each Express Route Circuit, ensuring optimal connectivity and performance for your Azure resources.

## Examples

### Basic info
Explore the status and settings of your Azure Express Route Circuits to understand their operational capabilities and provisioning state. This can assist in managing and optimizing your network connectivity.

```sql+postgres
select
  name,
  id,
  allow_classic_operations,
  circuit_provisioning_state
from
  azure_express_route_circuit;
```

```sql+sqlite
select
  name,
  id,
  allow_classic_operations,
  circuit_provisioning_state
from
  azure_express_route_circuit;
```

### List express route circuits with global reach enabled
Explore the express route circuits in your Azure environment that have global reach enabled. This is useful for assessing the scale of your network connectivity and understanding the associated costs.

```sql+postgres
select
  name,
  sku_tier,
  sku_name
from
  azure_express_route_circuit
where
  global_reach_enabled;
```

```sql+sqlite
select
  name,
  sku_tier,
  sku_name
from
  azure_express_route_circuit
where
  global_reach_enabled = 1;
```

### List premium express route circuits
Explore which express route circuits in your Azure environment are categorized as 'Premium'. This can be useful for understanding your network infrastructure and identifying areas for potential cost optimization.

```sql+postgres
select
  name,
  sku_tier,
  sku_name
from
  azure_express_route_circuit
where
  sku_tier = 'Premium';
```

```sql+sqlite
select
  name,
  sku_tier,
  sku_name
from
  azure_express_route_circuit
where
  sku_tier = 'Premium';
```