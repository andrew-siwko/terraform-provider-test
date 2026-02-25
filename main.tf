# main.tf
terraform {
  required_providers {
    mirror = {
      source = "registry.terraform.io/andrew/property-mirror"
    }
  }
}

resource "mirror_variable" "rhel_test" {
  value = "Running on RHEL 9.7"
}

resource "mirror_variable" "rhel_test_2" {
  value = "Also Running on RHEL 9.7"
}