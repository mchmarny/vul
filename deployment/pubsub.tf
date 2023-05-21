resource "google_pubsub_topic" "worker_topic" {
  name = var.pubsub_worker_queue
}

resource "google_pubsub_subscription" "worker_topic_sub" {
  name  = "${var.pubsub_worker_queue}-sub"
  topic = google_pubsub_topic.worker_topic.name

  ack_deadline_seconds = 600

  retry_policy {
    minimum_backoff = "300s"
  }

  push_config {
    push_endpoint = "${google_cloud_run_service.worker.status[0].url}/api/v1/process"

    attributes = {
      x-goog-version = "v1"
    }

    oidc_token {
      service_account_email = google_service_account.runner_service_account.email
      audience              = "${google_cloud_run_service.worker.status[0].url}/api/v1/process"
    }
  }
}
