#/bin/bash
NS=nsnsns

NAME=$(kubectl get pod -n $NS | grep -E 'template-controller' | awk '{print $1}')

echo "Exec Into '"$NAME"'"

#kubectl exec -it $NAME -n $NS /bin/sh
kubectl logs -n $NS $NAME
