kubectl apply -f sim-deployment.yaml 
kubectl apply -f sim-service.yaml 
kubectl get pods
kubectl get service
kubectl apply -f order_gateway.yaml 
kubectl apply -f order-gateway-service.yaml 
kubectl apply -f mdgateway.yaml 
kubectl apply -f md-gateway-service.yaml
