package session

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"dofus-bot/models"
)

type session struct {
	RestPos   *models.Pos       `json:"restPosition,omitempty"`
	Resources []models.Resource `json:"resources"`
}

const (
	sessionFile string = "sessions.json"
)

var (
	loaded           bool
	sessions         map[string]session
	selectedSessions []string
)

func readSessions() bool {
	sessions = make(map[string]session)

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

func Select() ([]models.Resource, models.Pos) {
	resources := make([]models.Resource, 0)
	restPos := models.Pos{}

	if !loaded && !readSessions() {
		return resources, restPos
	}

	if len(sessions) > 0 {
		logrus.Info("Saved sessions are:")

		sessionNames := make(sort.StringSlice, 0)
		for name := range sessions {
			sessionNames = append(sessionNames, name)
		}
		sort.Sort(sessionNames)

		idxToSession := []string{}
		for n, name := range sessionNames {
			fmt.Printf("%3d- %s (%d)\n", n+1, name, len(sessions[name].Resources))
			idxToSession = append(idxToSession, name)
		}

		fmt.Print("Session(s) to load, separated with comma (Left empty to quit):\n> ")
		buffer := ""
		fmt.Scanf("%s", &buffer)

		for _, choice := range strings.Split(buffer, ",") {
			if idx, _ := strconv.Atoi(choice); idx > 0 && idx <= len(sessions) {
				session := idxToSession[idx-1]
				logrus.Infof("loading session [%s]", session)
				selectedSessions = append(selectedSessions, session)
				resources = append(resources, sessions[session].Resources...)
			}
		}
	}

	if len(resources) < 1 {
		logrus.Info("no session loaded")
	} else if len(selectedSessions) > 1 {
		logrus.Warnf("more than 1 session (%d) loaded, no save available", len(selectedSessions))
	} else if pos := sessions[selectedSessions[0]].RestPos; pos != nil {
		restPos = *pos
	}

	return resources, restPos
}

func Save(restPosition *models.Pos, resources []models.Resource) {
	if !loaded && !readSessions() {
		return
	}

	// get session name
	selectedSession := ""
	switch len(selectedSessions) {
	case 0:
		fmt.Print("Save current session (Type new name or left empty to quit):\n> ")
		fmt.Scanln(&selectedSession)
	case 1:
		selectedSession = selectedSessions[0]
	default:
		logrus.Info("impossible to save multi-sessions")
		return
	}

	if selectedSession == "" {
		return
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
	sessions[selectedSession] = session{
		RestPos:   restPosition,
		Resources: resources,
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
