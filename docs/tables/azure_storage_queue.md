---
title: "Steampipe Table: azure_storage_queue - Query Azure Storage Queues using SQL"
description: "Allows users to query Azure Storage Queues, which provide reliable messaging for workflow processing and for communication between components of cloud services."
---

# Table: azure_storage_queue - Query Azure Storage Queues using SQL

Azure Storage Queues offer a simple way for components of a distributed application to communicate asynchronously. They are a part of Azure's scalable and secure cloud storage solution, providing reliable messaging for workflow processing and for communication between components of cloud services. Azure Storage Queues support a set of advanced messaging features, making them ideal for building flexible and reliable applications.

## Table Usage Guide

The 'azure_storage_queue' table provides insights into Azure Storage Queues within Azure's cloud storage solution. As a DevOps engineer, you can explore queue-specific details through this table, including metadata, approximate message count, and associated storage account information. Utilize it to uncover information about your queues, such as their message retention period, visibility timeout, and whether they are enabled for logging or not. The schema presents a range of attributes of the storage queue for your analysis, like the queue name, resource group, and associated tags.

## Examples

### List of queues and their corresponding storage accounts
This query allows you to identify the storage accounts associated with each queue in your Azure environment and their geographical locations. It can be used to manage and organize resources more effectively by understanding where data is stored and how it is distributed across different regions.

```sql
select
  name as queue_name,
  storage_account_name,
  region
from
  azure_storage_queue;
```


### List of storage queues without owner tag key
Discover the segments that lack an 'owner' tag within your Azure application security groups. This could be useful for identifying potential security gaps or for maintaining consistent tagging practices.

```sql
select
  name,
  tags
from
  azure_application_security_group
where
  not tags :: JSONB ? 'owner';
```