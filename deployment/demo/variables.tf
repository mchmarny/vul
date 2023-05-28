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

variable "rate_limit_threshold_count_per_min" {
  description = "Request rate limit threshold per minute"
  type        = number
  default     = 100
}

variable "app_ingress" {
  description = "App service ingress"
  type        = string
  default     = "internal-and-cloud-load-balancing"
  # all, internal, internal-and-cloud-load-balancing
}
