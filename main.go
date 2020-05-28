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
	"dofus-bot/session"
)

var (
	restPosition    models.Pos
	resources       []*models.Resource
	sessionModified bool
	mutex           sync.Mutex

	opts struct {
		Debug  bool `short:"d" long:"debug" description:"Run in debug mode"`
		Wait   bool `short:"w" long:"wait" description:"Simulate short reaction time"`
		NoRest bool `short:"r" long:"no-rest" description:"Disable rest"`
	}
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	if opts.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Info("Starting bot, press '+' to add resource, press 'END' to quit...")

	resources, restPosition = session.Select()

	// disable rest position if asked
	if opts.NoRest {
		restPosition = models.Pos{}
	}

	// focus window
	robotgo.MoveClick(150, 150, "left")

	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())

	// start listening for resource monitoring
	wg.Add(1)
	go listenResourceRegistration(cancel, wg)

	// start resource supervision
	wg.Add(1)
	go resourceSupervision(ctx, wg)

	wg.Wait()

	// save session
	// if sessionModified {
	session.Save(&restPosition, resources)
	// }

	logrus.Info("see you soon!")
}

func listenResourceRegistration(cancel context.CancelFunc, wg *sync.WaitGroup) {
	s := robotgo.Start()
	defer robotgo.End()
	defer cancel()
	defer wg.Done()

	for ev := range s {
		switch ev.Kind {
		case hook.KeyUp:
			switch ev.Rawcode {
			// quit
			case 65367:
				return

			// +
			case 61, 65323:
				addResource()

			default:
				logrus.Debugf("event: %s", strings.SplitAfter(ev.String(), "Event: ")[1])
			}
		case hook.MouseUp:
			switch ev.Button {
			case 2:
				addRestPosition()
			}
		}

	}
}

func addRestPosition() {
	mutex.Lock()
	defer mutex.Unlock()

	if !opts.NoRest {
		sessionModified = true
	}

	x, y := robotgo.GetMousePos()
	restPosition = models.Pos{X: x, Y: y}
	logrus.Infof("set rest position at [%vx%v]", x, y)
}

func addResource() {
	mutex.Lock()
	defer mutex.Unlock()

	sessionModified = true
	resources = append(resources, models.NewResourceUnderMouse())
}

func resourceSupervision(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	defer wg.Done()

	var collecting bool
	var resting bool
	var nextResource *models.Resource

	for {
		select {
		case <-ticker.C:
			justFinish := false

			if collecting {
				if collecting = nextResource.IsActive(); !collecting {
					logrus.Infof("[%s] collected", nextResource)
					justFinish = true
				}
			}

			if !collecting {
				var lastPos models.Pos
				if !justFinish && !restPosition.IsNull() {
					lastPos = restPosition
				} else if nextResource != nil {
					lastPos = nextResource.Pos
				}

				if nextResource, collecting = nearestRessource(lastPos); collecting {
					logrus.Infof("collecting [%s]...", nextResource)
					resting = false
					go nextResource.Collect(opts.Wait && !justFinish)
				} else if !resting && !restPosition.IsNull() {
					resting = true
					goRest()
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

func goRest() {
	logrus.Info("going to rest position")
	robotgo.MoveClick(restPosition.X, restPosition.Y, "left")
	time.Sleep(time.Millisecond * 20)
	robotgo.MoveMouse(150, 150)
}

func nearestRessource(lastPos models.Pos) (*models.Resource, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	var nextClosestResource *models.Resource
	var bestDistance int = math.MaxInt64

	for _, resource := range resources {
		// ignore inactive resources
		if !resource.IsActive() {
			continue
		}

		// if last resource is empty, then when cant do comparisons and return the first active resource
		if lastPos.IsNull() {
			return resource, true
		}

		// else, we compare distance to get the closest
		if dist := lastPos.DistanceTo(resource.Pos); dist <= bestDistance {
			nextClosestResource = resource
			bestDistance = dist
		}
	}

	return nextClosestResource, nextClosestResource != nil
}
