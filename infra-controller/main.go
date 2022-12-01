package main

import (
	"fmt"
	"io/ioutil"

	// "container/list"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Define listening port
const serverPort string = ":3333"
const kubectlCmd string = "kubectl"
const clusterctlCmd string = "clusterctl"

type KubeConfigList struct {
	Items []KubeConfig
}
type KubeConfig struct {
	Name       string `json:"Name"`
	KubeConfig string `json:"KubeConfig"`
}
type Message struct {
	Namespace string `json:"Namespace,omitempty"`
	Name      string `json:"Name,omitempty"`
	Phase     string `json:"Phase,omitempty"`
	Age       string `json:"Age,omitempty"`
}
type ClusterConfigurations struct {
	ClusterName              string `json:"ClusterName"`
	KubernetesVersion        string `json:"KubernetesVersion"`
	ControlPlaneMachineCount string `json:"ControlPlaneMachineCount"`
	KubernetesMachineCount   string `json:"KubernetesMachineCount"`
}
type ClusterRecord struct {
	Name                     string            `json:"name,omitempty"`
	InfraType                string            `json:"infraType,omitempty"`
	Labels                   map[string]string `json:"labels,omitempty"`
	Repository               string            `json:"repository,omitempty"`
	Provider                 string            `json:"provider,omitempty"`
	ProvisionMethod          string            `json:"provisionMethod,omitempty"`
	Namespace                string            `json:"namespace,omitempty"`
	KubernetesVersion        string            `json:"pubernetesVersion,omitempty"`
	ControlPlaneMachineCount string            `json:"controlPlaneMachineCount,omitempty"`
	KubernetesMachineCount   string            `json:"kubernetesMachineCount,omitempty"`
}
type ClusterRecordList struct {
	Items []ClusterRecord
}

type InfraRecord struct {
	Name                     string            `json:"name,omitempty"`
	InfraType                string            `json:"infraType,omitempty"`
	Labels                   map[string]string `json:"labels,omitempty"`
	Provider                 string            `json:"provider,omitempty"`
	ProvisionMethod          string            `json:"provisionMethod,omitempty"`
	Namespace                string            `json:"namespace,omitempty"`
	KubernetesVersion        string            `json:"kubernetesVersion,omitempty"`
	ControlPlaneMachineCount string            `json:"controlPlaneMachineCount,omitempty"`
	KubernetesMachineCount   string            `json:"kubernetesMachineCount,omitempty"`
}
type InfraRecordList struct {
	Items []InfraRecord
}

var kubeConfigList KubeConfigList
var currentClusterDeploymentPackages, backupClusterDeploymentPackages ClusterRecordList
var currentInfraDeploymentPackages, backupInfraDeploymentPackages InfraRecordList

func main() {
	// currentListCluster := list.newList()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	fmt.Println("KubeConfig file path" + os.Getenv("KUBECONFIG"))
	r.Get("/getcluster", func(w http.ResponseWriter, r *http.Request) {

		prg := "kubectl"
		arg1 := "get"
		arg2 := "cluster"
		arg3 := "-A"
		cmd := exec.Command(prg, arg1, arg2, arg3)
		stdout, err := cmd.Output()

		if err != nil {
			fmt.Println(err.Error())
			log.Fatal(err)
			return
		}

		var getClusterResult []Message
		trimmedString := strings.TrimSpace(string(stdout))
		listTrimmedString := strings.Split(trimmedString, "\n")

		for i, str := range listTrimmedString {
			if i != 0 {
				splitStr := strings.Fields(str)
				msg := Message{splitStr[0], splitStr[1], splitStr[2], splitStr[3]}
				msgMarshaled, _ := json.Marshal(msg)
				fmt.Println("msgMarshaled", string(msgMarshaled))
				getClusterResult = append(getClusterResult, msg)
			}
		}
		jsongetClusterResult, errorConvertJson := json.Marshal(getClusterResult)
		if errorConvertJson != nil {
			fmt.Println("error:", errorConvertJson)
		}

		w.Write([]byte(string(jsongetClusterResult)))
	})

	r.Get("/getkubeadmcontrolplanes", func(w http.ResponseWriter, r *http.Request) {

		prg := "kubectl"
		arg1 := "get"
		arg2 := "kubeadmcontrolplane"
		arg3 := "-A"
		cmd := exec.Command(prg, arg1, arg2, arg3)
		// Get the result from kubectl and send to Infra Controller
		stdout, err := cmd.Output()

		if err != nil {
			fmt.Println(err.Error())
			log.Fatal(err)
			return
		}

		var getClusterResult []Message
		trimmedString := strings.TrimSpace(string(stdout))
		listTrimmedString := strings.Split(trimmedString, "\n")

		for i, str := range listTrimmedString {
			if i != 0 {
				splitStr := strings.Fields(str)
				msg := Message{splitStr[0], splitStr[1], splitStr[2], splitStr[3]}
				// msgMarshaled, _ := json.Marshal(msg)

				getClusterResult = append(getClusterResult, msg)
			}
		}
		jsongetClusterResult, errorConvertJson := json.Marshal(getClusterResult)
		if errorConvertJson != nil {
			fmt.Println("error:", errorConvertJson)
		}

		w.Write([]byte(string(jsongetClusterResult)))
	})

	r.Get("/getkubeconfig", func(w http.ResponseWriter, r *http.Request) {
		var clusterName string
		clusterName = r.Header.Get("clustername")
		if len(clusterName) < 1 {
			fmt.Println("Missing clustername field in request")
		}
		prg := "clusterctl"
		arg1 := "get"
		arg2 := "kubeconfig"
		cmd := exec.Command(prg, arg1, arg2, clusterName)
		// Get the result from kubectl and send to Infra Controller
		stdout, err := cmd.Output()

		if err != nil {
			fmt.Println(err.Error())
			log.Fatal(err)
			return
		}
		var kubeConfigRaw = KubeConfig{Name: clusterName, KubeConfig: string(stdout)}
		jsongetClusterResult, errorConvertJson := json.Marshal(kubeConfigRaw)
		if errorConvertJson != nil {
			fmt.Println("error when convert JSON", jsongetClusterResult, errorConvertJson)

		}

		w.Write([]byte(string(stdout)))
	})

	r.Post("/createNewCluster", func(w http.ResponseWriter, r *http.Request) {

		// defer r.Body.Close()

		httpPostBody, err := ioutil.ReadAll(r.Body) //<--- here!

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(httpPostBody))
		var clusterConfig ClusterConfigurations
		err = json.Unmarshal(httpPostBody, &clusterConfig)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println((clusterConfig))

		// prg := "echo " + httpPostBody
		// arg := " | kubectl apply -f -"
		// cmd := exec.Command(prg, arg)
		// stdout, err := cmd.Output()

		// if err != nil {
		// 	fmt.Println(err.Error())
		// 	log.Fatal(err)
		// 	return
		// }
		w.Write([]byte(string("stdout")))
	})

	r.Post("/generateNewCluster", func(w http.ResponseWriter, r *http.Request) {

		var httpPostBody string = string("Test")

		prg := "echo " + httpPostBody
		arg := " | kubectl apply -f -"
		cmd := exec.Command(prg, arg)
		stdout, err := cmd.Output()

		if err != nil {
			fmt.Println(err.Error())
			log.Fatal(err)
			return
		}
		w.Write([]byte(string(stdout)))
	})

	r.Post("/updateClusterPackage", func(w http.ResponseWriter, r *http.Request) {
		// defer r.Body.Close()

		httpPostBody, err := ioutil.ReadAll(r.Body) //<--- here!

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(httpPostBody))
		// var packageDeployment ClusterRecord
		// err = json.Unmarshal(httpPostBody, &packageDeployment)

		// if err != nil {
		// 	fmt.Println(err)
		// }

		// prg := "echo " + httpPostBody
		// arg := " | kubectl apply -f -"
		// cmd := exec.Command(prg, arg)
		// stdout, err := cmd.Output()

		// if err != nil {
		// 	fmt.Println(err.Error())
		// 	log.Fatal(err)
		// 	return
		// }
		w.Write([]byte("received Cluster Package"))
	})

	r.Post("/updateInfraPackage", func(w http.ResponseWriter, r *http.Request) {
		// defer r.Body.Close()

		httpPostBody, err := ioutil.ReadAll(r.Body) //<--- here!

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(httpPostBody))
		// var packageDeployment ClusterRecord
		// err = json.Unmarshal(httpPostBody, &packageDeployment)

		// if err != nil {
		// 	fmt.Println(err)
		// }

		// prg := "echo " + httpPostBody
		// arg := " | kubectl apply -f -"
		// cmd := exec.Command(prg, arg)
		// stdout, err := cmd.Output()

		// if err != nil {
		// 	fmt.Println(err.Error())
		// 	log.Fatal(err)
		// 	return
		// }
		w.Write([]byte("received Infra Package"))
	})

	http.ListenAndServe(serverPort, r)
}

func isExistingInClusterList(list *ClusterRecordList, item *ClusterRecord) bool {
	if len((*list).Items) < 1 {
		return false
	}
	for _, iter := range (*list).Items {
		if iter.Name == item.Name {
			return true
		}
	}
	return false
}

func isExistingInInfraList(list *ClusterRecordList, item *ClusterRecord) bool {
	if len((*list).Items) < 1 {
		return false
	}
	for _, iter := range (*list).Items {
		if iter.Name == item.Name {
			return true
		}
	}
	return false
}
