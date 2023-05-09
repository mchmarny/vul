# Description: Creates a Google Container Registry for the project

resource "google_container_registry" "registry" {
  provider      = google-beta
  project       = var.project_id
  location      = "US"
}

resource "google_storage_bucket_iam_member" "admin" {
  bucket  = google_container_registry.registry.id
  role    = "roles/storage.objectAdmin"
  member  = "serviceAccount:${google_service_account.github_actions_user.email}"
}

resource "google_storage_bucket_iam_member" "creator" {
  bucket  = google_container_registry.registry.id
  role    = "roles/storage.objectCreator"
  member  = "serviceAccount:${google_service_account.github_actions_user.email}"
}
