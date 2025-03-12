vault {
  address = "http://vault:8200"
}

auto_auth {
    method {
      type = "approle"
      config = {
        role_id_file_path = "/etc/vault/role_id"
        secret_id_file_path = "/etc/vault/secret_id"
      }
    }
  
    sinks {
      sink {
        type = "file"
  
        config = {
          path = "/etc/vault/tokenca"
        }
      }
    }
}

template {
  source      = "/etc/vault/auth-template"
  destination = "/vault/secrets/auth-prod/secret.yaml"
}
template {
  source      = "/etc/vault/inventory-template"
  destination = "/vault/secrets/inventory-prod/secret.yaml"
}
template {
  source      = "/etc/vault/item-template"
  destination = "/vault/secrets/item-prod/secret.yaml"
}
template {
  source      = "/etc/vault/payment-template"
  destination = "/vault/secrets/payment-prod/secret.yaml"
}
template {
  source      = "/etc/vault/player-template"
  destination = "/vault/secrets/player-prod/secret.yaml"
}
