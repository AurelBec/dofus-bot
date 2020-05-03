package session

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	"dofus-bot/models"
)

const sessionFile string = "sessions.json"

var loaded bool
var sessions map[string][]models.Resource
var selectedSession string

func Load() bool {
	sessions = make(map[string][]models.Resource)

	jsonFile, err := os.Open(sessionFile)
	if err != nil {
		logrus.Errorf("failed to load session file [%v]: %v", sessionFile, err)
		return false
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logrus.Errorf("failed to parse session file [%v]: %v", sessionFile, err)
		return false
	}

	err = json.Unmarshal(byteValue, &sessions)
	if err != nil {
		logrus.Errorf("failed to unmarshal graph file [%v] : %v", sessionFile, err)
		return false
	}

	loaded = true
	return true
}

func Select() map[string]models.Resource {
	resources := make(map[string]models.Resource)

	if !loaded && !Load() {
		return resources
	}

	if len(sessions) > 0 {
		logrus.Info("Saved sessions are:")
		idxToSession := []string{}
		for id, session := range sessions {
			fmt.Printf("%3d- %s (%d)\n", len(idxToSession)+1, id, len(session))
			idxToSession = append(idxToSession, id)
		}

		fmt.Print("Which session load (Type 0 or left empty to quit):\n> ")
		idx := 0
		fmt.Scanf("%d", &idx)
		if idx > 0 && idx <= len(sessions) {
			selectedSession = idxToSession[idx-1]
			logrus.Infof("loading session [%s]", selectedSession)
			for _, resource := range sessions[selectedSession] {
				resources[resource.ID] = resource
			}
		}
	}

	if len(resources) < 1 {
		logrus.Info("no session loaded")
	}

	return resources
}

func Save(resources map[string]models.Resource) {
	if !loaded && !Load() {
		return
	}

	// get session name if missing
	if selectedSession == "" {
		fmt.Print("Save current session (Type new name or left empty to quit):\n> ")
		fmt.Scanln(&selectedSession)
		if selectedSession == "" {
			return
		}
	}

	// ask for override
	if _, exists := sessions[selectedSession]; exists {
		fmt.Printf("Override session [%s] (Type Yy(es) to confirm):\n> ", selectedSession)
		answer := ""
		fmt.Scanln(&answer)
		switch answer {
		case "Y", "Yes", "y", "yes":
			// continue
		default:
			return
		}
	}

	// update ressources
	sessions[selectedSession] = []models.Resource{}
	for _, resource := range resources {
		sessions[selectedSession] = append(sessions[selectedSession], resource)
	}

	// saved into file
	sessionsFile, err := os.Create(sessionFile)
	if err != nil {
		logrus.Errorf("failed to open session file [%v]: %v", sessionFile, err)
		return
	}
	defer sessionsFile.Close()

	if err := json.NewEncoder(sessionsFile).Encode(&sessions); err != nil {
		logrus.Errorf("failed to save session file [%v]: %v", sessionFile, err)
		return
	}

	logrus.Infof("session [%s] saved!", selectedSession)
}
