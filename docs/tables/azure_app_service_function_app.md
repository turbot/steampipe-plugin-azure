# Table: azure_app_service_function_app

A function app is the container that hosts the execution of individual functions.

## Examples

#### List of app functions which accepts HTTP traffic

```sql
select
  name,
  https_only,
  kind,
  location
from
  azure_app_service_function_app
where
  not https_only;
```


### List of all unreserved app function

```sql
select
  name,
  reserved,
  resource_group
from
  azure_app_service_function_app
where
  not reserved;
```


### Outbound IP addresses and possible outbound IP addresses info of each function app

```sql
select
  name,
  outbound_ip_addresses,
  possible_outbound_ip_addresses
from
  azure_app_service_function_app;
```


### List of app functions where client certificate mode is disabled.

```sql
select
  name,
  client_cert_enabled,
  kind,
  location
from
  azure_app_service_function_app
where
  not client_cert_enabled;
```