---
title: "Steampipe Table: azure_healthcare_service - Query Azure Healthcare Services using SQL"
description: "Allows users to query Azure Healthcare Services."
---

# Table: azure_healthcare_service - Query Azure Healthcare Services using SQL

Azure Healthcare Service is a managed service that provides built-in support for industry standard health data protocols and data types. It enables health organizations to ingest, manage, and persist health information in the cloud. This service supports the FHIR (Fast Healthcare Interoperability Resources) standard for exchanging healthcare information electronically.

## Table Usage Guide

The 'azure_healthcare_service' table provides insights into Azure Healthcare Services. As a DevOps engineer, explore service-specific details through this table, including the service type, provisioning state, access policies, and associated metadata. Utilize it to uncover information about services, such as those with public network access, the kind of service, and the provisioning state. The schema presents a range of attributes of the Azure Healthcare Service for your analysis, like the service name, resource group, subscription ID, and associated tags.

## Examples

### Basic info
Explore the fundamental characteristics of your Azure healthcare services. This query helps you understand the types of services you have, their authorities, and whether they allow credentials, providing insights into your overall healthcare service configuration.

```sql
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
Explore which healthcare services utilize the 'fhir-R4' type in Azure. This can be useful in identifying and managing services that employ this specific standard.

```sql
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
Explore the status and details of private connections for a healthcare service. This can be useful in managing and securing network connections within a healthcare service infrastructure.

```sql
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

### List diagnostic settings for healthcare service
Analyze the settings to understand the diagnostic configurations for a healthcare service. This is useful for managing and monitoring the health of the service.

```sql
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

### List Cosmos DB configuration settings
Review the configuration for Azure's Cosmos DB to determine the key vault key URI and offer throughput, which can be useful for assessing database performance and security settings.

```sql
select
  name,
  id,
  cosmos_db_configuration ->> 'keyVaultKeyUri' as key_vault_key_uri,
  cosmos_db_configuration -> 'offerThroughput' as offer_throughput
from
  azure_healthcare_service;
```