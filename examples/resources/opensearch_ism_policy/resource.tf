resource "opensearch_ism_policy" "demo" {
  policy_id = "demo"
  description = "demo"
  default_state = "start"

  states {
    name = "start"

    transitions {
      state_name = "bloated"
      conditions {
        min_size = "5gb"
      }
    }

    transitions {
      state_name = "important"
      conditions {
        min_doc_count = 200
      }
    }

    transitions {
      state_name = "dead"
      conditions {
        min_index_age = "30d"
      }
    }
  }

  states {
    name = "important"

    actions {
      timeout = "5m"
      retry {
        count =   5
        backoff = "exponential"
        delay =   "1m"
      }

      action = "index_priority"
      index_priority = 500
    }

    actions {
      timeout = "5m"
      retry {
        count =   5
        backoff = "exponential"
        delay =   "1m"
      }

      action = "replica_count"
      replica_count = 4
    }

    transitions {
      state_name = "bloated"
      conditions {
        min_size = "5gb"
      }
    }

    transitions {
      state_name = "dead"
      conditions {
        min_index_age = "30d"
      }
    }
  }

  states {
    name = "bloated"

    actions {
      timeout = "5m"
      retry {
        count =   5
        backoff = "exponential"
        delay =   "1m"
      }

      action = "close"
    }

    transitions {
      state_name = "dead"
      conditions {
        min_index_age = "30d"
      }
    }
  }

  states {
    name = "dead"

    actions {
      timeout = "5m"
      retry {
        count =   5
        backoff = "exponential"
        delay =   "1m"
      }

      action = "delete"
    }
  }

  ism_template {
    index_patterns = ["qa*"]
    priority = 100
  }

  ism_template {
    index_patterns = ["staging*"]
    priority = 100
  }
}