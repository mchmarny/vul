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

# updates SQL instance private IP and SQL instance name in secret-prod.yaml
# this is really only needed for the first deployment
# subsequent runs will still execute but have not impact on the updated prod config file
# Replacements in secret-prod.yaml:
#   PRIVATE_IP
#   SQL_INSTANCE
resource "null_resource" "sql_instance_config_update" {
  provisioner "local-exec" {
    command = "sed -i '' -e 's/PRIVATE_IP/${google_sql_database_instance.db_instance.private_ip_address}/g' -e 's/SQL_INSTANCE/${data.template_file.project_id.rendered}:${data.template_file.location.rendered}:${data.template_file.name.rendered}/g' ../../config/secret-prod.yaml"
  }
  depends_on = [google_sql_database_instance.db_instance]
  triggers = {
    # ensure this resource is runs every time
    always_run = "${timestamp()}"
  }
}
