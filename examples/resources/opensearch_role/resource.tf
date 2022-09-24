resource "opensearch_role" "qa" {
  name = "qa"

  index_permissions {
      index_patterns = ["qa*"]
      allowed_actions = ["read", "write"]
  }

  cluster_permissions = [
    "indices_monitor"
  ]
}