---
title: "Steampipe Table: azure_application_security_group - Query Azure Network Security Groups using SQL"
description: "Allows users to query Azure Application Security Groups"
---

# Table: azure_application_security_group - Query Azure Network Security Groups using SQL

An Azure Application Security Group is a network security feature provided by Azure Network Security Groups. It allows users to define fine-grained network policies based on workloads, centralized on applications, instead of explicit IP addresses. Application Security Groups provide a tool to manage the network security policy at scale and increases the manageability of security policies.

## Table Usage Guide

The 'azure_application_security_group' table provides insights into Application Security Groups within Azure Network Security Groups. As a Network Administrator, explore group-specific details through this table, including security rules, associated network interfaces, and associated metadata. Utilize it to uncover information about groups, such as those with certain security rules, the relationships between different groups, and the verification of security policies. The schema presents a range of attributes of the Application Security Group for your analysis, like the resource group, location, type, and associated tags.

## Examples

### Basic info
Discover the segments of your Azure application security groups, such as their names and regions, to better understand their distribution and organization within your resource groups.

```sql
select
  name,
  region,
  resource_group
from
  azure_application_security_group;
```


### List of application security group without application tag key
Explore which Azure Application Security Groups lack the 'application' tag key. This is useful for identifying potential gaps in your tagging strategy, which could impact resource tracking and management.

```sql
select
  name,
  tags
from
  azure_application_security_group
where
  not tags :: JSONB ? 'application';
```