resource "kubernetes_manifest" "mrshah_app_certificate" {
  manifest = {
    "apiVersion" = "cert-manager.io/v1"
    "kind" = "Certificate"
    "metadata" = {
      "name" = "app-certificate"
      "namespace" = "mrshah"
    }
    "spec" = {
      "dnsNames" = [
        "app.mrshah.space",
      ]
      "issuerRef" = {
        "group" = "cert-manager.io"
        "kind" = "ClusterIssuer"
        "name" = "cert-manager-global"
      }
      "secretName" = "app-certificate"
    }
  }
}
