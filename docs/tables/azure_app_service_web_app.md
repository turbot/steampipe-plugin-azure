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


### List web apps which accepts HTTP traffics (i.e HTTPS only is disabled)

```sql
select
  name,
  https_only,
  kind,
  region
from
  azure_app_service_web_app
where
  not https_only;
```


### List of web app where client certificate mode is disabled

```sql
select
  name,
  client_cert_enabled,
  kind,
  region
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
  region,
  resource_group
from
  azure_app_service_web_app;
```


### List web apps with latest HTTP version

```sql
select
  name,
  enabled,
  region
from
  azure_app_service_web_app
where
  (configuration -> 'properties' ->> 'http20Enabled')::boolean;
```


### List web apps that have FTP deployments set to disabled

```sql
select
  name,
  configuration -> 'properties' ->> 'ftpsState' as ftps_state
from
  azure_app_service_web_app
where
  configuration -> 'properties' ->> 'ftpsState' <> 'AllAllowed';
```


### List web apps that have managed service identity disabled

```sql
select
  name,
  enabled,
  region,
  identity
from
  azure_app_service_web_app
where
  identity = '{}';
```
