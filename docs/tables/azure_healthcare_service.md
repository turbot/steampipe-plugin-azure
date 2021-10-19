# Table: azure_healthcare_service

Azure Healthcare APIs provides pipelines that help you manage protected health information (PHI) data at scale. 

## Examples

### Basic info

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

```sql
select
  name,
  id,
  cosmos_db_configuration ->> 'keyVaultKeyUri' as key_vault_key_uri,
  cosmos_db_configuration -> 'offerThroughput' as offer_throughput
from
  azure_healthcare_service;
```
