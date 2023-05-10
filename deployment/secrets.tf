# API Key Secret 
resource "google_secret_manager_secret" "config_secret" {
  secret_id = "${var.name}-config"

  labels = {
    label = "config"
  }

  replication {
    automatic = true
  }

  depends_on = [
    google_project_service.default["secretmanager.googleapis.com"],
  ]
}

data "template_file" "config" {
  template = file("../config/secret-prod.yaml")
}

# API Key Secret version (holds data)
resource "google_secret_manager_secret_version" "config_secret_version" {
  secret = google_secret_manager_secret.config_secret.name

  secret_data = data.template_file.config.rendered
}
