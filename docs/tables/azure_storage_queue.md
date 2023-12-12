---
title: "Steampipe Table: azure_storage_queue - Query Azure Storage Queues using SQL"
description: "Allows users to query Azure Storage Queues, specifically to obtain detailed information about the queues, including metadata, message count, and status."
---

# Table: azure_storage_queue - Query Azure Storage Queues using SQL

Azure Storage Queues is a service in Microsoft Azure that provides reliable messaging between and within services. It allows for asynchronous message queuing between application components, whether they are running in the cloud, on the desktop, on-premises, or on mobile devices. Azure Storage Queues simplifies the development of large-scale distributed applications, providing a loosely coupled architecture for improved scalability and reliability.

## Table Usage Guide

The `azure_storage_queue` table provides insights into Azure Storage Queues within Microsoft Azure. As a developer or system administrator, you can explore queue-specific details through this table, including metadata, message count, and status. Utilize it to uncover information about queues, such as those with high message counts, and to monitor the status of queues for improved scalability and reliability.

## Examples

### List of queues and their corresponding storage accounts
Explore which Azure storage queues are linked to specific storage accounts and understand their geographical distribution. This can help in managing resources and optimizing storage strategies.

```sql+postgres
select
  name as queue_name,
  storage_account_name,
  region
from
  azure_storage_queue;
```

```sql+sqlite
select
  name as queue_name,
  storage_account_name,
  region
from
  azure_storage_queue;
```

### List of storage queues without owner tag key
Determine the areas in which Azure application security groups lack an 'owner' tag key. This helps to identify resources that may not be properly managed or tracked.

```sql+postgres
select
  name,
  tags
from
  azure_application_security_group
where
  not tags :: JSONB ? 'owner';
```

```sql+sqlite
select
  name,
  tags
from
  azure_application_security_group
where
  json_extract(tags, '$.owner') is null;
```