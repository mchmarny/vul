resource "google_sql_database_instance" "db_instance" {
  database_version    = "POSTGRES_14"
  name                = "${var.name}-instance"
  region              = var.location
  root_password       = var.db_password
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

    database_flags {
      name  = "cloudsql.enable_pgaudit"
      value = "on"
    }

    database_flags {
      name  = "pgaudit.log"
      value = "all"
    }

    ip_configuration {
      ipv4_enabled = true
      require_ssl  = false
    }
  }
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
