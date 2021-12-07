terraform {
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
      version = "~>2.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.0.0"
    }
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = ">= 1.7.0"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "kubernetes" {
  alias = "master"
  config_path = "~/.kube/config" 
  config_context = "master"
}

provider "kubernetes-alpha" {
  # alias = "kubernetes-alpha"
  config_path = "~/.kube/config"
  config_context = "master"
}


provider "kubectl" {
  config_path = "~/.kube/config"
  config_context = "master"
}

data "azurerm_kubernetes_cluster" "main" {
  name                = "k8stest"
  resource_group_name = "hcpResourceGroup"
}

provider "kubernetes" {
    alias = "azureCluster"
    host                   = data.azurerm_kubernetes_cluster.main.kube_config.0.host
    username               = data.azurerm_kubernetes_cluster.main.kube_config.0.username
    password               = data.azurerm_kubernetes_cluster.main.kube_config.0.password
    client_certificate     = base64decode(data.azurerm_kubernetes_cluster.main.kube_config.0.client_certificate)
    client_key             = base64decode(data.azurerm_kubernetes_cluster.main.kube_config.0.client_key)
    cluster_ca_certificate = base64decode(data.azurerm_kubernetes_cluster.main.kube_config.0.cluster_ca_certificate)
}

# merge kubeconfig
resource "null_resource" "mergekubeconfig" {
    provisioner "local-exec" {
      command = "az aks get-credentials --admin --name ${data.azurerm_kubernetes_cluster.main.name} --resource-group ${data.azurerm_kubernetes_cluster.main.resource_group_name}"
  }
}

# 1. CREATE namespace "kube-federation-system" in AKS CLUSTER
resource "kubernetes_namespace" "hcp" {
  provider = kubernetes.azureCluster
    metadata {
        name = "kube-federation-system"
    }
}

# 2. CREATE service account in AKS CLUSTER
resource "kubernetes_service_account" "hcp" {
  provider = kubernetes.azureCluster
    metadata {
        name = "${data.azurerm_kubernetes_cluster.main.name}"
        namespace = "kube-federation-system"
        annotations = {
          "kubernetes.io/service-account.name" = "default"
        }
    }
}

resource "kubernetes_secret" "hcp" {
  provider = kubernetes.azureCluster
    metadata {
        name = "${kubernetes_service_account.hcp.default_secret_name}"
        namespace = "kube-federation-system"
        annotations = {
          "kubernetes.io/service-account.name" = "default"
         }
    }
    type =  "kubernetes.io/service-account-token"
}

data "kubernetes_secret" "secretData" {
  provider = kubernetes.azureCluster
  metadata {
    name = "${kubernetes_service_account.hcp.default_secret_name}"
    namespace = "kube-federation-system"
    annotations = {
      "kubernetes.io/service-account.name" = "default"
    }
  }
}

# # 2-1. CREATE secret yamlFILE 
resource "null_resource" "createYAMLfile" {
    provisioner "local-exec" {
      command = "kubectl get secret ${kubernetes_service_account.hcp.default_secret_name} -n ${kubernetes_service_account.hcp.metadata.0.namespace} -o yaml > ${kubernetes_service_account.hcp.default_secret_name}.yaml --context ${data.azurerm_kubernetes_cluster.main.name}-admin"
    }
}

# 3. CREATE cluster role in AKS CLUSTER
resource "kubernetes_cluster_role" "hcp" {
  provider = kubernetes.azureCluster  
  metadata {
    name = "terraform-example"
  }

  rule {
      api_groups = [""]
      resources  = ["namespaces"]
      verbs      = ["get", "list", "update", "create", "patch"]
  }
}

# 4. CREATE cluster role binding in AKS CLUSTER
resource "kubernetes_cluster_role_binding" "hcp" {
  provider = kubernetes.azureCluster
  metadata {
    name = "terraform-example"
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.hcp.metadata[0].name
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.hcp.metadata[0].name
    namespace = kubernetes_namespace.hcp.metadata[0].name
  }
}




