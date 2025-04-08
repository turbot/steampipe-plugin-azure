---
title: "Steampipe Table: azure_nat_gateway - Query Azure NAT Gateways using SQL"
description: "Allows users to query Azure NAT Gateways, providing insights into the network traffic routing and potential anomalies."
folder: "Networking"
---

# Table: azure_nat_gateway - Query Azure NAT Gateways using SQL

Azure NAT Gateway is a service within Microsoft Azure that simplifies outbound-only internet connectivity for virtual networks. When configured on a subnet, all outbound connectivity uses your specified static public IP addresses. Overcome the challenges of outbound connectivity from your virtual networks with Azure NAT Gateway.

## Table Usage Guide

The `azure_nat_gateway` table provides insights into NAT Gateways within Microsoft Azure. As a Network Administrator, you can explore details about each NAT Gateway, including its configuration, associated resources, and status. Use this table to ensure your network's outbound connectivity is correctly routed and to quickly identify any potential issues.

## Examples

### Basic info
Explore the status and details of your Azure NAT Gateway configurations to understand their current state and type. This is beneficial for auditing and managing your network resources effectively.

```sql+postgres
select
  name,
  id,
  provisioning_state,
  sku_name,
  type
from
  azure_nat_gateway;
```

```sql+sqlite
select
  name,
  id,
  provisioning_state,
  sku_name,
  type
from
  azure_nat_gateway;
```

### List public IP address details for each nat gateway
Identify the public IP address details linked with each network address translation (NAT) gateway. This can help in managing network traffic and understanding the allocation method and IP version of each public IP address.

```sql+postgres
select
  n.name,
  i.ip_address as ip_address,
  i.ip_configuration_id as ip_configuration_id,
  i.public_ip_address_version as public_ip_address_version,
  i.public_ip_allocation_method as public_ip_allocation_method
from
  azure_nat_gateway as n,
  azure_public_ip as i,
  jsonb_array_elements(n.public_ip_addresses) as ip
where
  ip ->> 'id' = i.id;
```

```sql+sqlite
select
  n.name,
  i.ip_address as ip_address,
  i.ip_configuration_id as ip_configuration_id,
  i.public_ip_address_version as public_ip_address_version,
  i.public_ip_allocation_method as public_ip_allocation_method
from
  azure_nat_gateway as n,
  azure_public_ip as i,
  json_each(n.public_ip_addresses) as ip
where
  json_extract(ip.value, '$.id') = i.id;
```

### List subnet details associated with each nat gateway
Explore the connection between NAT gateways and their associated subnets in your Azure environment. This helps in understanding network flow and can assist in troubleshooting connectivity issues.

```sql+postgres
select
  n.name as name,
  s.name as subnet_name,
  s.virtual_network_name as virtual_network_name
from
  azure_nat_gateway as n,
  azure_subnet as s,
  jsonb_array_elements(n.subnets) as sb
where
  sb ->> 'id' = s.id;
```

```sql+sqlite
select
  n.name as name,
  s.name as subnet_name,
  s.virtual_network_name as virtual_network_name
from
  azure_nat_gateway as n,
  azure_subnet as s,
  json_each(n.subnets) as sb
where
  json_extract(sb.value, '$.id') = s.id;
```