cat <<EOF > ~/.terraformrc
provider_installation {
  dev_overrides {
    "andrew/property-mirror" = "$(pwd)"
  }
  direct {}
}
EOF

go build -o terraform-provider-property-mirror