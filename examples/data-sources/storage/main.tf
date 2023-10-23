# terraform {
#   required_providers {
#     vpsie = {
#         source = "registry.terraform.local/hashicorp/vpsie"
#     }
#   }
# }

# provider "vpsie" {
#   access_token = "yInm5a6lG3hruO6I-hCkzvwN8nHYOkNx0P-TAqAX5TiaTfAYjKZE-xXeif7QPvN3"
# }

# data "vpsie_storages" "all" {}

# output "storages" {
#     value = data.vpsie_storages.all
# }