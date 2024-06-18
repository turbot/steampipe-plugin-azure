---
title: "Steampipe Table: azure_private_endpoint - Query Azure Private Endpoints using SQL"
description: "Allows users to query Private Endpoints in Azure, providing detailed information about each private endpoint, including its associated network interfaces, IP configurations, and connection details."
---

# Table: azure_private_endpoint - Query Azure Private Endpoints using SQL

A Private Endpoint in Azure is a network interface that connects privately and securely to a service powered by Azure Private Link. This enables access to Azure services and resources over a private IP address in a virtual network (VNet), which helps to secure network traffic.

## Table Usage Guide

The `azure_private_endpoint` table provides insights into Private Endpoints within Azure. As an Infrastructure Engineer, explore detailed information about each private endpoint through this table, including its IP configurations, associated network interfaces, and connection details. Use this table to manage and optimize your private endpoint configurations, ensuring secure and efficient communication between your Azure resources.

## Examples

### Basic private endpoint information
Explore the configuration of your Azure private endpoint to gain insights into your private IP address details and associated network interfaces. This can help you understand your endpoint configurations and manage your network resources effectively.

```sql+postgres
select
  name,
  ip ->> 'name' as config_name,
  ip -> 'PrivateEndpointIPConfigurationProperties' ->> 'PrivateIPAddress' as private_ip_address,
  ip -> 'PrivateEndpointIPConfigurationProperties' as private_ip_configuration,
  ip -> 'properties' ->> 'Name' as private_ip_name,
  ip -> 'properties' ->> 'Type' as private_ip_type
from
  azure_private_endpoint
  cross join jsonb_array_elements(ip_configurations) as ip;
```

```sql+sqlite
select
  name,
  json_extract(ip.value, '$.name') as config_name,
  json_extract(ip.value, '$.PrivateEndpointIPConfigurationProperties.PrivateIPAddress') as private_ip_address,
  json_extract(ip.value, '$.PrivateEndpointIPConfigurationProperties') as private_ip_configuration,
  json_extract(ip.value, '$.Type') as private_ip_type
from
  azure_private_endpoint,
  json_each(ip_configurations) as ip;
```

### Find all private endpoints in a specific subnet
Determine the areas in which your Azure private endpoints are allocated within a specific subnet. This is useful for understanding how your network resources are distributed and identifying potential areas of congestion or security vulnerabilities.

```sql+postgres
select
  name,
  ip ->> 'name' as config_name,
  ip -> 'PrivateEndpointIPConfigurationProperties' ->> 'PrivateIPAddress' as private_ip_address
from
  azure_private_endpoint
  cross join jsonb_array_elements(ip_configurations) as ip
where
  ip -> 'PrivateEndpointIPConfigurationProperties' ->> 'PrivateIPAddress' like '10.66.0.%';
```

```sql+sqlite
select
  name,
  json_extract(ip.value, '$.name') as config_name,
  json_extract(ip.value, '$.PrivateEndpointIPConfigurationProperties.PrivateIPAddress') as private_ip_address
from
  azure_private_endpoint,
  json_each(ip_configurations) as ip
where
  json_extract(ip.value, '$.PrivateEndpointIPConfigurationProperties.PrivateIPAddress') like '10.66.0.%';
```

### Application security groups attached to each private endpoint
Explore which application security groups are linked to each private endpoint in your Azure environment. This can help in managing and improving the security posture of your network.

```sql+postgres
select
  name,
  jsonb_array_elements(application_security_groups) ->> 'id' as security_group_id
from
  azure_private_endpoint;
```

```sql+sqlite
select
  name,
  json_extract(security_group.value, '$.id') as security_group_id
from
  azure_private_endpoint,
  json_each(application_security_groups) as security_group;
```

### Custom DNS configurations
List the custom DNS configurations associated with each private endpoint.

```sql+postgres
select
  name,
  jsonb_array_elements(custom_dns_configs) ->> 'Fqdn' as fqdn,
  jsonb_array_elements(custom_dns_configs) ->> 'IPAddresses' as ip_addresses
from
  azure_private_endpoint;
```

```sql+sqlite
select
  name,
  json_extract(dns_config.value, '$.Fqdn') as fqdn,
  json_extract(dns_config.value, '$.IPAddresses') as ip_addresses
from
  azure_private_endpoint,
  json_each(custom_dns_configs) as dns_config;
```

### Extended location information
Retrieve the extended location information for each private endpoint.

```sql+postgres
select
  name,
  extended_location ->> 'name' as extended_location_name,
  extended_location ->> 'type' as extended_location_type
from
  azure_private_endpoint;
```

```sql+sqlite
select
  name,
  json_extract(extended_location, '$.name') as extended_location_name,
  json_extract(extended_location, '$.type') as extended_location_type
from
  azure_private_endpoint;
```
