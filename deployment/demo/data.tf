resource "google_sql_database_instance" "db_instance" {
  database_version = "POSTGRES_14"
  name             = data.template_file.name.rendered
  region           = data.template_file.location.rendered
  root_password    = data.template_file.db_password.rendered

  settings {
    activation_policy     = "ALWAYS"
    availability_type     = "ZONAL"
    tier                  = var.db_machine_type
    disk_autoresize       = true
    disk_autoresize_limit = 0
    disk_size             = 100
    disk_type             = "PD_SSD"

    insights_config {
      query_insights_enabled  = true
      query_string_length     = 1024
      record_application_tags = true
      record_client_address   = true
    }

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
      ipv4_enabled                                  = true
      private_network                               = data.google_compute_network.default.id
      enable_private_path_for_google_cloud_services = true
    }

    user_labels = {
      demo = "s3c"
    }
  }
}

# db
resource "google_sql_database" "db" {
  name     = data.template_file.name.rendered
  instance = google_sql_database_instance.db_instance.name
}

# accounts 
resource "google_sql_user" "github_actions_user" {
  name     = "${data.template_file.name.rendered}-github-sa"
  instance = google_sql_database_instance.db_instance.name
  type     = "CLOUD_IAM_SERVICE_ACCOUNT"
}

resource "google_sql_user" "runner_service_account" {
  name     = "${data.template_file.name.rendered}-run-sa"
  instance = google_sql_database_instance.db_instance.name
  type     = "CLOUD_IAM_SERVICE_ACCOUNT"
}
