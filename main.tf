# main.tf
terraform {
  required_providers {
    mirror = {
      source = "registry.terraform.io/andrew/property-mirror"
    }
  }
}

resource "mirror_variable" "rhel_test_1" {
  value = "Running on RHEL 9.7"
}

resource "mirror_variable" "rhel_test_2" {
  value = "Also Running on RHEL 9.7"
}

output "var_1" {
  value = mirror_variable.rhel_test_1.value
}

output "var_2" {
  value = mirror_variable.rhel_test_2.value
} 