select id, akas, title
from azure.azure_role_assignment_for_resource
where name = '{{ output.resource_id.value.split("/").pop() }}' and scope = '/subscriptions/{{ output.subscription_id.value }}'