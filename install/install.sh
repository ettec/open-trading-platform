
# If installing on a non microk8s cluster comment out the three lines below
shopt -s expand_aliases
alias kubectl=microk8s.kubectl
alias helm=microk8s.helm3


#Kafka

echo installing Kafka...

kubectl create ns kafka

helm repo add incubator http://storage.googleapis.com/kubernetes-charts-incubator

helm install kafka-opentp --wait --namespace kafka incubator/kafka

#install kafka cmd line client and setup orders topic

 
kubectl apply --wait -f kafka_cmdline_client.yaml


kubectl exec -it --namespace=kafka cmdlineclient -- /bin/bash -c "kafka-topics --zookeeper kafka-opentp-zookeeper:2181 --topic orders --create --partitions 1 --replication-factor 1"


#Postgres

echo installin Postgreql database...

kubectl create ns postgresql

helm install opentp --wait --namespace postgresql bitnami/postgresql --set-file pgHbaConfiguration=./pb_hba_no_sec.conf

export POSTGRES_PASSWORD=$(kubectl get secret --namespace postgresql opentp-postgresql -o jsonpath="{.data.postgresql-password}" | base64 --decode)

kubectl run opentp-postgresql-client --rm --tty -i --restart='Never' --namespace postgresql --image  ettec/opentp-ci-build:data-loader-client --env="PGPASSWORD=$POSTGRES_PASSWORD" --command -- psql --host opentp-postgresql -U postgres -d postgres -p 5432 -a -f ./opentp.db

#Envoy

echo installing Envoy...

kubectl create ns envoy

helm install opentp-envoy --wait --namespace=envoy stable/envoy -f envoy-config-helm-values.yaml 
kubectl patch service envoy --namespace envoy --type='json' -p='[{"op": "replace", "path": "/spec/sessionAffinity", "value": "ClientIP"}]'

#Opentp

echo installing Open Trading Platform...

helm install --wait --timeout 1200s otp-v1 ../helm-otp-chart/

#Instructions to start client
OTPPORT=$(kubectl get svc --namespace=envoy -o go-template='{{range .items}}{{range.spec.ports}}{{if .nodePort}}{{.nodePort}}{{"\n"}}{{end}}{{end}}{{end}}')

echo Open Trading Platform is running. To start a client point your browser at port $OTPPORT and login as trader1 







