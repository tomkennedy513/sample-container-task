package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	v1 "gitlab.eng.vmware.com/one-tanzu/container-app-resource/api/capps/v1"
	"sigs.k8s.io/yaml"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) == 1 {
		log.Println("No args provided, skipping adding env vars")
		return
	}

	workspaceDir := os.Getenv("TANZU_BUILD_WORKSPACE_DIR")
	containerAppPath := filepath.Join(workspaceDir, "output", "containerapp.yml")

	file, err := os.ReadFile(containerAppPath)
	if err != nil {
		log.Fatal(err)
	}

	var containerApp v1.ContainerApp
	err = yaml.Unmarshal(file, &containerApp)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("processing %s", containerApp.Name)

	for _, arg := range os.Args[1:] {
		split := strings.Split(arg, "=")
		if len(split) != 2 {
			log.Fatalf("invalid arg: %s. expected argument in format `NAME=VALUE`", arg)
		}
		name, value := split[0], split[1]

		log.Printf("will add env %s=%s", name, value)

		containerApp.Spec.NonSecretEnv = append(
			containerApp.Spec.NonSecretEnv,
			v1.NonSecretEnvVar{Name: name, Value: value},
		)

	}

	bytes, err := yaml.Marshal(containerApp)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(containerAppPath, bytes, 0666)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Done.")
}
