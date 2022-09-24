resource "opensearch_user" "dev" {
  username = "dev"
  password = "password"
}

resource "opensearch_user" "product" {
  username = "product"
  password = "password"
}

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

resource "opensearch_role" "staging" {
  name = "staging"

  index_permissions {
      index_patterns = ["staging*"]
      allowed_actions = ["read"]
  }
}

resource "opensearch_role_mapping" "qa" {
  role = "qa"
  users = ["dev", "product"]
  hosts = ["*"]
}

resource "opensearch_role_mapping" "staging" {
  role = "staging"
  users = ["product"]
  hosts = ["*"]
}