package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

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
const serverPort string = ":3334"
const kubectlCmd string = "kubectl"
const clusterctlCmd string = "clusterctl"

var providerApiServiceUrl string

// "http://127.0.0.1:3333/createNewCluster"
const CreateNewClusterEndpoint string = "/createNewCluster"

type KubeConfigList struct {
	Items []KubeConfig
}
type KubeConfig struct {
	Name        string    `json:"Name"`
	KubeConfig  string    `json:"KubeConfig"`
	CreatedTime time.Time `json:"CreatedTime"`
	UpdatedTime time.Time `json:"UpdatedTime"`
}
type Message struct {
	Namespace string `json:"Namespace,omitempty"`
	Name      string `json:"Name,omitempty"`
	Phase     string `json:"Phase,omitempty"`
	Age       string `json:"Age,omitempty"`
}
type ClusterConfigurations struct {
	Name       string            `json:"name"`
	Type       string            `json:"infraType"`
	Labels     map[string]string `json:"labels"`
	Repository string            `json:"repository"`
}
type ClusterConfigurationsList struct {
	Items []ClusterConfigurations
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
	CreatedTime              time.Time         `json:"createdTime,omitempty"`
	UpdatedTime              time.Time         `json:"updatedTime,omitempty"`
}
type InfraRecordList struct {
	Items []InfraRecord
}

var kubeConfigList KubeConfigList
var currentClusterDeploymentPackagesList, backupClusterDeploymentPackagesList ClusterRecordList
var currentInfraDeploymentPackagesList, backupInfraDeploymentPackagesList InfraRecordList
var currentClusterConfigList, backupClusterConfigsList ClusterConfigurationsList

