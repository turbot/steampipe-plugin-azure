# Table: azure_app_service_web_app

Azure App Service is an HTTP-based service for hosting web applications, REST APIs, and mobile back ends.

## Examples

### Outbound IP addresses and possible outbound IP addresses info of each web app

```sql
select
  name,
  outbound_ip_addresses,
  possible_outbound_ip_addresses
from
  azure_app_service_web_app;
```


### List of web app which accepts HTPP traffics (i.e HTTPS only is disabled)

```sql
select
  name,
  https_only,
  kind,
  location
from
  azure_app_service_web_app
where
  not https_only;
```


### List of web app where client certificate mode is disabled.

```sql
select
  name,
  client_cert_enabled,
  kind,
  location
from
  azure_app_service_web_app
where
  not client_cert_enabled;
```


### Host names of each web app

```sql
select
  name,
  host_names,
  kind,
  location,
  resource_group
from
  azure_app_service_web_app;
```
