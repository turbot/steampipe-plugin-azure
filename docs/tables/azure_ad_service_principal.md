# Table: azure_ad_service_principal

**Deprecated. Use [azuread_service_principal](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_service_principal) instead.**

An Azure service principal is an identity created for use with applications, hosted services, and automated tools to access Azure resources.

## Examples

### List of ad service principals where service principal account is disabled
Determine which ad service principals have their account disabled in Azure. This is useful for identifying potential inactive or unused resources within your Azure environment.

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
Determine the areas in which ad service principals do not require app role assignments. This is useful to identify potential areas of your Azure AD environment where security could be improved by requiring app role assignments.

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
Identify the roles of service principals within an application to gain insights into their permissions and status. This is useful for understanding who has access to what within your application and ensuring appropriate security measures are in place.

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
This query is useful to gain insights into the permissions related to the OAuth 2.0 protocol for an advertising service principal in Azure. It allows you to understand the consent descriptions, display names, IDs and the status (enabled or not) of these permissions, which is crucial for managing access and maintaining security.

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