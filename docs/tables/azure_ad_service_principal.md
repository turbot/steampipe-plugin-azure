# Table: azure_ad_service_principal

**Deprecated. Use [azuread_service_principal](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_service_principal) instead.**

An Azure service principal is an identity created for use with applications, hosted services, and automated tools to access Azure resources.

## Examples

### List of ad service principals where service principal account is disabled

```sql
select
  object_id,
  object_type,
  display_name,
  account_enabled
from
  azure_ad_service_principal
where
  not account_enabled;
```


### List of ad service principals where app role assignment is not required

```sql
select
  object_id,
  display_name,
  app_role_assignment_required
from
  azure_ad_service_principal
where
  not app_role_assignment_required;
```


### Application role info of service principals

```sql
select
  object_id,
  approle ->> 'allowedMemberTypes' as allowed_member_types,
  approle ->> 'description' as description,
  approle ->> 'displayName' as display_name,
  approle -> 'isEnabled' as isEnabled,
  approle ->> 'id' as id,
  approle ->> 'value' as id
from
  azure_ad_service_principal
  cross join jsonb_array_elements(app_roles) as approle;
```


### Oauth 2.0 permission info of ad service principal

```sql
select
  object_id,
  perm ->> 'adminConsentDescription' as admin_consent_description,
  perm ->> 'adminConsentDisplayName' as admin_consent_display_ame,
  perm ->> 'id' as id,
  perm ->> 'isEnabled' as is_enabled,
  perm ->> 'type' as type,
  perm ->> 'value' as value
from
  azure_ad_service_principal
  cross join jsonb_array_elements(oauth2_permissions) as perm;
```
