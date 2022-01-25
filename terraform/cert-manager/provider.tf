provider "kubectl" {
  config_path = "./kubeconfig.yaml"
}
provider "helm" {
  kubernetes {
    config_path = "./kubeconfig.yaml"
  }
}
terraform {
  required_providers {
    kubectl = {
      source  = "hashicorp/kubernetes"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "2.3.0"
    }
  }
}