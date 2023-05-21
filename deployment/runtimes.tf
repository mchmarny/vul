locals {
  # List of roles that will be assigned to the runner service account
  runner_roles = toset([
    "roles/cloudsql.client",
    "roles/iam.serviceAccountTokenCreator",
    "roles/monitoring.metricWriter",
    "roles/pubsub.publisher",
  ])
}

# Service Account under which the Cloud Run services will run
resource "google_service_account" "runner_service_account" {
  account_id   = "${var.name}-run-sa"
  display_name = "${var.name}-run-sa"
}

# Policy to allow access to secrets 
data "google_iam_policy" "secret_reader" {
  binding {
    role = "roles/secretmanager.secretAccessor"

    members = [
      "serviceAccount:${google_service_account.runner_service_account.email}",
    ]
  }
}

# Binding of the secret access policy to the service account under which 
# Cloud Run services is running
resource "google_secret_manager_secret_iam_policy" "config_secret_access_policy" {
  project     = var.project_id
  secret_id   = google_secret_manager_secret.config_secret.secret_id
  policy_data = data.google_iam_policy.secret_reader.policy_data
}

# Role binding to allow publisher to publish images
resource "google_project_iam_member" "runner_role_bindings" {
  for_each = local.runner_roles
  project  = var.project_id
  role     = each.value
  member   = "serviceAccount:${google_service_account.runner_service_account.email}"
}

# App Cloud Run service 
resource "google_cloud_run_service" "app" {
  name                       = "${var.name}-app"
  location                   = var.location
  project                    = var.project_id
  autogenerate_revision_name = true

  template {
    spec {
      containers {
        image = "${var.app_img_uri}:${data.template_file.version.rendered}"
        startup_probe {
          http_get {
            path = "/health"
          }
        }
        volume_mounts {
          name       = "config-secret"
          mount_path = "/secrets"
        }
        ports {
          name           = "http1"
          container_port = 8080
        }
        resources {
          limits = {
            cpu    = "2000m"
            memory = "512Mi"
          }
        }
        env {
          name  = "ADDRESS"
          value = ":8080"
        }
        env {
          name  = "CONFIG"
          value = "/secrets/${var.name}"
        }
      }
      volumes {
        name = "config-secret"
        secret {
          secret_name = google_secret_manager_secret.config_secret.secret_id
          items {
            key  = var.config_secret_version
            path = var.name
          }
        }
      }

      container_concurrency = 80
      timeout_seconds       = 300
      service_account_name  = google_service_account.runner_service_account.email
    }
    metadata {
      annotations = {
        "autoscaling.knative.dev/minScale"         = "1"
        "autoscaling.knative.dev/maxScale"         = "10"
        "run.googleapis.com/cloudsql-instances"    = "${var.db_conn_uri}"
        "run.googleapis.com/execution-environment" = "gen2"
      }
    }
  }

  metadata {
    annotations = {
      "run.googleapis.com/client-name" = "terraform"
      "run.googleapis.com/ingress"     = "all"
      # all, internal, internal-and-cloud-load-balancing
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_member" "app-public-access" {
  location = google_cloud_run_service.app.location
  project  = google_cloud_run_service.app.project
  service  = google_cloud_run_service.app.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}



# Worker Cloud Run service 
resource "google_cloud_run_service" "worker" {
  name                       = "${var.name}-worker"
  location                   = var.location
  project                    = var.project_id
  autogenerate_revision_name = true

  template {
    spec {
      containers {
        image = "${var.worker_img_uri}:${data.template_file.version.rendered}"
        startup_probe {
          http_get {
            path = "/health"
          }
        }
        volume_mounts {
          name       = "config-secret"
          mount_path = "/secrets"
        }
        ports {
          name           = "http1"
          container_port = 8080
        }
        resources {
          limits = {
            cpu    = "2000m"
            memory = "2048Mi"
          }
        }
        env {
          name  = "ADDRESS"
          value = ":8080"
        }
        env {
          name  = "CONFIG"
          value = "/secrets/${var.name}"
        }
      }
      volumes {
        name = "config-secret"
        secret {
          secret_name = google_secret_manager_secret.config_secret.secret_id
          items {
            key  = var.config_secret_version
            path = var.name
          }
        }
      }

      container_concurrency = 1
      timeout_seconds       = 600
      service_account_name  = google_service_account.runner_service_account.email
    }
    metadata {
      annotations = {
        "autoscaling.knative.dev/minScale"         = "0"
        "autoscaling.knative.dev/maxScale"         = "100"
        "run.googleapis.com/cloudsql-instances"    = "${var.db_conn_uri}"
        "run.googleapis.com/execution-environment" = "gen2"
      }
    }
  }

  metadata {
    annotations = {
      "run.googleapis.com/client-name" = "terraform"
      "run.googleapis.com/ingress"     = "internal-and-cloud-load-balancing"
      # all, internal, internal-and-cloud-load-balancing
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_member" "worker_service_run_invoker" {
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.runner_service_account.email}"
  location = google_cloud_run_service.worker.location
  project  = google_cloud_run_service.worker.project
  service  = google_cloud_run_service.worker.name
}