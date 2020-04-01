package papertrail

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
)

// papertrailApiDestinationsEndpoint represents the endpoint for interact with
// groups in papertrail API
const papertrailApiDestinationsEndpoint = papertrailApiBaseUrl + "destinations.json"

// checkIfDestinationExistById checks if a system exists on papertrail with the provided identifier
func checkIfDestinationExistById(destinationId int) (*bool, *Destination, error) {
	exists := false
	destinationIdUrl := strings.SplitAfter(papertrailApiDestinationsEndpoint, "destinations")[0] +
		"/" + strconv.Itoa(destinationId) + strings.SplitAfter(papertrailApiDestinationsEndpoint, "destinations")[1]
	getDestination, err  := apiOperation("GET", destinationIdUrl, nil)
	if err != nil {
		return nil, nil, err
	}
	if getDestination.StatusCode == 200 {
		exists = true
		var destination *Destination
		json.Unmarshal([]byte(getDestination.Body), &destination)
		log.Printf("Destination with id %d exists\n", destination.ID)
		return &exists, destination, nil
	}
	err = errors.New("Error: Response status code " + strconv.Itoa(getDestination.StatusCode))
	return nil, nil, err
}
