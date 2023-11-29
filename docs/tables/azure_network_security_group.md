---
title: "Steampipe Table: azure_network_security_group - Query Azure Network Security Groups using SQL"
description: "Allows users to query Azure Network Security Groups"
---

# Table: azure_network_security_group - Query Azure Network Security Groups using SQL

An Azure Network Security Group is a feature in Microsoft Azure that provides inbound and outbound network traffic filtering for various types of Azure resources. It acts as a virtual firewall, offering a layer of security by enabling you to configure network traffic rules. Network Security Groups can be associated with subnets, network interfaces, or both, providing control over traffic flowing in and out of Azure resources within a virtual network.

## Table Usage Guide

The 'azure_network_security_group' table provides insights into Network Security Groups within Azure Networking. As a Network Administrator, explore specific details through this table, including security rules, default rules, and associated subnets. Utilize it to uncover information about Network Security Groups, such as those with open inbound or outbound rules, the associated resources, and the verification of rule priorities. The schema presents a range of attributes of the Network Security Group for your analysis, like the group name, location, type, and associated tags.

## Examples

### Subnets and network interfaces attached to the network security groups
Discover the segments that are linked to your network security groups by analyzing network interfaces and subnets. This allows you to better understand and assess your Azure network configuration and security posture.

```sql
select
  name,
  split_part(nic ->> 'id', '/', 9) network_interface,
  split_part(vn ->> 'id', '/', 9) virtual_network,
  split_part(vn ->> 'id', '/', 11) subnets
from
  azure_network_security_group
  cross join jsonb_array_elements(network_interfaces) as nic,
  jsonb_array_elements(subnets) as vn;
```

### List the network security groups whose inbound is not restricted from the internet
Determine the network security groups in your Azure environment that have unrestricted inbound access from the internet. This can help you identify potential security risks and take necessary actions to secure your network.

```sql
select
  name,
  sg ->> 'name' as sg_name,
  sg -> 'properties' ->> 'access' as access,
  sg -> 'properties' ->> 'description' as description,
  sg -> 'properties' ->> 'destinationPortRange' as destination_port_range,
  sg -> 'properties' ->> 'direction' as direction,
  sg -> 'properties' ->> 'priority' as priority,
  sg -> 'properties' ->> 'sourcePortRange' as source_port_range,
  sg -> 'properties' ->> 'protocol' as protocol
from
  azure_network_security_group
  cross join jsonb_array_elements(security_rules) as sg
where
  (
    sg -> 'properties' ->> 'sourcePortRange' = '*'
    and sg -> 'properties' ->> 'destinationPortRange' = '*'
    and sg -> 'properties' ->> 'access' = 'Allow'
  );
```

### Default security group rules info
Gain insights into the default security rules of your Azure network security group. This query can help you understand the access, description, direction, priority, and protocol of each rule, which is crucial for maintaining network security and troubleshooting connectivity issues.

```sql
select
  name,
  sg -> 'name' as sg_name,
  sg -> 'properties' ->> 'access' as access,
  sg -> 'properties' ->> 'description' as description,
  sg -> 'properties' ->> 'destinationPortRange' as destination_port_range,
  sg -> 'properties' ->> 'direction' as direction,
  sg -> 'properties' ->> 'priority' as priority,
  sg -> 'properties' ->> 'sourcePortRange' as source_port_range,
  sg -> 'properties' ->> 'protocol' as protocol
from
  azure_network_security_group
  cross join jsonb_array_elements(default_security_rules) as sg;
```