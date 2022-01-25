resource "kubernetes_namespace" "network" {
  metadata {
    name = "network"
  }
}

resource "kubernetes_deployment" "bootnode" {
  depends_on = [
    kubernetes_namespace.network
  ]
  metadata {
    name      = "bootnode"
    namespace = "network"
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "bootnode"
      }
    }
    template {
      metadata {
        labels = {
          app = "bootnode"
        }
      }
      spec {
        container {
          image = "mrshah2/kp:network"
          name  = "bootnode-container"
          command = [
              "sh",
              "-c",
              "cp -r /nw_data/bootnode/* /qdata && bootnode -nodekey /qdata/nodekey"
            #   "bootnode", "-nodekey", "/qdata/nodekey"]
          ]
          port {
            container_port = 30301
          }
          volume_mount {
            mount_path = "/qdata"
            name = "data"
          }
        }
        volume {
          name = "data"
          empty_dir {
            
          }
        }
      }
    }
  }
}
resource "kubernetes_service" "bootnode" {
  depends_on = [
    kubernetes_namespace.network
  ]
  metadata {
    name      = "bootnode"
    namespace = "network"
  }
  spec {
    selector = {
      app = "bootnode"
    }
    type = "NodePort"
    port {
      node_port   = 30301
      port        = 30301
      target_port = 30301
    }
  }
}