# 5. Get & CREATE secret 
# resource "kubectl_manifest" "createSecret" {

#   yaml_body = <<YAML
#     apiVersion: v1
#     data:
#       ca.crt: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUU2VENDQXRHZ0F3SUJBZ0lSQUtGSittdmNFR1dWbDVnSHluVTRBMWt3RFFZSktvWklodmNOQVFFTEJRQXcKRFRFTE1Ba0dBMVVFQXhNQ1kyRXdJQmNOTWpFd056SXdNREUxT0RFeFdoZ1BNakExTVRBM01qQXdNakE0TVRGYQpNQTB4Q3pBSkJnTlZCQU1UQW1OaE1JSUNJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBZzhBTUlJQ0NnS0NBZ0VBCnE0OVFtWU0xNDhsK2VzelJmbC9uWkpjblQzd296cW9KQzdJVDdOOGJxK1JUK2lvcHFXZVUzY3dnTm5PS0ZwM0IKKy8yOXpGWTBqVzZZRXVTUEVpSWFxak9KL0xaUUwyU0t5bytOUTVsV0lpaEQrREM1aEtVZDNLTnJXbDZKRjRRYgpCaDZlT1FKUkFpRXlTNi9uNkY4c1daOWNHTERmcjZkSGkveTRHazhoM085c0NYTDV0d3ZGUVcxeklrRlR1ZUFxCmZVRlk2VjUvTWcrcWFZZS9QY1VhRm1PWHpTY1Y5RUpoZVorRjYxQU5sVG9CVEt0YU1IMHkyNFYzYjNYVEVwSjYKblMxSzI4ejFqa2ErVEpXQ1BrU1RCaEdvYzArYVF4LzVlTDY1RnVFS3FkZDlQcTJnd3p0bVhOTkx6Smpra3BRVQp5VmVnUWFFRUk3Sk5hZFZkb3RVdm5ENXJLcElCZ2VXYUtuVmdlMHBwNlRiWEFBUFJPNS9YWEl5WkV0RUlVbVcvCmJUdzdUeXZ5MEkxdVNQOU9BRXhQbGVNRzBiK2Vvakg4Y1NzMk80T3FscmlrbTRHSDFTdFRkc3VFaTRXd1JNQlMKT1lGOXBCUVkzTUFFVTN1Z04yVEczbUlVUVA2cmxiQWRIVU16VE52NHo5Q3VjbjAzbG5JQVdQQ0dmb2hLdzliUgo1STJUTy9UV1BLMkQrcmh6S2w4VzN6WmhVQ0doN2R6VE9jTGZ4NC9aTVA1b1B4RDRRRThBZ1dpTW9PMVpuV2R4CmZYdjdhOWYrVVdrb3psNU1oLy9zOTNEZ09rNlhGdVF0QTM5MWJNYWVuejZGcVVuRjZrZnd0QVFMUWs3T1h6a2YKUE5FNGtGZk5yY1hseDBneHNWSTVrMWNNL0lYQjRMZGl4MkFYenVGcjBtRUNBd0VBQWFOQ01FQXdEZ1lEVlIwUApBUUgvQkFRREFnS2tNQThHQTFVZEV3RUIvd1FGTUFNQkFmOHdIUVlEVlIwT0JCWUVGSG9lY1hwZWQ0TmJodUZiClhFM3lmaXMxTWg1OU1BMEdDU3FHU0liM0RRRUJDd1VBQTRJQ0FRQmxwNWhiRGd0bWpMNVREeVlFUk5lM25RdUkKYjhPYnVPUjE3K2xOUnRpNHBxeEdwZWVjSXhicWUwZFc0d1V4NU9mVUtJa3BOTVA1NlFHWGRoRE5aK3V4b0VrZwpoczRvVy9vYzdoUlkvVkhRZGExcmgra0ZNczdseWdrUFdQbmZWWThrbXROczRLT2EzVjBLWi9TS0JzV2RrcW9mCnVGeTdIZm1rOGtrVlpDcWVBcEtnaFI4eEtKMk5jMnNtV1RJb3FvN0J5WVlJOU9JZ3FkUzJVQ1FSOFVkVzNTNkQKdlQxTnRmWGFRMDMvYWVZV1pZTHIrMFh2Szh5aTVZczY3OXJqbVYyazg5SC9sVjZON2FCLzUveVI1aEdRUUtmNwpsaVB6LzZqUHlFNjRFeHJvUVVycVNyNWJXSmwvVFBCOENuaHg2LzVkRXE1MVdJSUFRZXR6NkdRY1pNLzRVT0YzClp5YllSdStjYjZCT2hsOVRYMnJ2ZW9oYjRSUzNKSDhsR1h5SHhZSGxYUFFsdWJVVHZvaG4zNHhPOTRUdngvVEUKN1JqZERsaVhRWjFRbXJYRUlJcHVxcTdtSEZIMWdKYVJPYzBKUWdyQVlPV2ZiWDU5Y3padlI3SHlra01kK0UwMApDemIvNUNPVU1icTVxQ1orNWRUcUVyc3VINXRFZHdwdC9ZaDZhMk5PZFlwMVR0Q1I2S2VCeXIxMnFSb1FiK0c2CmhpUHBWbXU3WmI5Rkhqc0FnQUh1SFo0QTFqMHNWODZOc0ZWQzFiUkF3aTVKc0lZWmxXcnhQQmlONjhpU3B4MzgKNWUvVGZNQmlLR0tnMFRpZ0o2WG9QNkhmSFErdG9oNkpCU1hTOWswUmpnbWN6aWVmaDlQdFdOd29qUEpIRTEvSgpENkp2WWVYUm9oeWdPZXlvUmc9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="
#       namespace: "kube-federation-system"
#       token: "ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklqbHJNV1ZYYWs5WFVWTm9PR2hmVW1oVWMxTm9SVlpyU1hGYU5WZ3dXRVYwWkV4dFIwUjNXV0pTTXpnaWZRLmV5SnBjM01pT2lKcmRXSmxjbTVsZEdWekwzTmxjblpwWTJWaFkyTnZkVzUwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXVZVzFsYzNCaFkyVWlPaUpyZFdKbExXWmxaR1Z5WVhScGIyNHRjM2x6ZEdWdElpd2lhM1ZpWlhKdVpYUmxjeTVwYnk5elpYSjJhV05sWVdOamIzVnVkQzl6WldOeVpYUXVibUZ0WlNJNkltczRjM1JsYzNRdGRHOXJaVzR0Ym10bk1tMGlMQ0pyZFdKbGNtNWxkR1Z6TG1sdkwzTmxjblpwWTJWaFkyTnZkVzUwTDNObGNuWnBZMlV0WVdOamIzVnVkQzV1WVcxbElqb2laR1ZtWVhWc2RDSXNJbXQxWW1WeWJtVjBaWE11YVc4dmMyVnlkbWxqWldGalkyOTFiblF2YzJWeWRtbGpaUzFoWTJOdmRXNTBMblZwWkNJNkltWXdZakZtWlRnMExURXpaR010TkdNMU5DMWlOMll4TFRrNU1qWXlORGsxTmpWak55SXNJbk4xWWlJNkluTjVjM1JsYlRwelpYSjJhV05sWVdOamIzVnVkRHByZFdKbExXWmxaR1Z5WVhScGIyNHRjM2x6ZEdWdE9tUmxabUYxYkhRaWZRLm9lVFdsYll6d0tDb3dvRUlfZjhSLVc5TmhzcFhZdHoyVWRfMk1IVmY2UldUcEwybFVkbVhUdzFUNnIzT2paQUdyVE9XZDFpblRjT1lKSWh3bW1uamprc01JNWNGWnJkcFlPRmF4emF4MWtfbEwzbGxHM3V2NnZORkMxQ3Q3MUVKbTJ1YTltSXZyNC1vbVpFX1JveW9JWWc4Q1JSVXVyOE8zRFd4ajhFb1d3TnBTbV9vQkE4RjdCaGc4SXdIWk5LaWJGRFNxaWw5NHJqSUdCOHNNVXAzNUJvbDY5S2lTMTVuOVZPdWJpS1RzQlBQekhiM2FJS2dlZWQtQ282d0xoNnRWdUNIS0t6U2FhcTBqMUVyeXZBSW9LVm5uY3RDa2RYWkxUa2U1Z2xMVmhEUEQwc1B1eTRIdWtRZUd5c1M2N1FUaEJmbW5zalA5eE9IWmdob0lpcXZHNlRQZGw3U2F1cVFFeVVhUDkyaGh6LW5OTkppMFRZWlhnZXBsZGhwX2o0WE5mN3ZuQzFHeFBhOE1RZjUtbk95eUxoNjJ2S2VjQ2djUmlDVFRySXNmRXVQRlB4UWY2LTZxUHlBZFdyM3IzRVNNTFhoQjF0SUlBWEhZX0d5YTR0Vnp2NXhsYzM4eFNpZWN4QjFLYS1jak9OSnRlUExtUl84aWZ0cGNfQWhtN0ZhMmtKTF9ESVhsMmFwOFBISmRCMUp4TzRpS1JCbGdsTTlyd3NXbFd5VllINjFqWjNkLUpsTTZJeGsxdGc3cC16YmY1NURKenN0cjFZbWR1YlV4bzBSdXVuT0ZxcG42amdadDN6RWZMbDVlYTJ3S2VlWmFVWWtESW51dEZfZWFTdktIcEQwdkFUVkIzeXRrRjRBNnVEOEpvZmh1cjFNcDNGTXk2M0NGYU15d1dj"
#     kind: Secret
#     metadata:
#       annotations:
#         kubernetes.io/service-account.name: default
#       name: "k8stest-token-vvj6f"
#       namespace: "kube-federation-system"
#       selfLink: "/api/v1/namespaces/kube-federation-system/secrets/k8stest-token-vvj6f"
#     type: kubernetes.io/service-account-token
#   YAML
# }

