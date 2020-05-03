package main

import (
	"context"
	"math"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/jessevdk/go-flags"
	hook "github.com/robotn/gohook"
	"github.com/sirupsen/logrus"

	"dofus-bot/models"
)

var resources map[string]models.Resource
var mutex sync.Mutex

var opts struct {
	Debug  bool `short:"d" long:"debug" description:"Run in debug mode"`
	Invert bool `short:"i" long:"invert" description:"Invert resource selection"`
	React  bool `short:"r" long:"react" description:"Simulate short reaction time"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	if opts.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Info("Starting bot, press '+' to add resource, press 'END' to quit...")
	if opts.Invert {
		logrus.Warn("Running in inverted mode")
	}

	resources = make(map[string]models.Resource)

	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())

	// start listening for resource monitoring
	wg.Add(1)
	go listenResourceRegistration(cancel, wg)

	// start resource supervision
	wg.Add(1)
	go resourceSupervision(ctx, wg)

	wg.Wait()
	logrus.Info("see you soon!")
}

func listenResourceRegistration(cancel context.CancelFunc, wg *sync.WaitGroup) {
	s := robotgo.Start()
	defer robotgo.End()
	defer cancel()
	defer wg.Done()

	// var lastResourceAdded models.Resource

	for ev := range s {
		switch ev.Kind {
		case hook.KeyUp:
			switch ev.Rawcode {
			// quit
			case 65367:
				return

			// +
			case 61, 65323:
				// lastResourceAdded = addResource()
				addResource()

			// z, - (cancel last add)
			// case 122, 45:
			// removeResource(lastResourceAdded.ID)

			default:
				logrus.Debugf("event: %s", strings.SplitAfter(ev.String(), "Event: ")[1])
			}
		}
	}
}

func addResource() models.Resource {
	mutex.Lock()
	defer mutex.Unlock()

	addedResource := models.NewResourceUnderMouse(opts.Invert)
	resources[addedResource.ID] = addedResource
	return addedResource
}

func removeResource(lastResourceAddedID string) {
	mutex.Lock()
	defer mutex.Unlock()

	if lastResourceAddedID != "" {
		logrus.Warnf("[%s] removed", lastResourceAddedID)
		delete(resources, lastResourceAddedID)
	}
}

func resourceSupervision(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	defer wg.Done()

	var collecting bool
	var nextResource models.Resource

	for {
		select {
		case <-ticker.C:
			if _, exists := resources[nextResource.ID]; !exists && nextResource.ID != "" && collecting {
				logrus.Infof("stop collecting [%s]", nextResource.ID)
				collecting = false
			}

			if collecting {
				if collecting = nextResource.IsActive(); !collecting {
					logrus.Infof("[%s] collected", nextResource.ID)
				}
			}

			if !collecting {
				if nextResource, collecting = nearestRessource(nextResource); collecting {
					logrus.Infof("collecting [%s]...", nextResource.ID)
					nextResource.Collect(opts.React)
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

func nearestRessource(previousResource models.Resource) (models.Resource, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	var nextClosestResource models.Resource

	bestDistance := math.Inf(0)
	for _, resource := range resources {
		// ignore inactive resources
		if !resource.IsActive() {
			continue
		}

		// if last resource is empty, then when cant do comparisons and return the first active resource
		if previousResource.ID == "" {
			return resource, true
		}

		// else, we compare distance to get the closest
		if dist := previousResource.SquareDistanceTo(resource); dist < bestDistance {
			nextClosestResource = resource
			bestDistance = dist
		}
	}

	return nextClosestResource, nextClosestResource.ID != ""
}
