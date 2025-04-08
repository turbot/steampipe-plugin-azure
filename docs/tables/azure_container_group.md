---
title: "Steampipe Table: azure_container_group - Query Azure Container Groups using SQL"
description: "Allows users to query Azure Container Groups, providing detailed information about each container group's configuration, status, and metadata."
folder: "Containers"
---

# Table: azure_container_group - Query Azure Container Groups using SQL

Azure Container Groups is a service within Microsoft Azure that allows you to manage multiple containers as a single entity. It provides a way to deploy, manage, and scale containers together, simplifying the process of managing multi-container applications. Azure Container Groups helps you to deploy applications quickly and efficiently, without the need to manage the underlying infrastructure.

## Table Usage Guide

The `azure_container_group` table provides insights into Container Groups within Microsoft Azure. As a DevOps engineer, explore group-specific details through this table, including container configurations, statuses, and associated metadata. Utilize it to uncover information about container groups, such as those with specific configurations, the statuses of various container groups, and the verification of metadata.

## Examples

### Basic info
Analyze the settings to understand the configuration of your Azure container groups. This can help in managing and optimizing your resources by identifying the regions, restart policies, and other key details.

```sql+postgres
select
  name,
  id,
  provisioning_state,
  restart_policy,
  sku,
  region
from
  azure_container_group;
```

```sql+sqlite
select
  name,
  id,
  provisioning_state,
  restart_policy,
  sku,
  region
from
  azure_container_group;
```

### Get encryption details of each group
This query helps to analyze the encryption details of each group within your Azure Container service. It is useful for assessing your security setup and ensuring that encryption keys are properly configured and up-to-date across all regions.

```sql+postgres
select
  name,
  encryption_properties ->> 'VaultBaseURL' as vault_base_url,
  encryption_properties ->> 'KeyName' as key_name,
  encryption_properties ->> 'KeyVersion' as key_version,
  region
from
  azure_container_group;
```

```sql+sqlite
select
  name,
  json_extract(encryption_properties, '$.VaultBaseURL') as vault_base_url,
  json_extract(encryption_properties, '$.KeyName') as key_name,
  json_extract(encryption_properties, '$.KeyVersion') as key_version,
  region
from
  azure_container_group;
```

### List groups that have restart policy set to `OnFailure`
Identify the groups in your Azure Container service that have been configured to restart only when a failure occurs. This could be beneficial in managing resources and avoiding unnecessary restarts.

```sql+postgres
select
  name,
  restart_policy,
  provisioning_state,
  type
from
  azure_container_group
where
  restart_policy = "OnFailure";
```

```sql+sqlite
select
  name,
  restart_policy,
  provisioning_state,
  type
from
  azure_container_group
where
  restart_policy = 'OnFailure';
```

### Count groups by operation type
Analyze the distribution of Azure container groups based on their operating system type. This can help in understanding the usage pattern of different OS types within your Azure container groups.

```sql+postgres
select
  os_type,
  count(name) as group_count
from
  azure_container_group
group by
  os_type;
```

```sql+sqlite
select
  os_type,
  count(name) as group_count
from
  azure_container_group
group by
  os_type;
```

### Get IP address details of each group
Discover the segments that provide information about IP addresses associated with each group. This is useful in understanding the network connectivity and accessibility of these groups within the Azure container ecosystem.

```sql+postgres
select
  name,
  ip_address -> 'Ports' as ports,
  ip_address ->> 'Type' as ip_address_type,
  ip_address ->> 'IP' as ip,
  ip_address ->> 'DNSNameLabel' as dns_name_label,
  ip_address ->> 'Fqdn' as fqdn
from
  azure_container_group;
```

```sql+sqlite
select
  name,
  json_extract(ip_address, '$.Ports') as ports,
  json_extract(ip_address, '$.Type') as ip_address_type,
  json_extract(ip_address, '$.IP') as ip,
  json_extract(ip_address, '$.DNSNameLabel') as dns_name_label,
  json_extract(ip_address, '$.Fqdn') as fqdn
from
  azure_container_group;
```

### Get image registry credential details of each group
Determine the credentials of image registries for each container group in Azure. This is useful for managing and verifying access to different image registries.

```sql+postgres
select
  name,
  i ->> 'Server' as server,
  i ->> 'Username' as username,
  i ->> 'Password' as password,
  i ->> 'Identity' as identity,
  i ->> 'IdentityURL' as identity_url
from
  azure_container_group,
  jsonb_array_elements(image_registry_credentials) as i;
```

```sql+sqlite
select
  name,
  json_extract(i.value, '$.Server') as server,
  json_extract(i.value, '$.Username') as username,
  json_extract(i.value, '$.Password') as password,
  json_extract(i.value, '$.Identity') as identity,
  json_extract(i.value, '$.IdentityURL') as identity_url
from
  azure_container_group,
  json_each(image_registry_credentials) as i;
```

### Get DNS configuration details of each group
This query allows you to gain insights into the DNS configuration details for each Azure container group. It's particularly useful for system administrators who need to manage or troubleshoot network settings across multiple container groups.

```sql+postgres
select
  name,
  id,
  dns_config -> 'NameServers' as name_servers,
  dns_config ->> 'SearchDomains' as search_domains,
  dns_config ->> 'Options' as options
from
  azure_container_group;
```

```sql+sqlite
select
  name,
  id,
  json_extract(dns_config, '$.NameServers') as name_servers,
  json_extract(dns_config, '$.SearchDomains') as search_domains,
  json_extract(dns_config, '$.Options') as options
from
  azure_container_group;
```