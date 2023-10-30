# Table: azure_api_management_backend

The "Backend" in API Management represents the configuration of the web service/API on the backend that the gateway will forward the request to, after the request has been accepted and processed by APIM's policies. An API Management Backend provides details about the destination of the API requests. Backends in APIM can be, for instance, a web service running on Azure App Service, a virtual machine, or even an external API hosted outside of Azure. They play a critical role in the end-to-end API request and response lifecycle within APIM.

## Examples

### Basic Info

```sql
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

### List the backend credential contract properties

```sql
select
  name,
  id,
  credentials -> 'authorization' ->> 'parameter' as parameter,
  credentials -> 'authorization' ->> 'scheme' as scheme
from
  azure_api_management_backend;
```

### Get the TLS configuration for a particular backend service

```sql
select
  name,
  id,
  tls -> 'validateCertificateChain' as tls_validate_certificate_chain,
  tls -> 'validateCertificateName' as tls_validate_certificate_name,
from
  azure_api_management_backend;
```

### List backends that follow http protocol

```sql
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