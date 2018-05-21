package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// global flags
var hostFlag = flag.String("hostname", "", "hostname")
var typeFlag = flag.String("type", "", "type: roles | vars")
var oldFileFlag = flag.String("old", "", "old file to diff")
var newFileFlag = flag.String("new", "", "new file to diff")

func main() {
	var rolesToRun []string
	flag.Parse()

	oldLines := getLines(*oldFileFlag)
	newLines := getLines(*newFileFlag)

	var purgeSuffix string
	if *typeFlag == "roles" {
		purgeSuffix = "_purge"
	}

	// check for removed/changed lines
	for line := range oldLines {
		if _, ok := newLines[line]; !ok {
			rolesToRun = append(rolesToRun, fmt.Sprintf("%s%s", line, purgeSuffix))
		}
	}

	// check for added/changed lines
	for line := range newLines {
		if _, ok := oldLines[line]; !ok {
			rolesToRun = append(rolesToRun, line)
		}
	}

	rolesToRun = trimAndUniq(rolesToRun)
	log.Printf(">>> [command]: roles to run: %v\n", rolesToRun)
	for _, role := range rolesToRun {
		log.Printf(">>> [command]: executing ansible for role: %s\n", role)
		runAnsible(role)
	}

	os.Rename(*newFileFlag, *oldFileFlag)
}

func getLines(fileName string) map[string]bool {
	f, err := os.Open(fileName)

	// TODO: what if old file doesn't exist yet
	if err != nil {
		log.Fatalf("could not open file %s: %v", fileName, err)
	}

	retMap := make(map[string]bool)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		retMap[scanner.Text()] = true
	}

	return retMap
}

func runAnsible(role string) error {

	var ansibleDir string = "/demo/ansible"
	var ansiblePlaybook string = "/opt/ansible/ansible/bin/ansible-playbook"
	var inventory string = "inventory"
	var playbook string = "playbook.yml"
	var roleToRun string = fmt.Sprintf("ansible_role_to_run=%s", role)
	var command string

	command = ansiblePlaybook
	args := []string{"-i", inventory, "-l", *hostFlag, "-e", roleToRun, playbook}

	cmd := exec.Command(command, args...)
	cmd.Dir = ansibleDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf(">>> [command]: %s %s\n", command, strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
	log.Printf(">>> [command]: done\n")

	return nil
}

func trimAndUniq(s []string) []string {
	result := []string{}
	encountered := map[string]bool{}
	for _, v := range s {
		v := strings.Split(v, ".")[0]
		if _, ok := encountered[v]; !ok {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}
