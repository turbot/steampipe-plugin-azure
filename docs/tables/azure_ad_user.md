# Table: azure_ad_user

**Deprecated. Use [azuread_user](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_user) instead.**

Azure Active Directory (Azure AD) is Microsoft's cloud-based identity and access management service, which helps employees sign in and access resources.

## Examples

### Basic active directory user info
Explore user details within your Azure Active Directory to gain insights into their status and contact information. This can be particularly useful for managing user access and maintaining up-to-date records.

```sql
select
  display_name,
  user_principal_name,
  given_name,
  mail,
  account_enabled,
  object_id
from
  azure_ad_user;
```


### List of guest users in the active directory
Determine the areas in which guest users are active within your directory. This can help in managing user access and maintaining security protocols.

```sql
select
  display_name,
  user_principal_name,
  mail,
  user_type,
  usage_location
from
  azure_ad_user
where
  user_type = 'Guest';
```


### Password profile info of each user
Determine the areas in which user password policies are enforced and where users are required to change their passwords at their next login. This helps to understand and manage user security within your Azure Active Directory.

```sql
select
  display_name,
  user_principal_name,
  additional_properties -> 'passwordProfile' -> 'enforceChangePasswordPolicy' as enforce_change_password_policy,
  additional_properties -> 'passwordProfile' -> 'forceChangePasswordNextLogin' as change_password_next_login
from
  azure_ad_user;
```