package main

import (
	"fmt"
	"io/ioutil"
	"time"

	// work "github.com/gocraft/work"
	// "container/list"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Define listening port
const serverPort string = ":3333"
const kubectlCmd string = "kubectl"
const clusterctlCmd string = "clusterctl"

var kubeConfig string

type KubeConfigMessage struct {
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
	CreatedTime              time.Time         `json:"createdTime,omitempty"`
	UpdatedTime              time.Time         `json:"updatedTime,omitempty"`
}

var listYamlFileClusterAPI []string

func main() {
	// currentListCluster := list.newList()
	kubeConfig = getEnv("KUBECONFIG", "$HOME/.kube/config")
	fmt.Println("Env KUBECONFIG", kubeConfig)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// fmt.Println("KubeConfig file path" + os.Getenv("KUBECONFIG"))
	r.Get("/getcluster", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received Get Cluster Request")
		prg := "kubectl"
		arg1 := "get"
		arg2 := "cluster"
		arg3 := "-A"
		argKubeConfig := "--kubeconfig"
		cmd := exec.Command(prg, arg1, arg2, arg3, argKubeConfig, kubeConfig)
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
		argKubeConfig := "--kubeconfig"
		cmd := exec.Command(prg, arg1, arg2, arg3, argKubeConfig, kubeConfig)
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
		// argKubeConfig := "--kubeconfig " + kubeConfig
		cmd := exec.Command(prg, arg1, arg2, clusterName)
		// Get the result from kubectl and send to Infra Controller
		stdout, err := cmd.Output()

		if err != nil {
			fmt.Println(err.Error())
			log.Fatal(err)
			return
		}
		var kubeConfigRaw = KubeConfigMessage{Name: clusterName, KubeConfig: string(stdout)}
		jsongetClusterResult, errorConvertJson := json.Marshal(kubeConfigRaw)
		if errorConvertJson != nil {
			fmt.Println("error when convert JSON", jsongetClusterResult, errorConvertJson)

		}

		w.Write([]byte(string(stdout)))
	})

	r.Post("/createNewCluster", func(w http.ResponseWriter, r *http.Request) {

		// defer r.Body.Close()
		fmt.Println("Received create new Cluster Request")
		httpPostBody, err := ioutil.ReadAll(r.Body) //<--- here!

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(httpPostBody))
		var clusterConfig ClusterRecord
		err = json.Unmarshal(httpPostBody, &clusterConfig)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println((clusterConfig))
		fmt.Println("Before Applying cluster YAML FIle")
		clusterYamlFile, ok := generateClusterYamlFile(clusterConfig)
		if ok {
			prg := "kubectl"
			arg1 := "apply -f"
			argKubeConfig := "--kubeconfig " + kubeConfig
			cmd := exec.Command(prg, arg1, clusterYamlFile, argKubeConfig)
			// Get the result from kubectl and send to Infra Controller
			stdout1, err := cmd.Output()

			if err != nil {
				fmt.Println(err.Error())
				log.Fatal(err)
			}
			fmt.Println("Output kubectl apply -f ", string(stdout1))
			listYamlFileClusterAPI = append(listYamlFileClusterAPI, clusterYamlFile)
		}

		// prg := "echo " + httpPostBody
		// arg := " | kubectl apply -f -"
		// cmd := exec.Command(prg, arg)
		// stdout, err := cmd.Output()

		// if err != nil {
		// 	fmt.Println(err.Error())
		// 	log.Fatal(err)
		// 	return
		// }
		w.Write([]byte(string("Creating cluster: ") + clusterConfig.Name))
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
	fmt.Println("Start Server at port", serverPort)
	http.ListenAndServe(serverPort, r)
}

// ==============================FUNCTIONS============================
func getEnv(key string, defaultValue string) string {
	fmt.Println("Get Env KUBECONFIG", os.Getenv("KUBECONFIG"))
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}

func createTempFolder(nameCluster string) string {
	dname, err := os.MkdirTemp("", nameCluster)
	if err != nil {
		panic(err)
	}
	return dname
}
func generateClusterYamlFile(record ClusterRecord) (string, bool) {
	fmt.Println("Generate Cluster Yaml File", record.Name)
	arg := "clusterctl"
	arg1 := "generate"
	arg2 := "cluster"
	// clusterctl generate cluster capi-quickstart   --kubernetes-version v1.23.5   --control-plane-machine-count=3   --worker-machine-count=3   > capi-quickstart.yaml
	argK8sVersion := "--kubernetes-version v1.23.5"
	argControlPlaneMachineCount := "--control-plane-machine-count=" + record.ControlPlaneMachineCount
	argWorkerMachineCount := "--worker-machine-count=" + record.KubernetesMachineCount
	cmd := exec.Command(arg, arg1, arg2, record.Name, argK8sVersion, argControlPlaneMachineCount, argWorkerMachineCount)
	stdout, err := cmd.Output()
	fmt.Println("stdout", string(stdout))
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
		return string(stdout), false
	}
	// Create folder
	tempFolder := createTempFolder(record.Name)
	templateClusterFile := filepath.Join(tempFolder, record.Name)
	err = os.WriteFile(templateClusterFile, stdout, 0777)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
		return "error", false
	}

	return templateClusterFile, true
}
