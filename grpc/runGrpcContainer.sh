kubectl run -i --tty --rm $1 --image=192.168.1.200:5000/grpcshell:1 --restart=Never -- sh
