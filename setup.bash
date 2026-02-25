cat <<EOF > ~/.terraformrc
provider_installation {
  dev_overrides {
    "registry.terraform.io/andrew/property-mirror" = "$(pwd)"
  }
  direct {}
}
EOF