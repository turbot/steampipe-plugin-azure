resource "null_resource" "named_test_resource" {
  provisioner "local-exec" {
    command = "az disk-access delete -g {{ output.resource_name.value }} -n {{ output.resource_name.value }}"
  }
}
