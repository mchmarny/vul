terraform {
 backend "gcs" {
   bucket  = "vul-terraform-state"
   prefix  = "demo"
 }
}