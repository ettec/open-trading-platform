usage() { echo "Usage: $0 [-v <version number>] [-r <docker repo>] [-m]  where -v is the version number, omit this flag to install latest ci build. -r is the docker repo, omit this flag to use the default repo. -m flag should be used if installing on microk8s  " 1>&2; exit 1; }

DOCKERREPO="ettec"

while getopts ":v:r:m" o; do
    case "${o}" in
        v)
            VERSION=${OPTARG}
            ;;
        m)
            USEMICROK8S="true"
            ;;
        r)
            DOCKERREPO=${OPTARG}
            ;;    
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))




TAG=$VERSION
if [ -z "$VERSION" ]; then 
	printf "Installing latest Open Trading Platform build\n"; 
	TAG="latest"
else 
       printf "Installing Open Trading Platform version $VERSION\n"; 
fi

if [ "$USEMICROK8S" = "true" ];  then
 echo installing into MicroK8s cluster
 shopt -s expand_aliases
 alias kubectl=microk8s.kubectl
 alias helm=microk8s.helm3
else
 echo installing into kubernetes cluster using kubectl current context
fi


DIRECTORY=$(cd `dirname $0` && pwd)
cd $DIRECTORY 

#Binami chart repo
echo adding bitnami chart repo...
helm repo add bitnami https://charts.bitnami.com/bitnami

#Kafka

echo installing Kafka...

kubectl create ns kafka

helm install kafka-opentp  --wait --namespace kafka  bitnami/kafka --version 26.2.0 --set "listeners.client.protocol=PLAINTEXT"
if [ $? -ne 0 ]; then
   echo "Failed to install kafka"
   exit 1		
fi


#Kafka topics
echo creating kafka topics
kubectl run kafka-opentp-client --restart='Never' --image docker.io/bitnami/kafka:3.6.0-debian-11-r0 --namespace kafka --command -- sleep infinity
if [ $? -ne 0 ]; then
   echo "Failed to start kafka client"
   exit 1		
fi

kubectl wait pods -n kafka -l run=kafka-opentp-client --for condition=Ready --timeout=90s

#Orders Topic
kubectl exec --tty -i kafka-opentp-client --namespace kafka -- bash -c "kafka-topics.sh --create --topic orders --bootstrap-server kafka-opentp.kafka.svc.cluster.local:9092"

#Postgres

echo installing Postgresql database...

kubectl create ns postgresql

helm install opentp --wait --namespace postgresql bitnami/postgresql --version 13.1.5 --set-file primary.pgHbaConfiguration=./pb_hba_no_sec.conf --set volumePermissions.enabled=true

if [ $? -ne 0 ]; then
   echo "Failed to install postgres"
   exit 1		
fi



echo loading data into Postgresql database...
export POSTGRES_PASSWORD=$(kubectl get secret --namespace postgresql opentp-postgresql -o jsonpath="{.data.postgres-password}" | base64 --decode)

kubectl run opentp-postgresql-client --rm --tty -i --restart='Never' --namespace postgresql --image  ${DOCKERREPO}/otp-dataload:${TAG} --env="POSTGRESQL_PASSWORD=$POSTGRES_PASSWORD" --command -- psql --host opentp-postgresql -U postgres -d postgres -p 5432 -a -f ./opentp.db

if [ $? -ne 0 ]; then
   echo "Failed to load initial data set"
   exit 1		
fi

#Envoy

echo installing Envoy...

kubectl create ns envoy
helm install opentp-envoy --wait --namespace=envoy ./charts/envoy -f envoy-config-helm-values.yaml 
if [ $? -ne 0 ]; then
   echo "Failed to install envoy"
   exit 1		
fi

kubectl patch service envoy --namespace envoy --type='json' -p='[{"op": "replace", "path": "/spec/sessionAffinity", "value": "ClientIP"}]'


#Opentp

echo installing Open Trading Platform...


helm install --wait --timeout 1200s otp-${VERSION} ../helm-otp-chart/ --set dockerRepo=${DOCKERREPO} --set dockerTag=${TAG}
if [ $? -ne 0 ]; then
   echo "Failed to install open trading platfrom"
   exit 1		
fi


#Instructions to start client
OTPPORT=$(kubectl get svc --namespace=envoy -o go-template='{{range .items}}{{range.spec.ports}}{{if .nodePort}}{{.nodePort}}{{"\n"}}{{end}}{{end}}{{end}}')

echo
echo Open Trading Platform is running. To start a client point your browser at port $OTPPORT and login as trader1 







