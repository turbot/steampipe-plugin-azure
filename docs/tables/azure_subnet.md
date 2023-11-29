---
title: "Steampipe Table: azure_subnet - Query Azure Virtual Networks Subnets using SQL"
description: "Allows users to query Azure Virtual Networks Subnets."
---

# Table: azure_subnet - Query Azure Virtual Networks Subnets using SQL

A subnet is a range within a virtual network where Azure resources like VMs and PaaS services can be deployed and accessed from the internet, other networks, and the internet. Subnets help to segment the virtual network into one or more sub-networks, providing a range of IP addresses, network security policies, and route tables. Each subnet in Azure is associated with a route table, which defines the rules for packet routing.

## Table Usage Guide

The 'azure_subnet' table provides insights into subnets within Azure Virtual Networks. As a DevOps engineer, explore subnet-specific details through this table, including IP configurations, network security group details, and associated metadata. Utilize it to uncover information about subnets, such as those with private endpoints, the associated route table, and the service endpoint policies. The schema presents a range of attributes of the subnet for your analysis, like the subnet ID, address prefix, associated network security group, and associated route table.

## Examples

### Virtual network and IP address range of each subnet
Analyze the settings to understand the relationship between your virtual network and IP address range for each subnet. This can help you effectively manage your network resources and ensure optimal performance and security.

```sql
select
  name,
  virtual_network_name,
  address_prefix,
  region,
  resource_group
from
  azure_subnet;
```

### Route table associated with each subnet
Explore which route tables are associated with each subnet in your Azure environment. This can help you understand and manage the routing of network traffic within your virtual network.

```sql
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

### Network security group associated with each subnet
Explore which network security groups are linked to each Azure subnet. This is beneficial for understanding your network's security layout and identifying any potential vulnerabilities or misconfigurations.

```sql
select
  name subnet_name,
  virtual_network_name,
  split_part(network_security_group_id, '/', 9) as network_security_name
from
  azure_subnet;
```

### Service endpoints info of each subnet
Explore which locations are associated with each subnet service in Azure. This can help in understanding the geographical distribution of your services and planning for potential regional expansion or redundancy.

```sql
select
  name,
  endpoint -> 'locations' as location,
  endpoint -> 'service' as service
from
  azure_subnet
  cross join jsonb_array_elements(service_endpoints) as endpoint;
```