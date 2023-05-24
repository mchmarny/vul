resource "random_password" "root_password" {
  length  = 32
  special = true
}

resource "google_sql_database_instance" "db_instance" {
  database_version    = "POSTGRES_14"
  name                = "${var.name}-instance"
  region              = var.location
  root_password       = random_password.root_password.result
  deletion_protection = "true"

  settings {
    activation_policy     = "ALWAYS"
    availability_type     = "ZONAL"
    tier                  = var.db_machine_type
    disk_autoresize       = true
    disk_autoresize_limit = 0
    disk_size             = 100
    disk_type             = "PD_SSD"

    database_flags {
      name  = "cloudsql.iam_authentication"
      value = "on"
    }

    ip_configuration {
      authorized_networks {
        value = "0.0.0.0/0"
      }
      ipv4_enabled = true
      require_ssl  = true
    }
  }
}

# certs 

data "google_sql_ca_certs" "ca_certs" {
  instance = google_sql_database_instance.db_instance.name
}

resource "google_sql_ssl_cert" "client_cert" {
  common_name = var.name
  instance    = google_sql_database_instance.db_instance.name
}

resource "local_file" "private_key" {
    content  = google_sql_ssl_cert.client_cert.private_key
    filename = "../keys/private_key.key"
}

resource "local_file" "public_key" {
    content  = google_sql_ssl_cert.client_cert.cert
    filename = "../keys/public_cert.pem"
}

resource "local_file" "server_ca" {
    content  = google_sql_ssl_cert.client_cert.server_ca_cert
    filename = "../keys/server_ca.pem"
}

# db
resource "google_sql_database" "db" {
  name     = var.name
  instance = google_sql_database_instance.db_instance.name
}

# accounts 
resource "google_sql_user" "github_actions_user" {
  name     = "${var.name}-github-actions-user"
  instance = google_sql_database_instance.db_instance.name
  type     = "CLOUD_IAM_SERVICE_ACCOUNT"
}

resource "google_sql_user" "runner_service_account" {
  name     = "${var.name}-run-sa"
  instance = google_sql_database_instance.db_instance.name
  type     = "CLOUD_IAM_SERVICE_ACCOUNT"
}
