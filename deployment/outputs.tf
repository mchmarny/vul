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

output "WORKER_SERVICE_URL" {
  value = google_cloud_run_service.worker.status[0].url
}

output "WORKER_SERVICE_IMG" {
  value = "${var.worker_img_uri}:${data.template_file.version.rendered}"
}
