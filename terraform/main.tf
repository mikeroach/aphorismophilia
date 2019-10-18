provider "kubernetes" {
  version          = "1.9.0"
  load_config_file = true
  config_path      = "../secrets/ephemeral-kubeconfig" // Locally generated via Makefile-invoked gcloud
}

resource "kubernetes_namespace" "aphorismophilia-ephemeral-namespace" {
  metadata {
    name = var.service_aphorismophilia_namespace
  }
}

module "aphorismophilia-ephemeral" {
  source                            = "github.com/mikeroach/aphorismophilia-terraform?ref=v12"
  dns_domain                        = var.dns_domain
  dns_hostname                      = var.dns_hostname
  dockerhub_credentials             = var.dockerhub_credentials
  service_aphorismophilia_version   = var.service_aphorismophilia_version
  service_aphorismophilia_namespace = kubernetes_namespace.aphorismophilia-ephemeral-namespace.metadata[0].name
}