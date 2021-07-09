resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 3600"
  }
}

resource "null_resource" "named_test_resource" {
  depends_on = [
    null_resource.delay
  ]
  provisioner "local-exec" {
    command = "az network vnet-gateway delete -g {{ resourceName }} -n {{ resourceName }}"
  }
}
