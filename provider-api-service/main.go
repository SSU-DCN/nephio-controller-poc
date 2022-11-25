package main

import (
	"fmt"
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
		var kubeConfigRaw = KubeConfigMessage{Name: clusterName, KubeConfig: string(stdout)}
		jsongetClusterResult, errorConvertJson := json.Marshal(kubeConfigRaw)
		if errorConvertJson != nil {
			fmt.Println("error when convert JSON", jsongetClusterResult, errorConvertJson)

		}

		w.Write([]byte(string(stdout)))
	})

	http.ListenAndServe(serverPort, r)
}
