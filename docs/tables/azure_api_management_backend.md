---
title: "Steampipe Table: azure_api_management_backend - Query Azure API Management Backends using SQL"
description: "Allows users to query Azure API Management Backends, specifically providing access to configuration details, service details, and backend settings."
---

# Table: azure_api_management_backend - Query API Management Backend Configurations using SQL

In Azure API Management (APIM), a "Backend" represents the configuration of a web service or API that is the destination for API requests processed by APIM's policies. This includes various types of backends, such as services running on Azure App Service, virtual machines, or external APIs hosted outside of Azure. Understanding the backend configurations is crucial for managing the API request and response lifecycle within APIM.

## Table Usage Guide

The `azure_api_management_backend` table provides insights into the backend configurations within Azure API Management. Use this table to gain detailed information about how API requests are routed and managed, including the protocol used, service details, and security configurations.

## Examples

### Basic Info
Gain a general understanding of the backend configurations in your API Management service. This query is essential for a quick overview of how your APIs are set up in terms of endpoints and protocols.

```sql+postgres
select
  name,
  id,
  protocol,
  service_name,
  url,
  type
from
  azure_api_management_backend;
```

```sql+sqlite
select
  name,
  id,
  protocol,
  service_name,
  url,
  type
from
  azure_api_management_backend;
```

### List the Backend Credential Contract Properties
Review the authorization credentials for backends. This is crucial for understanding and managing security configurations for your API backends.

```sql+postgres
select
  name,
  id,
  credentials -> 'authorization' ->> 'parameter' as parameter,
  credentials -> 'authorization' ->> 'scheme' as scheme
from
  azure_api_management_backend;
```

```sql+sqlite
select
  name,
  id,
  json_extract(json_extract(credentials, '$.authorization'), '$.parameter') as parameter,
  json_extract(json_extract(credentials, '$.authorization'), '$.scheme') as scheme
from
  azure_api_management_backend;
```

### Get the TLS Configuration for a Particular Backend Service
Examine the TLS (Transport Layer Security) configurations for backend services. This query helps ensure that secure communication protocols are in place.

```sql+postgres
select
  name,
  id,
  tls -> 'validateCertificateChain' as tls_validate_certificate_chain,
  tls -> 'validateCertificateName' as tls_validate_certificate_name
from
  azure_api_management_backend;
```

```sql+sqlite
select
  name,
  id,
  json_extract(tls, '$.validateCertificateChain') as tls_validate_certificate_chain,
  json_extract(tls, '$.validateCertificateName') as tls_validate_certificate_name
from
  azure_api_management_backend;
```

### List Backends That Follow HTTP Protocol
Identify backends using the HTTP protocol. This can be important for reviewing API security, as HTTP lacks the encryption of HTTPS.

```sql+postgres
select
  name,
  id,
  protocol,
  service_name,
  url,
  type
from
  azure_api_management_backend
where
  protocol = 'http';
```

```sql+sqlite
select
  name,
  id,
  protocol,
  service_name,
  url,
  type
from
  azure_api_management_backend
where
  protocol = 'http';
```
