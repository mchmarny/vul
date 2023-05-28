# Description: Creates a Google Artifact Registry for the project

# Artifact Registry
resource "google_artifact_registry_repository" "registry" {
  provider      = google-beta
  project       = data.template_file.project_id.rendered
  location      = data.template_file.location.rendered
  repository_id = data.template_file.name.rendered
  format        = "DOCKER"
}

# Role binding to allow publisher to publish images
resource "google_artifact_registry_repository_iam_member" "registry_role_binding" {
  provider   = google-beta
  project    = data.template_file.project_id.rendered
  location   = data.template_file.location.rendered
  repository = google_artifact_registry_repository.registry.name
  role       = "roles/artifactregistry.repoAdmin"
  member     = "serviceAccount:${google_service_account.github_actions_user.email}"
}