func main() {
	// currentListCluster := list.newList()
	providerApiServiceUrl := "http://provider-api-service:3333"
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	fmt.Println("KubeConfig file path" + os.Getenv("KUBECONFIG"))
	fmt.Println("Print PROVIDER_API_SVC_SERVICE_HOST: ", providerApiServiceUrl)

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(string("Received Request Infra COntroller")))
	})
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
	r.Get("/testSendPackage", func(w http.ResponseWriter, r *http.Request) {
		config := ClusterRecord{
			"default", "minimal", map[string]string{"none": "none"}, "default", "default", "default", "default", "v1.24.0", "1", "1", time.Now(), time.Now(),
		}
		go sendRequestCreateNewCluster(config, providerApiServiceUrl)
		w.Write([]byte(string("received")))
	})
	r.Post("/updateClusterPackage", func(w http.ResponseWriter, r *http.Request) {
		// defer r.Body.Close()

		httpPostBody, err := ioutil.ReadAll(r.Body) //<--- here!

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(httpPostBody))

		var clusterConfigs ClusterConfigurations
		err = json.Unmarshal(httpPostBody, &clusterConfigs)
		// clusterConfigs.CreatedTime = time.Now()
		// clusterConfigs.UpdatedTime = clusterConfigs.CreatedTime
		fmt.Println("Print Cluster Configurations", clusterConfigs)
		// Add to List if not exist
		if !isExistingInClusterConfigList(&currentClusterConfigList, &clusterConfigs) {
			currentClusterConfigList.Items = append(currentClusterConfigList.Items, clusterConfigs)
			config := mappingValueofClusterToClusterRecord(clusterConfigs, currentInfraDeploymentPackagesList)
			go sendRequestCreateNewCluster(config, providerApiServiceUrl)
		} else {
			replaceExistingClusterConfiguration(&currentClusterConfigList, clusterConfigs)
		}
		// if ( currentClusterConfigList)
		// if (isExistingInInfraList)
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
		var infraDeployment InfraRecord
		err = json.Unmarshal(httpPostBody, &infraDeployment)
		fmt.Println("Print infraDeployment", infraDeployment)
		if !isExistingInInfraList(&currentInfraDeploymentPackagesList, &infraDeployment) {
			currentInfraDeploymentPackagesList.Items = append(currentInfraDeploymentPackagesList.Items, infraDeployment)
		} else {
			// Replace
			replaceExistingInfraPackage(&currentInfraDeploymentPackagesList, infraDeployment)
		}
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
func replaceExistingClusterPackage(list *InfraRecordList, item InfraRecord) bool {
	if len((*list).Items) < 1 {
		return false
	}
	for i, iter := range (*list).Items {
		if iter.Name == item.Name {
			(*list).Items[i] = item
			return true
		}
	}
	return false
}

func isExistingInInfraList(list *InfraRecordList, item *InfraRecord) bool {
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

func replaceExistingInfraPackage(list *InfraRecordList, item InfraRecord) bool {
	if len((*list).Items) < 1 {
		return false
	}
	for i, iter := range (*list).Items {
		if iter.Name == item.Name {
			(*list).Items[i] = item
			return true
		}
	}
	return false
}

func isExistingInClusterConfigList(list *ClusterConfigurationsList, item *ClusterConfigurations) bool {
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

func replaceExistingClusterConfiguration(list *ClusterConfigurationsList, item ClusterConfigurations) bool {
	if len((*list).Items) < 1 {
		return false
	}
	for i, iter := range (*list).Items {
		if iter.Name == item.Name {
			(*list).Items[i] = item
			return true
		}
	}
	return false
}
func delete_at_index(slice []any, index int) []any {

	// Append function used to append elements to a slice
	// first parameter as the slice to which the elements
	// are to be added/appended second parameter is the
	// element(s) to be appended into the slice
	// return value as a slice
	return append(slice[:index], slice[index+1:]...)
}
func delete_at_index_ClusterConfiguration(slice []ClusterConfigurations, index int) []ClusterConfigurations {

	// Append function used to append elements to a slice
	// first parameter as the slice to which the elements
	// are to be added/appended second parameter is the
	// element(s) to be appended into the slice
	// return value as a slice
	return append(slice[:index], slice[index+1:]...)
}

func findAndDeleteClusterPackage(list ClusterConfigurationsList, itemToDelete ClusterConfigurations) ClusterConfigurationsList {

	for index, item := range list.Items {
		if item.Name != itemToDelete.Name {
			var newClusterConfigurationsList ClusterConfigurationsList

			newClusterConfigurationsList.Items = delete_at_index_ClusterConfiguration(list.Items, index)
			return newClusterConfigurationsList
		}
	}
	return list
}
func delete_at_index_InfraRecord(slice []InfraRecord, index int) []InfraRecord {

	// Append function used to append elements to a slice
	// first parameter as the slice to which the elements
	// are to be added/appended second parameter is the
	// element(s) to be appended into the slice
	// return value as a slice
	return append(slice[:index], slice[index+1:]...)
}

func findAndDeleteInfraPackage(list InfraRecordList, itemToDelete ClusterConfigurations) InfraRecordList {

	for index, item := range list.Items {
		if item.Name != itemToDelete.Name {
			var newInfraRecordList InfraRecordList

			newInfraRecordList.Items = delete_at_index_InfraRecord(list.Items, index)
			return newInfraRecordList
		}
	}
	return list
}

func mappingValueofClusterToClusterRecord(newCluster ClusterConfigurations, listInfra InfraRecordList) ClusterRecord {
	var newClusterRecord ClusterRecord
	// Copy value to new variable
	newClusterRecord.Name = newCluster.Name
	newClusterRecord.InfraType = newCluster.Type
	newClusterRecord.Labels = newCluster.Labels
	newClusterRecord.Repository = newCluster.Repository
	// mapping value to new cluster record
	// First finding first infra record associate with defined in cluster configuration
	// If not found, use default value
	// Default value:
	// Name                     string            `json:"name,omitempty"`
	// InfraType                string            `json:"infraType,omitempty"`
	// Labels                   map[string]string `json:"labels,omitempty"`
	// Provider                 string            `json:"provider,omitempty"`
	// ProvisionMethod          string            `json:"provisionMethod,omitempty"`
	// Namespace                string            `json:"namespace,omitempty"`
	// KubernetesVersion        string            `json:"kubernetesVersion,omitempty"`
	// ControlPlaneMachineCount string            `json:"controlPlaneMachineCount,omitempty"`
	// KubernetesMachineCount   string            `json:"kubernetesMachineCount,omitempty"`
	mappingInfraRecord := InfraRecord{"default", "minimal", map[string]string{"none": "none"}, "default", "default", "default", "v1.24.0", "1", "1", time.Now(), time.Now()}
	for _, infraItem := range listInfra.Items {
		if newClusterRecord.InfraType == infraItem.Name {
			mappingInfraRecord = infraItem
			break
		}
	}
	// Assign value to newclusterRecord
	newClusterRecord.Provider = mappingInfraRecord.Provider
	newClusterRecord.ProvisionMethod = mappingInfraRecord.ProvisionMethod
	newClusterRecord.ControlPlaneMachineCount = mappingInfraRecord.ControlPlaneMachineCount
	newClusterRecord.KubernetesMachineCount = mappingInfraRecord.KubernetesMachineCount
	newClusterRecord.KubernetesVersion = mappingInfraRecord.KubernetesVersion
	return newClusterRecord
	// (*listClusterRecord).Items = append((*&listClusterRecord).Items, newClusterRecord)
}

func sendRequestCreateNewCluster(config ClusterRecord, url string) bool {
	// Send request
	jsonBytesData, error := json.Marshal(config)
	if error != nil {
		return false
	}
	request, error := http.NewRequest("POST", url+CreateNewClusterEndpoint, bytes.NewBuffer(jsonBytesData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		return false
	}

	fmt.Println("Sent Cluster package. response Status:", response.Status)

	return true
}
func getEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}
