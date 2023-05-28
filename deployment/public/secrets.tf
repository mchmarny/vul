# API Key Secret 
resource "google_secret_manager_secret" "config_secret" {
  secret_id = "${data.template_file.name.rendered}-config"

  labels = {
    label = "config"
  }

  replication {
    automatic = true
  }

  depends_on  = [null_resource.sql_instance_config_update]
}

data "template_file" "config" {
  template = file("../../config/secret-prod.yaml")
}

# API Key Secret version (holds data)
resource "google_secret_manager_secret_version" "config_secret_version" {
  secret      = google_secret_manager_secret.config_secret.name
  secret_data = data.template_file.config.rendered
}