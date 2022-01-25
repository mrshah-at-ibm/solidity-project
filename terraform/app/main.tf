resource "kubernetes_namespace" "mrshah" {
  metadata {
    name = "mrshah"
  }
}

resource "kubernetes_role" "allow-app" {
  metadata {
    name = "allow-app"
    namespace = "mrshah"
    # labels = {
    #   test = "MyRole"
    # }
  }

  rule {
    api_groups     = [""]
    resources      = ["secrets"]
    resource_names = ["privatekey"]
    verbs          = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }
  rule {
    api_groups = [""]
    resources  = ["configmaps"]
    resource_names = ["config"]
    verbs      = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }
}

resource "kubernetes_role_binding" "allow-app" {
  metadata {
    name      = "allow-app"
    namespace = "mrshah"
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "Role"
    name      = "allow-app"
  }
  # subject {
  #   kind      = "User"
  #   name      = "admin"
  #   api_group = "rbac.authorization.k8s.io"
  # }
  subject {
    kind      = "ServiceAccount"
    name      = "default"
    namespace = "mrshah"
  }
  # subject {
  #   kind      = "Group"
  #   name      = "system:masters"
  #   api_group = "rbac.authorization.k8s.io"
  # }
}

resource "kubernetes_config_map" "config" {
  metadata {
    name = "config"
    namespace = kubernetes_namespace.mrshah.metadata.0.name
  }

  data = {
    "config" = "${file("${path.module}/config.yaml")}"
  }
}

resource "kubernetes_deployment" "mrshah" {
  metadata {
    name      = "app"
    namespace = kubernetes_namespace.mrshah.metadata.0.name
  }
  spec {
    replicas = 2
    selector {
      match_labels = {
        app = "mrshah-app"
      }
    }
    template {
      metadata {
        labels = {
          app = "mrshah-app"
        }
      }
      spec {
        container {
          image = "mrshah2/kp:app"
          name  = "app-container"
          port {
            container_port = 3000
          }
          env {
            name = "INCLUSTER"
            value = "true"
          }
          env {
            name = "NAMESPACE"
            value = kubernetes_namespace.mrshah.metadata.0.name
          }
        }
      }
    }
  }
}
resource "kubernetes_service" "mrshah" {
  metadata {
    name      = "app"
    namespace = kubernetes_namespace.mrshah.metadata.0.name
  }
  spec {
    selector = {
      app = kubernetes_deployment.mrshah.spec.0.template.0.metadata.0.labels.app
    }
    type = "NodePort"
    port {
      node_port   = 30000
      port        = 3000
      target_port = 3000
    }
  }
}