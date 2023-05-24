# Description: Outputs for the deployment

output "SA_EMAIL" {
  value       = google_service_account.github_actions_user.email
  description = "Service account to use in GitHub Actions."
}

output "PROVIDER_ID" {
  value       = google_iam_workload_identity_pool_provider.github_provider.name
  description = "Provider ID to use in Auth Actions."
}

output "REGISTRY_URI" {
  value       = "${google_artifact_registry_repository.registry.location}-docker.pkg.dev/${data.google_project.project.name}/${google_artifact_registry_repository.registry.name}"
  description = "Artifact Registry location."
}

output "APP_SERVICE_URL" {
  value = google_cloud_run_service.app.status[0].url
}

output "APP_SERVICE_IMG" {
  value = "${var.app_img_uri}:${data.template_file.version.rendered}"
}

output "SECRET_VERSION" {
  value = "${google_secret_manager_secret.config_secret.secret_id}:latest"
}

output "DB_IP" {
  value = google_sql_database_instance.db_instance.ip_address.0.ip_address
}

output "DB_NAME" {
  value = "postgres://${google_sql_database_instance.db_instance.ip_address.0.ip_address}:5432/${var.name}?sslmode=verify-ca&sslrootcert=ca.pem&sslcert=cert.pem&sslkey=key.pem"
}

