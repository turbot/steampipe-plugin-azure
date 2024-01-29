---
title: "Steampipe Table: azure_subnet - Query Azure Subnets using SQL"
description: "Allows users to query Azure Subnets, providing detailed information about each subnet within Azure Virtual Networks."
---

# Table: azure_subnet - Query Azure Subnets using SQL

Azure Subnets are subdivisions of Azure Virtual Networks, which provide a range of IP addresses that can be used by resources. They allow for the segmentation of networks within Azure, which can enhance security and traffic management. Subnets can be associated with Network Security Groups and Route Tables to further customize network traffic rules.

## Table Usage Guide

The `azure_subnet` table provides insights into Azure Subnets within Azure Virtual Networks. As a network administrator, you can explore subnet-specific details through this table, including associated Network Security Groups, Route Tables, and IP configurations. Utilize it to manage and monitor your network segmentation, ensuring optimal security and traffic flow within your Azure environment.

## Examples

### Virtual network and IP address range of each subnet
Determine the areas in which your Azure virtual networks are deployed and gain insights into the IP address range of each subnet. This can help in managing network configurations and ensuring optimal resource allocation across different regions.

```sql+postgres
select
  name,
  virtual_network_name,
  address_prefix,
  resource_group
from
  azure_subnet;
```

```sql+sqlite
select
  name,
  virtual_network_name,
  address_prefix,
  resource_group
from
  azure_subnet;
```

### Route table associated with each subnet
Determine the areas in which subnets and their associated route tables exist in Azure. This information can be useful to understand the routing of network traffic within your Azure environment.

```sql+postgres
select
  st.name subnet_name,
  st.virtual_network_name,
  rt.name route_table_name,
  jsonb_array_elements(rt.routes) -> 'properties' ->> 'addressPrefix' as route_address_prefix,
  jsonb_array_elements(rt.routes) -> 'properties' ->> 'nextHopType' as route_next_hop_type
from
  azure_route_table as rt
  join azure_subnet st on rt.id = st.route_table_id;
```

```sql+sqlite
select
  st.name as subnet_name,
  st.virtual_network_name,
  rt.name as route_table_name,
  json_extract(route.value, '$.properties.addressPrefix') as route_address_prefix,
  json_extract(route.value, '$.properties.nextHopType') as route_next_hop_type
from
  azure_route_table as rt,
  json_each(rt.routes) as route
join
  azure_subnet as st on rt.id = st.route_table_id;
```

### Network security group associated with each subnet
Explore the association between each subnet and its network security group to understand how your Azure network's security is structured. This can help identify potential vulnerabilities or areas for improvement in your network's security configuration.

```sql+postgres
select
  name subnet_name,
  virtual_network_name,
  split_part(network_security_group_id, '/', 9) as network_security_name
from
  azure_subnet;
```

```sql+sqlite
Error: SQLite does not support split_part function.
```

### Service endpoints info of each subnet
Analyze the settings to understand the service endpoints for each subnet within your Azure environment. This can be useful to identify which services are accessible in specific locations, helping to manage network security and connectivity.

```sql+postgres
select
  name,
  endpoint -> 'locations' as location,
  endpoint -> 'service' as service
from
  azure_subnet
  cross join jsonb_array_elements(service_endpoints) as endpoint;
```

```sql+sqlite
select
  name,
  json_extract(endpoint.value, '$.locations') as location,
  json_extract(endpoint.value, '$.service') as service
from
  azure_subnet,
  json_each(service_endpoints) as endpoint;
```