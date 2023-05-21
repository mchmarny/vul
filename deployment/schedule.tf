resource "google_cloud_scheduler_job" "image_queue_schedule" {
  name             = "${var.name}-image-queue-schedule"
  description      = "Queues all images for re-processing"
  schedule         = "0 */1 * * *"
  time_zone        = "America/Los_Angeles"
  attempt_deadline = "900s"
  region           = var.location

  retry_config {
    retry_count = 1
  }

  http_target {
    http_method = "POST"
    uri         = "${google_cloud_run_service.worker.status[0].url}/api/v1/queue"

    oidc_token {
      service_account_email = google_service_account.runner_service_account.email
      audience              = "${google_cloud_run_service.worker.status[0].url}/api/v1/queue"
    }
  }
}
