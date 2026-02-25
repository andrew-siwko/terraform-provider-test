# Create the file in your user profile directory
$rcPath = "$env:USERPROFILE\terraform.rc"
$rcContent = @"
provider_installation {
  dev_overrides {
    "registry.terraform.io/andrew/property-mirror" = "C:/Users/asiwk/OneDrive/python/terraform-provider-test"
  }
  direct {}
}
"@

$rcContent | Out-File -FilePath $rcPath -Encoding ascii