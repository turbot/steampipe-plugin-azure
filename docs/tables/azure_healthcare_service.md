---
title: "Steampipe Table: azure_healthcare_service - Query Azure Healthcare Services using SQL"
description: "Allows users to query Azure Healthcare Services, providing insights into the health and performance of healthcare services and potential anomalies."
---

# Table: azure_healthcare_service - Query Azure Healthcare Services using SQL

Azure Healthcare Services is a service within Microsoft Azure that allows users to manage and monitor health data in the cloud. It provides a centralized way to set up and manage healthcare services, including data protection, access control, and compliance features. Azure Healthcare Services helps users stay informed about the health and performance of their healthcare services and take appropriate actions when predefined conditions are met.

## Table Usage Guide

The `azure_healthcare_service` table provides insights into Healthcare Services within Microsoft Azure. As a healthcare data analyst, explore service-specific details through this table, including data protection measures, access control settings, and compliance features. Utilize it to uncover information about services, such as those with potential security risks, the access control settings of each service, and the compliance status of each service.

## Examples

### Basic info
Explore the characteristics and settings of your Azure Healthcare Services. This query can be useful for understanding the configuration and type of each service, which is essential for effective management and utilization of these services.

```sql+postgres
select
  name,
  id,
  kind,
  type,
  allow_credentials,
  audience,
  authority
from
  azure_healthcare_service;
```

```sql+sqlite
select
  name,
  id,
  kind,
  type,
  allow_credentials,
  audience,
  authority
from
  azure_healthcare_service;
```

### List healthcare services of fhir-R4 kind 
Determine the areas in which specific healthcare services of 'fhir-R4' kind are utilized within the Azure platform. This can be helpful in assessing the usage and distribution of this particular type of service.

```sql+postgres
select
  name,
  id,
  type,
  kind
from
  azure_healthcare_service
where
  kind = 'fhir-R4';
```

```sql+sqlite
select
  name,
  id,
  type,
  kind
from
  azure_healthcare_service
where
  kind = 'fhir-R4';
```

### List private endpoint connection details for healthcare service
Gain insights into the private connection details of your healthcare service. This query is useful for understanding the connection's state and type, which can assist in troubleshooting or optimizing your service's network configuration.

```sql+postgres
select
  name,
  id,
  p ->> 'PrivateEndpointConnectionId' as private_endpoint_connection_id,
  p ->> 'ProvisioningState' as private_endpoint_provisioning_state,
  p ->> 'PrivateEndpointConnectionName' as private_endpoint_connection_name,
  p ->> 'PrivateEndpointConnectionType' as private_endpoint_connection_type
from
  azure_healthcare_service,
  jsonb_array_elements(private_endpoint_connections) as p;
```

```sql+sqlite
select
  name,
  id,
  json_extract(p.value, '$.PrivateEndpointConnectionId') as private_endpoint_connection_id,
  json_extract(p.value, '$.ProvisioningState') as private_endpoint_provisioning_state,
  json_extract(p.value, '$.PrivateEndpointConnectionName') as private_endpoint_connection_name,
  json_extract(p.value, '$.PrivateEndpointConnectionType') as private_endpoint_connection_type
from
  azure_healthcare_service,
  json_each(private_endpoint_connections) as p;
```

### List diagnostic settings for healthcare service
Explore the diagnostic settings of your healthcare service to gain insights into its configuration and performance. This can be beneficial for identifying potential issues or areas for improvement in your service's setup.

```sql+postgres
select
  name,
  id,
  d ->> 'id' as diagnostic_setting_id,
  d ->> 'name' as diagnostic_setting_name,
  d ->> 'type' as diagnostic_setting_type,
  d ->> 'properties' as diagnostic_setting_properties
from
  azure_healthcare_service,
  jsonb_array_elements(diagnostic_settings) as d;
```

```sql+sqlite
select
  name,
  id,
  json_extract(d.value, '$.id') as diagnostic_setting_id,
  json_extract(d.value, '$.name') as diagnostic_setting_name,
  json_extract(d.value, '$.type') as diagnostic_setting_type,
  json_extract(d.value, '$.properties') as diagnostic_setting_properties
from
  azure_healthcare_service,
  json_each(diagnostic_settings) as d;
```

### List Cosmos DB configuration settings
Gain insights into the configuration settings of your Azure Cosmos DB within the healthcare service. This analysis can help optimize the database performance and security by understanding key vault key URI and offer throughput details.

```sql+postgres
select
  name,
  id,
  cosmos_db_configuration ->> 'keyVaultKeyUri' as key_vault_key_uri,
  cosmos_db_configuration -> 'offerThroughput' as offer_throughput
from
  azure_healthcare_service;
```

```sql+sqlite
select
  name,
  id,
  json_extract(cosmos_db_configuration, '$.keyVaultKeyUri') as key_vault_key_uri,
  json_extract(cosmos_db_configuration, '$.offerThroughput') as offer_throughput
from
  azure_healthcare_service;
```