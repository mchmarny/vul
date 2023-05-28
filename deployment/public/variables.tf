variable "git_repo" {
  description = "GitHub Repo"
  type        = string
  default     = "mchmarny/vul"
}

variable "config_secret_version" {
  description = "secret manager version to use"
  type        = string
  default     = "latest"
}

variable "db_machine_type" {
  description = "Cloud SQL tier"
  type        = string
  default     = "db-custom-2-7680"
}

variable "queue_image_schedule" {
  description = "Cloud Scheduler cron schedule for queueing images"
  type        = string
  default     = "30 */8 * * *" # every 8 hours
}

variable "img_reg_uri" {
  description = "App image URI"
  type        = string
  default     = "us-west1-docker.pkg.dev/s3cme1/vul"
}

variable "app_ingress" {
  description = "App service ingress"
  type        = string
  default     = "all"
}
