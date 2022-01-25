resource "kubernetes_ingress" "app_ingress" {
  metadata {
    name = "app-ingress"
    namespace = "mrshah"
    annotations = {
      "kubernetes.io/ingress.class" = "nginx"
      "nginx.ingress.kubernetes.io/ssl-redirect" = "false"
      "nginx.ingress.kubernetes.io/use-regex" = "true"
    }

  }

  spec {
    backend {
      service_name = "app"
      service_port = 3000
    }

    rule {
      host = "app.mrshah.space"
      http {
        path {
          backend {
            service_name = "app"
            service_port = 3000
          }

          path = "/*"
        }
     }
    }

    tls {
      secret_name = "app-certificate"
    }

  }
}
