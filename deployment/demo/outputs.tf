# Description: Outputs for the deployment

output "PROVIDER_ID" {
  value       = google_iam_workload_identity_pool_provider.github_provider.name
}

output "PROVIDER_EMAIL" {
  value       = google_service_account.github_actions_user.email
}

output "REGISTRY_URI" {
  value       = "${google_artifact_registry_repository.registry.location}-docker.pkg.dev"
}

output "REGISTRY_FOLDER" {
  value       = "${data.google_project.project.name}/${google_artifact_registry_repository.registry.name}"
}

output "APP_SERVICE_DOMAIN" {
  value = data.template_file.app_domain.rendered
}

output "APP_SERVICE_LB_IP" {
  value = module.lb-http.external_ip
}

output "APP_SERVICE_IMG" {
  value = "${data.template_file.location.rendered}-docker.pkg.dev/${data.template_file.project_id.rendered}/${data.template_file.name.rendered}/app:${data.template_file.version.rendered}"
}

output "DB_IP" {
  value = google_sql_database_instance.db_instance.private_ip_address
}

output "DB_INSTANCE" {
  value = "${data.template_file.project_id.rendered}:${data.template_file.location.rendered}:${data.template_file.name.rendered}"
}

