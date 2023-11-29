---
title: "Steampipe Table: azure_private_dns_zone - Query Azure Private DNS Zones using SQL"
description: "Allows users to query Azure Private DNS Zones."
---

# Table: azure_private_dns_zone - Query Azure Private DNS Zones using SQL

Azure Private DNS is a service that provides reliable resolution of domain names in a Virtual Network, without the need for you to create and manage custom DNS solution. This service supports Azure services like VMs and Azure Kubernetes Service (AKS) clusters within a virtual network to securely and privately resolve and connect to the services running in the virtual network. It provides name resolution for virtual machines (VMs) within a VNet and between VNets.

## Table Usage Guide

The 'azure_private_dns_zone' table provides insights into Private DNS Zones within Azure DNS. As a DevOps engineer, explore zone-specific details through this table, including record sets, virtual network links, and associated metadata. Utilize it to uncover information about zones, such as those with private DNS records, the virtual networks linked to the zones, and the verification of DNS records. The schema presents a range of attributes of the Private DNS Zone for your analysis, like the zone name, resource group, subscription ID, and associated tags.

## Examples

### Basic info
Explore which private DNS zones are present in your Azure infrastructure, including their associated resource groups and any attached tags. This can help you track and manage your resources more effectively.

```sql
select
  name,
  resource_group,
  tags
from
  azure_private_dns_zone;
```

### List private DNS zones with record sets
Analyze the configuration of your Azure private DNS zones to identify those with more than one record set. This can be useful in pinpointing specific locations where multiple resources might be sharing the same DNS zone.

```sql
select
  name,
  resource_group
from
  azure_private_dns_zone
where
  number_of_record_sets > 1;
```

### List private DNS zones linked to no virtual networks
Explore which private DNS zones in Azure are not linked to any virtual networks. This can help identify potential areas for optimization or detect configuration errors.

```sql
select
  name,
  resource_group
from
  azure_private_dns_zone
where
  number_of_virtual_network_links_with_registration = 0;
```