# resource "kubernetes_secret" "kubernetes_tls_secret" {
#   provider = kubernetes
#   metadata {
#     name      = "my_secret"
#     namespace = "${data.kubernetes_secret.secretData.data.namespace}"
#   }

#   data = {
#     // avoid trying to decode plaintext data
#     host  = "${data.azurerm_kubernetes_cluster.main.kube_config.0.host}"
#     token = "${base64decode(data.kubernetes_secret.secretData.data["token"])}"
#   }

#   type = "kubernetes.io/service-account-token"
# }

# resource "kubernetes_manifest" "createSecret" {
#   provider = kubernetes-alpha
#   manifest = {
#     "apiVersion" = "v1"
#     "data" = {
#       "token" = "ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklqbHJNV1ZYYWs5WFVWTm9PR2hmVW1oVWMxTm9SVlpyU1hGYU5WZ3dXRVYwWkV4dFIwUjNXV0pTTXpnaWZRLmV5SnBjM01pT2lKcmRXSmxjbTVsZEdWekwzTmxjblpwWTJWaFkyTnZkVzUwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXVZVzFsYzNCaFkyVWlPaUpyZFdKbExXWmxaR1Z5WVhScGIyNHRjM2x6ZEdWdElpd2lhM1ZpWlhKdVpYUmxjeTVwYnk5elpYSjJhV05sWVdOamIzVnVkQzl6WldOeVpYUXVibUZ0WlNJNkltczRjM1JsYzNRdGRHOXJaVzR0Ym10bk1tMGlMQ0pyZFdKbGNtNWxkR1Z6TG1sdkwzTmxjblpwWTJWaFkyTnZkVzUwTDNObGNuWnBZMlV0WVdOamIzVnVkQzV1WVcxbElqb2laR1ZtWVhWc2RDSXNJbXQxWW1WeWJtVjBaWE11YVc4dmMyVnlkbWxqWldGalkyOTFiblF2YzJWeWRtbGpaUzFoWTJOdmRXNTBMblZwWkNJNkltWXdZakZtWlRnMExURXpaR010TkdNMU5DMWlOMll4TFRrNU1qWXlORGsxTmpWak55SXNJbk4xWWlJNkluTjVjM1JsYlRwelpYSjJhV05sWVdOamIzVnVkRHByZFdKbExXWmxaR1Z5WVhScGIyNHRjM2x6ZEdWdE9tUmxabUYxYkhRaWZRLm9lVFdsYll6d0tDb3dvRUlfZjhSLVc5TmhzcFhZdHoyVWRfMk1IVmY2UldUcEwybFVkbVhUdzFUNnIzT2paQUdyVE9XZDFpblRjT1lKSWh3bW1uamprc01JNWNGWnJkcFlPRmF4emF4MWtfbEwzbGxHM3V2NnZORkMxQ3Q3MUVKbTJ1YTltSXZyNC1vbVpFX1JveW9JWWc4Q1JSVXVyOE8zRFd4ajhFb1d3TnBTbV9vQkE4RjdCaGc4SXdIWk5LaWJGRFNxaWw5NHJqSUdCOHNNVXAzNUJvbDY5S2lTMTVuOVZPdWJpS1RzQlBQekhiM2FJS2dlZWQtQ282d0xoNnRWdUNIS0t6U2FhcTBqMUVyeXZBSW9LVm5uY3RDa2RYWkxUa2U1Z2xMVmhEUEQwc1B1eTRIdWtRZUd5c1M2N1FUaEJmbW5zalA5eE9IWmdob0lpcXZHNlRQZGw3U2F1cVFFeVVhUDkyaGh6LW5OTkppMFRZWlhnZXBsZGhwX2o0WE5mN3ZuQzFHeFBhOE1RZjUtbk95eUxoNjJ2S2VjQ2djUmlDVFRySXNmRXVQRlB4UWY2LTZxUHlBZFdyM3IzRVNNTFhoQjF0SUlBWEhZX0d5YTR0Vnp2NXhsYzM4eFNpZWN4QjFLYS1jak9OSnRlUExtUl84aWZ0cGNfQWhtN0ZhMmtKTF9ESVhsMmFwOFBISmRCMUp4TzRpS1JCbGdsTTlyd3NXbFd5VllINjFqWjNkLUpsTTZJeGsxdGc3cC16YmY1NURKenN0cjFZbWR1YlV4bzBSdXVuT0ZxcG42amdadDN6RWZMbDVlYTJ3S2VlWmFVWWtESW51dEZfZWFTdktIcEQwdkFUVkIzeXRrRjRBNnVEOEpvZmh1cjFNcDNGTXk2M0NGYU15d1dj"
#     }
#     "kind" = "Secret"
#     "metadata" = {
#       "annotations" = {
#         "kubernetes.io/service-account.name" = "default"
#       }
#       "name" = "${data.kubernetes_secret.secretData.metadata.0.name}"
#       "namespace" = "${data.kubernetes_secret.secretData.metadata.0.namespace}"
#     }
#     "type" = "kubernetes.io/service-account-token"
#   }
# }

