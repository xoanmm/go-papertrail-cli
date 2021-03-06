package papertrail

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

// papertrailApiDestinationsEndpoint represents the endpoint for interact with
// groups in papertrail API
const papertrailApiDestinationsEndpoint = papertrailApiBaseUrl + "destinations.json"

// checkIfDestinationExistById checks if a system exists on papertrail with the provided identifier
func checkIfDestinationExistById(destinationId int) (*Destination, error) {
	destinationIdUrl := strings.SplitAfter(papertrailApiDestinationsEndpoint, "destinations")[0] +
		"/" + strconv.Itoa(destinationId) + strings.SplitAfter(papertrailApiDestinationsEndpoint, "destinations")[1]
	getDestination, err := apiOperation("GET", destinationIdUrl, nil)
	if err != nil {
		return nil, err
	}
	if getDestination.StatusCode == 200 {
		var destination *Destination
		json.Unmarshal([]byte(getDestination.Body), &destination)
		log.Printf("Destination with id %d exists\n", destination.ID)
		return destination, nil
	}
	err = convertStatusCodeToError(getDestination.StatusCode, "Destination", "Obtaining")
	return nil, err
}
