#cluster=gpu
link=stats/summary
#link=metrics/resource
#link=metrics/cadvisor
#link=stats
#link=pods
#link=metrics
#APISERVER=$(kubectl config view --minify --context $cluster | grep server | cut -f 2- -d ":" | tr -d " ")
SECRET_NAME=$(kubectl get secrets -n hcp | grep hcp-metric-collector | cut -f1 -d ' ')
TOKEN=$(kubectl describe secret $SECRET_NAME -n hcp | grep -E '^token' | cut -f2 -d':' | tr -d " ")
curl https://10.0.5.84:10250/$link --header "Authorization: Bearer $TOKEN" --insecure
