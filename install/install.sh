# enable aliases in script
shopt -s expand_aliases
alias kubectl=microk8s.kubectl
alias helm=microk8s.helm3

# K8s generic cmds  (works with kubeadm cluster/minikube)

#Kafka
kubectl create ns kafka

helm repo add incubator http://storage.googleapis.com/kubernetes-charts-incubator

helm install kafka-opentp --namespace kafka incubator/kafka

kubectl apply -f kafka_cmdline_client.yaml


#Postgres:

kubectl create ns postgresql

helm install opentp --namespace postgresql bitnami/postgresql --set-file pgHbaConfiguration=./pb_hba_no_sec.conf


export POSTGRES_PASSWORD=$(kubectl get secret --namespace postgresql opentp-postgresql -o jsonpath="{.data.postgresql-password}" | base64 --decode)


#Envoy:

kubectl create ns envoy

helm install opentp-envoy --namespace=envoy stable/envoy -f envoy-config-helm-values.yaml 
kubectl patch service envoy --namespace envoy --type='json' -p='[{"op": "replace", "path": "/spec/sessionAffinity", "value": "ClientIP"}]'


#Opentp app:

helm install otp-v1 ../otpchart








