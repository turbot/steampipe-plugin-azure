---
title: "Steampipe Table: azure_express_route_circuit - Query Azure ExpressRoute Circuits using SQL"
description: "Allows users to query Azure ExpressRoute Circuits"
---

# Table: azure_express_route_circuit - Query Azure ExpressRoute Circuits using SQL

Azure ExpressRoute is a cloud integration solution for creating private data connections between your on-premises infrastructure and Microsoft Azure. These connections do not go over the public Internet, providing higher security, reliability, and speeds with lower latencies than typical connections over the Internet. ExpressRoute connections are ideal for data migration, replication for business continuity, disaster recovery, and other high-availability strategies.

## Table Usage Guide

The 'azure_express_route_circuit' table provides insights into ExpressRoute Circuits within Azure Networking. As a network engineer, explore circuit-specific details through this table, including peering locations, service provider details, and associated metadata. Utilize it to uncover information about circuits, such as those with high bandwidth usage, the peering relationships between circuits, and the verification of service key. The schema presents a range of attributes of the ExpressRoute Circuit for your analysis, like the circuit ARN, creation date, attached peering locations, and associated tags.

## Examples

### Basic info
Explore which Azure Express Route Circuits allow classic operations and analyze their provisioning states to understand their current status and configuration. This can be useful for identifying any circuits that may require updates or changes.

```sql
select
  name,
  id,
  allow_classic_operations,
  circuit_provisioning_state
from
  azure_express_route_circuit;
```

### List express route circuits with global reach enabled
Analyze the settings to understand which Azure Express Route Circuits have global reach enabled. This can be useful to determine the areas in which your network traffic can extend globally, thus optimizing your network strategy.

```sql
select
  name,
  sku_tier,
  sku_name
from
  azure_express_route_circuit
where
  global_reach_enabled;
```

### List premium express route circuits
Discover the segments that are using premium tier Express Route Circuits in Azure. This can be beneficial for assessing the distribution of resources and optimizing cost management within your cloud infrastructure.

```sql
select
  name,
  sku_tier,
  sku_name
from
  azure_express_route_circuit
where
  sku_tier = 'Premium';
```