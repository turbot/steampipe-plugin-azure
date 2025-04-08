---
title: "Steampipe Table: azure_network_security_group - Query Azure Network Security Groups using SQL"
description: "Allows users to query Network Security Groups in Azure, providing insights into the security rules and configurations of the network infrastructure."
folder: "Networking"
---

# Table: azure_network_security_group - Query Azure Network Security Groups using SQL

A Network Security Group in Azure is a security feature that acts as a virtual firewall for your network in Azure, using inbound and outbound rules to allow or deny network traffic to resources. It provides granular access control over network traffic by defining network security rules that allow or deny traffic based on traffic direction, protocol, source address and port, and destination address and port. This is a fundamental layer of security for virtual networks in Azure.

## Table Usage Guide

The `azure_network_security_group` table provides insights into Network Security Groups within Azure. As a security analyst or network administrator, you can explore the details of these groups through this table, including security rules, configurations, and associated metadata. Utilize this table to uncover information about the security posture of your network, such as the rules that are allowing or denying traffic, the protocols used, and the source and destination addresses and ports.

## Examples

### Subnets and network interfaces attached to the network security groups
Explore the relationships between network security groups, their attached network interfaces, and the subnets within the virtual networks. This can help in understanding the network topology and identifying potential security vulnerabilities.

```sql+postgres
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

```sql+sqlite
Error: SQLite does not support split or string_to_array functions.
```

### List the network security groups whose inbound is not restricted from the internet
Explore which network security groups are not restricting inbound access from the internet. This is useful in identifying potential security vulnerabilities within your network infrastructure.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(sg.value, '$.name') as sg_name,
  json_extract(sg.value, '$.properties.access') as access,
  json_extract(sg.value, '$.properties.description') as description,
  json_extract(sg.value, '$.properties.destinationPortRange') as destination_port_range,
  json_extract(sg.value, '$.properties.direction') as direction,
  json_extract(sg.value, '$.properties.priority') as priority,
  json_extract(sg.value, '$.properties.sourcePortRange') as source_port_range,
  json_extract(sg.value, '$.properties.protocol') as protocol
from
  azure_network_security_group,
  json_each(security_rules) as sg
where
  (
    json_extract(sg.value, '$.properties.sourcePortRange') = '*'
    and json_extract(sg.value, '$.properties.destinationPortRange') = '*'
    and json_extract(sg.value, '$.properties.access') = 'Allow'
  );
```

### Default security group rules info
Discover the details of default security group rules within your Azure network security group. This query can help you understand the access, direction, and protocol of each rule, which can be useful for auditing and optimizing your network security settings.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(sg.value, '$.name') as sg_name,
  json_extract(sg.value, '$.properties.access') as access,
  json_extract(sg.value, '$.properties.description') as description,
  json_extract(sg.value, '$.properties.destinationPortRange') as destination_port_range,
  json_extract(sg.value, '$.properties.direction') as direction,
  json_extract(sg.value, '$.properties.priority') as priority,
  json_extract(sg.value, '$.properties.sourcePortRange') as source_port_range,
  json_extract(sg.value, '$.properties.protocol') as protocol
from
  azure_network_security_group,
  json_each(default_security_rules) as sg;
```