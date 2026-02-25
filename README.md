# Simple Terraform Provider
This code implements a terraform provider which stores the tf resource values in the state.


## Notes
* Created and clone repo
* In the repo directory:
  * go mod init github.com/andrew/terraform-provider-test
  * go get github.com/hashicorp/terraform-plugin-framework
  * go get github.com/hashicorp/terraform-plugin-framework/providerserver
  * go get github.com/hashicorp/terraform-plugin-go/tfprotov6
  * go mod tidy
  * go build -o terraform-provider-test.exe
  * setup.ps1 to create terraform.rc
  * setup.bash to create terraform.rc
  