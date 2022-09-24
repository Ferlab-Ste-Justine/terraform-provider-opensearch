resource "opensearch_user" "ops" {
  username = "ops"
  password = "password"
  backend_roles = ["all_access"]
}