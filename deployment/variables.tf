# Description: List of variables which can be provided ar runtime to override the specified defaults 

variable "project_id" {
  description = "GCP Project ID"
  type        = string
  nullable    = false
}

variable "name" {
  description = "Base name to derive everythign else from"
  type        = string
  nullable    = false
}

variable "location" {
  description = "Deployment location"
  type        = string
  nullable    = false
}

variable "git_repo" {
  description = "GitHub Repo"
  type        = string
  nullable    = false
}

variable "domain_name" {
  description = "Domain name"
  type        = string
  nullable    = false
}

variable "app_img_uri" {
  description = "app container image to deploy"
  type        = string
  nullable    = false
}

variable "config_secret_version" {
  description = "secret manager version to use"
  type        = string
  default     = "latest"
}

variable "db_conn_uri" {
  description = "Cloud SQL instance connection URI"
  type        = string
  nullable    = false
}

variable "db_machine_type" {
  description = "Cloud SQL tier"
  type        = string
  default     = "db-custom-2-7680"
}
