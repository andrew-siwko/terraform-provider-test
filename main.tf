# main.tf
terraform {
  required_providers {
    hashivar = {
      source = "registry.terraform.io/andrew/property-mirror"
    }
  }
}

resource "hashivar_variable" "rhel_test" {
  value = "Running on RHEL 9"
}