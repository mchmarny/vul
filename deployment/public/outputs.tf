# Description: Outputs for the deployment

output "APP_SERVICE_URL" {
  value = google_cloud_run_service.app.status[0].url
}
