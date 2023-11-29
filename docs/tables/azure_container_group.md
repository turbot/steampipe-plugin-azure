---
title: "Steampipe Table: azure_container_group - Query Azure Container Instances using SQL"
description: "Allows users to query Azure Container Groups."
---

# Table: azure_container_group - Query Azure Container Instances using SQL

Azure Container Instances offers the fastest and simplest way to run a container in Azure, without having to provision any virtual machines and without having to adopt a higher-level service. It is a solution for any scenario that can operate in isolated containers, without orchestration. Run event-driven applications, quickly deploy from your container development pipelines, and run data processing and build jobs.

## Table Usage Guide

The 'azure_container_group' table provides insights into Container Groups within Azure Container Instances. As a DevOps engineer, explore Container Group-specific details through this table, including the containers within the group, the image they are using, the commands they are running, and associated metadata. Utilize it to uncover information about Container Groups, such as their current state, the events that have occurred within them, and the configurations they have been given. The schema presents a range of attributes of the Container Group for your analysis, like the group name, creation date, associated containers, and associated tags.

## Examples

### Basic info
Explore the configuration of your Azure Container Groups to understand their provisioning states and restart policies. This is useful for assessing the performance and management of your resources across different regions.

```sql
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
Uncover the details of encryption for each group within your Azure Container service. This will help you assess the security measures in place and ensure that each group is properly protected.

```sql
select
  name,
  encryption_properties ->> 'VaultBaseURL' as vault_base_url,
  encryption_properties ->> 'KeyName' as key_name,
  encryption_properties ->> 'KeyVersion' as key_version,
  region
from
  azure_container_group;
```

### List groups that have restart policy set to `OnFailure`
Discover the segments that have their restart policy set to 'OnFailure' in the Azure Container Group. This can be useful in assessing system resilience and planning for potential system failures.

```sql
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

### Count groups by operation type
Analyze the distribution of Azure container groups based on their operating system type. This can provide insights into the predominant OS types used within your container groups, aiding in system optimization and resource planning.

```sql
select
  os_type,
  count(name) as group_count
from
  azure_container_group
group by
  os_type;
```

### Get IP address details of each group
Explore which containers in your Azure environment are associated with specific IP addresses. This can help you manage your network configuration and identify potential bottlenecks or security risks.

```sql
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

### Get image registry credential details of each group
Explore the authentication details for image registries used by different container groups. This can be useful to ensure proper security measures are in place and to manage access to your image repositories.

```sql
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

### Get DNS configuration details of each group
Explore the DNS configuration details for each container group in Azure. This can help you understand how your container groups are configured for network communication, aiding in network troubleshooting and optimization.

```sql
select
  name,
  id,
  dns_config -> 'NameServers' as name_servers,
  dns_config ->> 'SearchDomains' as search_domains,
  dns_config ->> 'Options' as options
from
  azure_container_group;
```