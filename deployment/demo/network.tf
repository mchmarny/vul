# External IP to use in SSL cert on LB
resource "google_compute_global_address" "http_lb_address" {
  project = data.template_file.project_id.rendered
  name    = "${data.template_file.name.rendered}-ip"

  address_type = "EXTERNAL"
  ip_version   = "IPV4"

  lifecycle {
    prevent_destroy = true
  }
}

# Serverless Network Endpoint Groups (NEGs) for HTTPS LB to CLoud Run services
module "lb-http" {
  source  = "GoogleCloudPlatform/lb-http/google//modules/serverless_negs"
  version = "9.0.0"

  project = data.template_file.project_id.rendered
  name    = "${data.template_file.name.rendered}-lb-http"

  create_address = false
  address        = google_compute_global_address.http_lb_address.address

  ssl                             = true
  managed_ssl_certificate_domains = [data.template_file.app_domain.rendered]
  https_redirect                  = true

  backends = {
    default = {
      description             = null
      enable_cdn              = true
      security_policy         = google_compute_security_policy.policy.name
      edge_security_policy    = null
      custom_request_headers  = null
      custom_response_headers = null
      compression_mode        = null
      protocol                = "HTTPS"
      port_name               = "http"
      compression_mode        = "DISABLED"
      connection_draining_timeout_sec = 300

      groups = [
        {
          group = google_compute_region_network_endpoint_group.serverless_neg.id
        }
      ]

      iap_config = {
        enable               = false
        oauth2_client_id     = null
        oauth2_client_secret = null
      }

      log_config = {
        enable      = true
        sample_rate = 1.0
      }

      cdn_policy = {
        cacheMode                    = "CACHE_ALL_STATIC"
        maxTtl                       = 86400
        signed_url_cache_max_age_sec = 3600
        cache_key_policy = {
          include_protocol     = true
          include_query_string = true
          include_host         = true
        }
      }
    }
  }

  depends_on = [
    google_project_service.compute_engine,
    google_compute_global_address.http_lb_address,
  ]
}

# Region network endpoint group for Cloud Run sercice in that region
resource "google_compute_region_network_endpoint_group" "serverless_neg" {
  name                  = "${data.template_file.name.rendered}-neg"
  network_endpoint_type = "SERVERLESS"
  region                = data.template_file.location.rendered

  cloud_run {
    service = google_cloud_run_service.app.name
  }
}


# Cloud Armor policies 
resource "google_compute_security_policy" "policy" {
  name     = "${data.template_file.name.rendered}-security-policy"
  provider = google-beta

  rule {
    action      = "deny(403)"
    description = "owasp-crs-v030001 protocolattack"
    priority    = 901
    match {
      expr {
        expression = "evaluatePreconfiguredExpr('protocolattack-canary')"
      }
    }
  }

  rule {
    action      = "deny(403)"
    description = "owasp-crs-v030001 sessionfixation"
    priority    = 902
    match {
      expr {
        expression = "evaluatePreconfiguredExpr('sessionfixation-canary')"
      }
    }
  }

  rule {
    action      = "deny(403)"
    description = "owasp-crs-v030001 scannerdetection"
    priority    = 903
    match {
      expr {
        expression = "evaluatePreconfiguredExpr('scannerdetection-canary')"
      }
    }
  }

  rule {
    action      = "deny(403)"
    description = "owasp-crs-v030001 rce"
    priority    = 904
    match {
      expr {
        expression = "evaluatePreconfiguredExpr('rce-canary')"
      }
    }
  }

  rule {
    action      = "deny(403)"
    description = "owasp-crs-v030001 XSS"
    priority    = 905
    match {
      expr {
        expression = "evaluatePreconfiguredExpr('xss-canary')"
      }
    }
  }

  rule {
    action      = "deny(403)"
    description = "common crawlers"
    priority    = 1200
    match {
      expr {
        expression = "request.path.matches('/Autodiscover|/bin/|/ecp/|/owa/|/vendor/|/ReportServer|/_ignition|/index.php')"
      }
    }
  }

  rule {
    action   = "throttle"
    priority = 2147483647
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["*"]
      }
    }
    description = "default rule"

    rate_limit_options {
      conform_action = "allow"
      exceed_action  = "deny(429)"

      enforce_on_key = ""

      enforce_on_key_configs {
        enforce_on_key_type = "IP"
      }

      rate_limit_threshold {
        count        = var.rate_limit_threshold_count_per_min
        interval_sec = 60
      }
    }
  }
}
