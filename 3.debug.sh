#/bin/bash
NS=hybrid

NAME=$(kubectl get pod -n $NS | grep -E 'hcp-policy-engine' | awk '{print $1}')

echo "Exec Into '"$NAME"'"

#kubectl exec -it $NAME -n $NS /bin/sh
kubectl logs -n $NS $NAME
