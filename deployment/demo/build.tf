# PubSub topic for image queue
resource "google_pubsub_topic" "image_queue_topic" {
  name = "${data.template_file.name.rendered}-image-queue"
}

# Build Pool
resource "google_cloudbuild_worker_pool" "pool" {
  name     = "${data.template_file.name.rendered}-import-pool"
  location = data.template_file.location.rendered
  worker_config {
    disk_size_gb   = 100
    machine_type   = "e2-standard-2"
    no_external_ip = false
  }
  network_config {
    peered_network = data.google_compute_network.default.id
  }
}

# Build Trigger
resource "google_cloudbuild_trigger" "process_image_trigger" {
  location = data.template_file.location.rendered
  name     = "${data.template_file.name.rendered}-process-image"
  filename = "workflow/process.yaml"

  pubsub_config {
    topic = google_pubsub_topic.image_queue_topic.id
  }

  source_to_build {
    uri       = "https://www.github.com/${var.git_repo}"
    ref       = "refs/heads/main"
    repo_type = "GITHUB"
  }

  substitutions = {
    _IMPORT_IMAGE_VERSION = data.template_file.version.rendered
    _DIGEST               = "$(body.message.data)"
    _SNYK_TOKEN           = data.template_file.snyk_token.rendered
    _POOL_ID              = google_cloudbuild_worker_pool.pool.id
  }
}

resource "google_cloudbuild_trigger" "queue_image_trigger" {
  location = data.template_file.location.rendered
  name     = "${data.template_file.name.rendered}-queue-image"
  filename = "workflow/queue.yaml"

  source_to_build {
    uri       = "https://www.github.com/${var.git_repo}"
    ref       = "refs/heads/main"
    repo_type = "GITHUB"
  }

  substitutions = {
    _IMAGE_QUEUE_NAME = "${data.template_file.name.rendered}-image-queue"
    _POOL_ID          = google_cloudbuild_worker_pool.pool.id
  }
}

# Scheduler
resource "google_service_account" "schedule_account" {
  account_id = "${data.template_file.name.rendered}-schedule-account"
}

resource "google_project_iam_binding" "project" {
  project = data.template_file.project_id.rendered
  role    = "roles/cloudbuild.builds.editor"
  members = [
    "serviceAccount:${google_service_account.github_actions_user.email}",
  ]
}

resource "google_cloud_scheduler_job" "queue_image_trigger_schedule" {
  name             = "${data.template_file.name.rendered}-queue-image-schedule"
  schedule         = var.queue_image_schedule
  time_zone        = "America/Los_Angeles"
  attempt_deadline = "900s"
  region           = data.template_file.location.rendered

  retry_config {
    retry_count = 1
  }

  http_target {
    http_method = "POST"
    uri         = "https://cloudbuild.googleapis.com/v1/projects/${data.template_file.project_id.rendered}/locations/${data.template_file.location.rendered}/triggers/${google_cloudbuild_trigger.queue_image_trigger.trigger_id}:run"

    oauth_token {
      service_account_email = google_service_account.github_actions_user.email
    }
  }
}
