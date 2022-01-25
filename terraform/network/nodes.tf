resource "kubernetes_deployment" "nodes" {
  depends_on = [
    kubernetes_namespace.network,
    kubernetes_service.bootnode,
    kubernetes_deployment.bootnode
  ]

  count = 5
  metadata {
    name      = "node-${count.index + 1}"
    namespace = "network"
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "node-${count.index + 1}"
      }
    }
    template {
      metadata {
        labels = {
          app = "node-${count.index + 1}"
        }
      }
      spec {
        container {
          image = "mrshah2/kp:network"
          image_pull_policy = "Always"
          name  = "node-${count.index + 1}-container"
          command = [
              "sh",
              "-c",
              "cp -r /nw_data/qdata_${count.index} /qdata && start.sh --bootnode=\"enode://0ea4ced154bdcc26f7a36aca36b7d0f403ad378a3305b196acd638c017ec8fc95d862bb98be461d254fd91f34a621420ee5b871011ece06aa4a8c9225dc8934d@bootnode.network:30301\" --raftInit  --networkid 2018"
          ]
          port {
            container_port = 8545
          }
          port {
            container_port = 8546
          }
          port {
            container_port = 30303
          }
          port {
            container_port = 31000
          }
        }
      }
    }
  }
}
resource "kubernetes_service" "nodes" {

  count = 5

  depends_on = [
    kubernetes_namespace.network,
    kubernetes_service.bootnode,
    kubernetes_deployment.bootnode
  ]
#   for_each = resource.kubernetes_deployment.nodes

  metadata {
    name      = "node-${count.index + 1}"
    namespace = "network"
  }
  spec {
    selector = {
      app = "node-${count.index + 1}"
    }
    type = "NodePort"
    port {
      name = "port1"
      node_port   = 32001 + count.index
      port        = 32001
      target_port = 8545
    }
    port {
      name = "port2"
      node_port   = 32501 + count.index
      port        = 32501
      target_port = 8546
    }
    port {
      name = "port3-udp"
      port        = 30303
      target_port = 30303
      protocol = "UDP"
    }
    port {
      name = "port3"
      port        = 30303
      target_port = 30303
    }
     port {
      name = "raft"
      port        = 31000
      target_port = 31000
    }
 
  }
}
