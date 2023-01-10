provider "kubernetes" {
  version = "2.16.1"
  /* TODO: Consider consuming cluster API endpoint from remote state lookup and
  issue token-based auth via service account direct in Terraform. On the other
  hand, we still want to support the "local development against production-like
  live environments" use case via kubectl/Telepresence/Okteto/etc. */
  config_path = "../secrets/ephemeral-kubeconfig" // Locally generated via Makefile-invoked gcloud

  ignore_annotations = [
    "^cloud.google.com\\/neg.*",
  ]
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