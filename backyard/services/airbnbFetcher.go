package services

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os/exec"
    "strings"

    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"
    "github.com/hospedate/backyard/repositories"
)

const STATUS_CHANNEL_CAPACITY = 3
const EVENTS_CHANNEL_CAPACITY = 50000

type NewPropertyEvent struct {
    Property *models.Property
}

type AirbnbFetcher interface {
    Start()
    Stop()
    GetStatus() models.AirbnbFetcherStatus
    SendPropertyEvent(property *models.Property)
}

type AirbnbFetcherService struct {
    statusChannel  chan models.AirbnbFetcherStatus
    propertiesRepo repositories.PropertiesRepository
    logger         log.Logger
    eventsChannel  chan NewPropertyEvent
    status         models.AirbnbFetcherStatus
    disabled       bool
    pathPrefix     string
}

type AirbnbFetcherScriptResult struct {
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Images      []string `json:"images"`
}

var STATUS_STARTED models.AirbnbFetcherStatus = "STARTED"
var STATUS_STOPPED models.AirbnbFetcherStatus = "STOPPED"

func NewAirbnbFetcher(propertiesRepo repositories.PropertiesRepository, disabled bool, pathPrefix string) *AirbnbFetcherService {
    logger := log.GetOrCreateLogger("AirbnbFetcher", "INFO")

    return &AirbnbFetcherService{
        propertiesRepo: propertiesRepo,
        logger:         logger,
        statusChannel:  make(chan models.AirbnbFetcherStatus, STATUS_CHANNEL_CAPACITY),
        eventsChannel:  make(chan NewPropertyEvent, EVENTS_CHANNEL_CAPACITY),
        status:         STATUS_STOPPED,
        disabled:       disabled,
        pathPrefix:     pathPrefix,
    }
}

func (abs *AirbnbFetcherService) Start() {
    if abs.disabled {
        abs.logger.Info("Skipping Airbnb-fetcher service initialization (disabled = true)")
    } else {
        // prevent initialization more than once
        if abs.status == STATUS_STARTED {
            return
        }
        abs.status = STATUS_STARTED

        go abs.mainLoop() // main loop runs in a separate go routine

        abs.logger.Info("Airbnb-fetcher service initialized")
    }
}

func (abs *AirbnbFetcherService) mainLoop() {
    for { // infinite loop unless stopped
        select {
        case event := <-abs.eventsChannel:
            abs.logger.Infof("Received a Property %v", event.Property.Id)
            abs.runFetcher(event)
        case newStatus := <-abs.statusChannel:
            if newStatus == STATUS_STOPPED {
                abs.status = STATUS_STOPPED
                abs.logger.Info("Stopping!")
                return // go routine dies
            }
        }
    }
}

func (abs *AirbnbFetcherService) Stop() {
    // Sending if no one is listening would block the main thread
    // so we avoid sending a message if the service was already
    // stopped.
    if abs.status == STATUS_STOPPED {
        return
    }
    abs.statusChannel <- STATUS_STOPPED
}

func (abs *AirbnbFetcherService) GetStatus() models.AirbnbFetcherStatus {
    return abs.status
}

func (abs *AirbnbFetcherService) SendPropertyEvent(property *models.Property) {
    // prevent sending to the channel if the service is stopped
    // so we avoid blocking the main thread
    if abs.status == STATUS_STARTED {
        abs.eventsChannel <- NewPropertyEvent{Property: property}
    }
}

func (abs *AirbnbFetcherService) runFetcher(event NewPropertyEvent) {
    property := event.Property
    cmdPath := fmt.Sprintf("%v/airbnb-fetcher/main.py", abs.pathPrefix)

    abs.logger.Infof(
        "Calling airbnb-fetcher for property/room: %v/%v",
        property.Id,
        property.AirbnbRoomId,
    )

    cmd := exec.Command(
        "python",
        cmdPath,
        "--room_id",
        *property.AirbnbRoomId,
    )

    // Get the pipe for the output of the script
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        abs.logger.Error("Error getting the stdout pipe", err.Error())
        return
    }

    // Start the command
    if err := cmd.Start(); err != nil {
        abs.logger.Error("Error starting the airbnb-fetcher command", err.Error())
        return
    }

    // Read the output
    output, err := ioutil.ReadAll(stdout)
    if err != nil {
        abs.logger.Error("Error reading stdout from airbnb-fetcher", err.Error())
        return
    }

    // Wait for the command to finish
    err = cmd.Wait()
    if err != nil {
        abs.logger.Error("Error executing airbnb-fetcher", err.Error())
        return
    }

    if strings.TrimSpace(string(output)) == "{}" {
        abs.logger.Infof(
            "Property %v has an invalid room id (%v)",
            property.Id,
            property.AirbnbRoomId,
        )
        return
    } else {
        abs.logger.Infof("the output was \"%v\"", string(output))
    }

    // Parse the JSON output
    var result AirbnbFetcherScriptResult
    err = json.Unmarshal(output, &result)
    if err != nil {
        abs.logger.Error("Error parsing stdout from airbnb-fetcher", err.Error())
        return
    }

    // Print the output
    abs.logger.Infof(
        "Got title for property %v: \"%v\"",
        property.Id,
        safeSlice(result.Title, 20),
    )
    abs.logger.Infof(
        "Got description for property %v: %v",
        property.Id,
        safeSlice(result.Description, 20),
    )
    abs.logger.Infof("Got %v images for property %v", len(result.Images), property.Id)

    newStatus := "active"
    editRequest := models.PropertyEditRequest{
        Id:          property.Id,
        Title:       &result.Title,
        Description: &result.Description,
        Images:      result.Images,
        Status:      &newStatus,
    }
    _, err = abs.propertiesRepo.Edit(editRequest)

    if err != nil {
        abs.logger.Error("Error updating property from airbnb-fetcher", err.Error())
        return
    }

    abs.logger.Infof("Property %v successfully updated", property.Id)
}

func safeSlice(s string, top int) string {
    if len(s) > top {
        return s[:top]
    }
    return s
}