# TODO: check namespace exists or not
# resource "kubernetes_namespace" "ns" {
#   provider = kubernetes.master
#     metadata {
#         name = "kube-federation-system"
#     }
# }

resource "kubernetes_secret" "sc" {
  provider = kubernetes.master
    metadata {
      name = "${kubernetes_service_account.hcp.default_secret_name}"
      namespace = "kube-federation-system"
      annotations = {
        "kubernetes.io/service-account.name" = "default"
      }
    }
    
    data = {
      "${kubernetes_service_account.hcp.default_secret_name}.yaml" = file("${path.cwd}/${kubernetes_service_account.hcp.default_secret_name}.yaml")
    }
    type = "kubernetes.io/service-account-token"
}


resource "kubernetes_manifest" "kubefedcluster" {
  provider = kubernetes-alpha

  manifest = {
        "apiVersion" = "apiextensions.k8s.io/v1"
        "kind" = "CustomResourceDefinition"
        "metadata" = {
            "name" = "kubefedclusters.core.kubefed.io"
        }
        "spec" = {
            "conversion" = {
            "strategy" = "None"
            }
            "group" = "core.kubefed.io"
            "names" = {
            "kind" = "KubeFedCluster"
            "listKind" = "KubeFedClusterList"
            "plural" = "kubefedclusters"
            "singular" = "kubefedcluster"
            }
            "scope" = "Namespaced"
            "versions" = [
            {
                "additionalPrinterColumns" = [
                {
                    "jsonPath" = ".metadata.creationTimestamp"
                    "name" = "age"
                    "type" = "date"
                },
                {
                    "jsonPath" = ".status.conditions[?(@.type=='Ready')].status"
                    "name" = "ready"
                    "type" = "string"
                },
                ]
                "name" = "v1beta1"
                "schema" = {
                "openAPIV3Schema" = {
                    "properties" = {
                    "apiVersion" = {
                        "type" = "string"
                    }
                    "kind" = {
                        "type" = "string"
                    }
                    "metadata" = {
                        "type" = "object"
                    }
                    "spec" = {
                        "properties" = {
                        "apiEndpoint" = {
                            "type" = "string"
                        }
                        "caBundle" = {
                            "format" = "byte"
                            "type" = "string"
                        }
                        "disabledTLSValidations" = {
                            "items" = {
                            "type" = "string"
                            }
                            "type" = "array"
                        }
                        "proxyURL" = {
                            "type" = "string"
                        }
                        "secretRef" = {
                            "properties" = {
                            "name" = {
                                "type" = "string"
                            }
                            }
                            "required" = [
                            "name",
                            ]
                            "type" = "object"
                        }
                        }
                        "required" = [
                        "apiEndpoint",
                        "secretRef",
                        ]
                        "type" = "object"
                    }
                    "status" = {
                        "properties" = {
                        "conditions" = {
                            "items" = {
                            "properties" = {
                                "lastProbeTime" = {
                                "format" = "date-time"
                                "type" = "string"
                                }
                                "lastTransitionTime" = {
                                "format" = "date-time"
                                "type" = "string"
                                }
                                "message" = {
                                "type" = "string"
                                }
                                "reason" = {
                                "type" = "string"
                                }
                                "status" = {
                                "type" = "string"
                                }
                                "type" = {
                                "type" = "string"
                                }
                            }
                            "required" = [
                                "lastProbeTime",
                                "status",
                                "type",
                            ]
                            "type" = "object"
                            }
                            "type" = "array"
                        }
                        "region" = {
                            "type" = "string"
                        }
                        "zones" = {
                            "items" = {
                            "type" = "string"
                            }
                            "type" = "array"
                        }
                        }
                        "required" = [
                        "conditions",
                        ]
                        "type" = "object"
                    }
                    }
                    "required" = [
                    "spec",
                    ]
                    "type" = "object"
                }
                }
                "served" = true
                "storage" = true
                "subresources" = {
                "status" = {}
                }
            },
            ]
        }
    }
}


resource "kubernetes_manifest" "createKubefedcluster" {
  provider = kubernetes-alpha

  manifest = {
        "apiVersion" = "core.kubefed.io/v1beta1"
        "kind" = "KubeFedCluster"
        "metadata" = {
            "name" = "${data.azurerm_kubernetes_cluster.main.name}"
            "namespace" = "kube-federation-system"
        }
        "spec" = {
          "apiEndpoint" = "https://10.0.5.83:6443"
          # "secretRef" = {
          #   "name" = "${data.azurerm_kubernetes_cluster.main.name}-axxwz"
          # }
        }
  }
}
