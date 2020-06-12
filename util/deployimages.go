package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Repos struct {
	Repositories []string `yaml:"repositories"`
}

type Tags struct {
	Name string   `yaml:"name"`
	Tags []string `yaml:"tags"`
}

type upgradeCandidate struct {
	deployment v1.Deployment
	newImage string

}


func main() {

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 1 {
		panic("usage: deployimages <docker registry address>")
	}

	dockerServerAddress := argsWithoutProg[0]
	dockerUrl := "https://" + dockerServerAddress
	apiVersion := "v2"

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	r, err := http.Get(dockerUrl + "/" + apiVersion + "/_catalog")
	if err != nil {
		panic(err)
	}

	repoString, err := ioutil.ReadAll(r.Body)

	repos := &Repos{}

	yaml.Unmarshal(repoString, repos)

	imageToLatestVersion := map[string]int{}

	for _, repo := range repos.Repositories {
		tagsListUrl := dockerUrl + "/" + apiVersion + "/" + repo + "/tags/list"
		r, err := http.Get(tagsListUrl)
		if err != nil {
			panic(err)
		}
		tagsString, err := ioutil.ReadAll(r.Body)
		tags := &Tags{}
		yaml.Unmarshal(tagsString, tags)

		highestVersion := 0
		for _, tag := range tags.Tags {
			version, err := strconv.Atoi(tag)
			if err != nil {
				continue
			}
			if version > highestVersion {
				highestVersion = version
			}
		}

		imageToLatestVersion[dockerServerAddress+"/"+repo] = highestVersion

	}

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	list, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	var upgradeCandidates []upgradeCandidate

	fmt.Printf("Listing deployments eligible for upgrade in namespace %q:\n", apiv1.NamespaceDefault)
	for _, d := range list.Items {
		image := d.Spec.Template.Spec.Containers[0].Image

		if !strings.Contains(image, dockerServerAddress) {
			continue
		}

		imageSplit := strings.Split(image, ":")
		untaggedImage := imageSplit[0] + ":" + imageSplit[1]

		if len(imageSplit) < 3 {
			log.Printf("fff %v", 3)
		}

		deployedVersion, err := strconv.Atoi(imageSplit[2])
		if err != nil {
			panic(err)
		}

		latestDockerVersion := imageToLatestVersion[untaggedImage]

		if latestDockerVersion > deployedVersion {
			upgradeCandidates = append(upgradeCandidates, upgradeCandidate{
				deployment: d,
				newImage:   untaggedImage + ":" + strconv.Itoa(latestDockerVersion),
			})
			fmt.Printf(" %v image: %v deployed version: %v   latest version: %v\n", d.Name, untaggedImage, deployedVersion, latestDockerVersion)
		}
	}

	if askForConfirmation("Upgrade all deployments listed above?") {
		for _, c := range upgradeCandidates {
			c.deployment.Spec.Template.Spec.Containers[0].Image = c.newImage
			_, err := deploymentsClient.Update(&c.deployment)
			if err != nil {
				panic(fmt.Errorf("failed to update deployment %v: %w", c.deployment.Name, err))
			}
		}
	}

}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}