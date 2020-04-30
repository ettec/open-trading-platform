package k8s

import (
	"flag"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/klog"
	"log"
	"net"
	"os"
	"path/filepath"
)

func GetK8sClientSet(external bool) *kubernetes.Clientset {
	var clientSet *kubernetes.Clientset
	if external {
		var kubeconfig *string
		if home := homeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}

		// create the clientSet
		clientSet, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

	} else {

		var config *rest.Config
		var err error
		if tproot, exists := os.LookupEnv("TELEPRESENCE_ROOT"); exists {
			config, err = InClusterTelepresenceConfig(tproot)
			if err != nil {
				panic(err.Error())
			}
		} else {
			config, err = rest.InClusterConfig()
			if err != nil {
				panic(err.Error())
			}
		}

		// creates the clientSet
		clientSet, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	}
	return clientSet
}

func InClusterTelepresenceConfig(tproot string) (*rest.Config, error) {

	log.Printf("Using telepresence root %v", tproot)

	tokenFile := tproot + "/var/run/secrets/kubernetes.io/serviceaccount/token"
	rootCAFile := tproot + "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"

	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		return nil, rest.ErrNotInCluster
	}

	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}

	tlsClientConfig := rest.TLSClientConfig{}

	if _, err := certutil.NewPool(rootCAFile); err != nil {
		klog.Errorf("Expected to load root CA config from %s, but got err: %v", rootCAFile, err)
	} else {
		tlsClientConfig.CAFile = rootCAFile
	}

	return &rest.Config{
		// TODO: switch to using cluster DNS.
		Host:            "https://" + net.JoinHostPort(host, port),
		TLSClientConfig: tlsClientConfig,
		BearerToken:     string(token),
		BearerTokenFile: tokenFile,
	}, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
