resource "google_storage_bucket" "state" {
  name          = "${data.template_file.name.rendered}-terraform-state"
  force_destroy = false
  location      = data.template_file.location.rendered
  storage_class = "STANDARD"
  versioning {
    enabled = true
  }
}

resource "google_storage_bucket" "data_dump_bucket" {
  name          = "${data.template_file.name.rendered}-data-dumps"
  location      = data.template_file.location.rendered
  storage_class = "STANDARD"
  force_destroy = true

  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_binding" "data_dump_bucket_sql_binding" {
  bucket = google_storage_bucket.data_dump_bucket.name
  role   = "roles/storage.admin"
  members = [
    "serviceAccount:${google_sql_database_instance.db_instance.service_account_email_address}",
  ]
}
