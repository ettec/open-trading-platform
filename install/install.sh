# Usage: install.sh <TAG> <use own cluster>    e.g. install.sh 1.0.8  leave it blank to install latest ci build
VERSION=$1
DOCKERREPO="ettec/opentp:"
TAG=-$VERSION
if [ -z "$VERSION" ]; then 
	echo "installing latest Open Trading Platform ci build"; 
	DOCKERREPO="ettec/opentp-ci-build:"
	TAG=""
else 

       echo "installing Open Trading Platform version $VERSION"; 
fi

# If installing on a non microk8s cluster comment out the three lines below
USINGOWNCLUSTER=$2
if [ "$USINGOWNCLUSTER" = "true" ];  then
 echo installing into kubernetes cluster using kubectl current context
else
 echo installing into MicroK8s cluster
 shopt -s expand_aliases
 alias kubectl=microk8s.kubectl
 alias helm=microk8s.helm3
fi


DIRECTORY=$(cd `dirname $0` && pwd)
cd $DIRECTORY 

#Kafka

echo installing Kafka...

kubectl create ns kafka

helm repo add incubator http://storage.googleapis.com/kubernetes-charts-incubator

helm install kafka-opentp --wait --namespace kafka incubator/kafka

#install kafka cmd line client 

 
kubectl apply --wait -f kafka_cmdline_client.yaml




#Postgres

echo installing Postgreql database...

kubectl create ns postgresql

helm repo add bitnami https://charts.bitnami.com/bitnami
helm install opentp --wait --namespace postgresql bitnami/postgresql --set-file pgHbaConfiguration=./pb_hba_no_sec.conf --set volumePermissions.enabled=true


export POSTGRES_PASSWORD=$(kubectl get secret --namespace postgresql opentp-postgresql -o jsonpath="{.data.postgresql-password}" | base64 --decode)

kubectl run opentp-postgresql-client --rm --tty -i --restart='Never' --namespace postgresql --image  ${DOCKERREPO}data-loader-client${TAG} --env="PGPASSWORD=$POSTGRES_PASSWORD" --command -- psql --host opentp-postgresql -U postgres -d postgres -p 5432 -a -f ./opentp.db

#Envoy

echo installing Envoy...

kubectl create ns envoy

helm repo add stable https://kubernetes-charts.storage.googleapis.com
helm install opentp-envoy --wait --namespace=envoy stable/envoy -f envoy-config-helm-values.yaml 
kubectl patch service envoy --namespace envoy --type='json' -p='[{"op": "replace", "path": "/spec/sessionAffinity", "value": "ClientIP"}]'

#Orders topic
kubectl exec -it --namespace=kafka cmdlineclient -- /bin/bash -c "kafka-topics --zookeeper kafka-opentp-zookeeper:2181 --topic orders --create --partitions 1 --replication-factor 1"


#Opentp

echo installing Open Trading Platform using tag $TAG...


helm install --wait --timeout 1200s otp-v1 ../helm-otp-chart/ --set dockerRepo=${DOCKERREPO} --set dockerTag=${TAG}

#Instructions to start client
OTPPORT=$(kubectl get svc --namespace=envoy -o go-template='{{range .items}}{{range.spec.ports}}{{if .nodePort}}{{.nodePort}}{{"\n"}}{{end}}{{end}}{{end}}')

echo Open Trading Platform is running. To start a client point your browser at port $OTPPORT and login as trader1 







