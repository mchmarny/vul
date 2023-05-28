# Description: This file contains the file resources for the deployment

# data.template_file.version.rendered
data "template_file" "version" {
  template = file("../../.version")
}

# data.template_file.name.rendered
data "template_file" "name" {
  template = yamldecode(file("../../config/secret-prod.yaml"))["name"]
}

# data.template_file.project_id.rendered
data "template_file" "project_id" {
  template = yamldecode(file("../../config/secret-prod.yaml"))["project_id"]
}

# data.template_file.location.rendered
data "template_file" "location" {
  template = yamldecode(file("../../config/secret-prod.yaml"))["location"]
}

# data.template_file.db_password.rendered
data "template_file" "db_password" {
  template = yamldecode(file("../../config/secret-prod.yaml"))["store"]["password"]
}

# data.template_file.snyk_token.rendered
data "template_file" "snyk_token" {
  template = yamldecode(file("../../config/secret-prod.yaml"))["import"]["snyk_token"]
}

# data.template_file.app_domain.rendered
data "template_file" "app_domain" {
  template = yamldecode(file("../../config/secret-prod.yaml"))["app"]["domain"]
}
