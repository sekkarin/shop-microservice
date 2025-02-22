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
  source      = "/etc/vault/secrets-template"
  destination = "/vault/secrets/.env"
}
