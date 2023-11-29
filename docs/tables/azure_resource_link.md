---
title: "Steampipe Table: azure_resource_link - Query Azure Resource Links using SQL"
description: "Allows users to query Azure Resource Links."
---

# Table: azure_resource_link - Query Azure Resource Links using SQL

Azure Resource Links are a feature within Microsoft Azure that allows you to link resources across different resource groups and even across different subscriptions. This feature provides a way to visualize and manage the dependencies between resources, which can be helpful for tasks like application mapping and audit. It also enables you to set up and manage links for various Azure resources, including virtual machines, databases, web applications, and more.

## Table Usage Guide

The 'azure_resource_link' table provides insights into Resource Links within Microsoft Azure. As a DevOps engineer, explore link-specific details through this table, including the source and target of each link, as well as associated metadata. Utilize it to uncover information about the relationships between different resources, such as those spanning across different resource groups or subscriptions. The schema presents a range of attributes of the Resource Link for your analysis, like the link id, source id, target id, and notes.

## Examples

### Basic Info
Explore the connections between different Azure resources. This can be useful in understanding the structure of your Azure environment and identifying potential dependencies or bottlenecks.

```sql
select
  name,
  id,
  type,
  source_id,
  target_id
from
  azure_resource_link;
```

### List resource links with virtual machines
Explore the connections between various resources and virtual machines within your Azure environment. This query can be useful to understand the relationships and dependencies in your infrastructure, providing valuable insights for resource management and optimization.

```sql
select
  name,
  id,
  source_id,
  target_id
from
  azure_resource_link
where
  source_id LIKE '%virtualmachines%';
```