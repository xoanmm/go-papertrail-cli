package papertrail

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

// ApiOperation is a generic function to interact with the papertrail API, in which
// a series of headers necessary for the interaction with this API are established.
// Through the parameters it is possible to indicate the type of operation, the body to be sent
// and the specific URL of the API
func ApiOperation(method string, url string, bodyToSend io.Reader) (*ApiResponse, error) {
	req, err := http.NewRequest(method, url, bodyToSend)
	req.Header.Add(papertrailTokenName, papertrailToken)
	req.Header.Add("Content-Type", "application/json")
	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return  nil, err
	}
	return &ApiResponse{
		Body:           body,
		StatusCode: 	resp.StatusCode,
		err:			err,
	}, nil
}

// checkNecessaryPapertrailConditions checks if the conditions to provide a token to interact
// with papertrail are met, as well as that a valid action is provided (c/create or o/obtain)
func checkNecessaryPapertrailConditions(action string, systemType string, ipAddress string, destinationId int,
	destinationPort int) {
	if len(papertrailToken) == 0 {
		log.Fatalf("Error getting value of PAPERTRAIL_API_TOKEN, " +
			"it's necessary to define this variable with your papertrail's API token")
	}
	checkValidActionsConditions(action)
	checkValidSystemTypeConditions(systemType, ipAddress, destinationId, destinationPort, action)
}

func checkValidActionsConditions(action string) {
	validActions := []string{"c", "create", "o", "obtain", "d", "delete"}
	_, found := Find(validActions, action)
	if !found {
		log.Fatalf("Not valid option provided for action to perform, the only valid values are: \n" +
			"\t'c' or 'create': create new groups or search\n" +
			"\t'o'or 'obtain': obtain logs in base of parameters provided\n")
	}
}

func systemTypeIsHostname(systemType string) bool {
	if systemType == "h" || systemType == "hostname" {
		return true
	}
	return false
}

func systemTypeIsIpAddress(systemType string) bool {
	if systemType == "i" || systemType == "ip-address" {
		return true
	}
	return false
}

func checkValidSystemTypeConditions(systemType  string, ipAddress string, destinationId int,
	destinationPort int, actionType string) {
	validSystemTypes := []string{"h", "hostname", "i", "ip-address"}
	_, found := Find(validSystemTypes, systemType)
	if !found {
		log.Fatalf("Not valid option provided for system, the only valid values are: \n" +
			"\t'h' or 'hostname': system based in hostname\n" +
			"\t'i'or 'ip-address': system based in ip-address\n")
	} else if !ActionIsDelete(actionType) {
		if systemTypeIsHostname(systemType) {
			if destinationId != 0 && destinationPort != 0 {
				log.Fatalf("If the system is a hostname-type system, only destination " +
					"id or destination port can be specified\n")
			} else if !((destinationId != 0 && destinationPort == 0) ||
				(destinationId == 0 && destinationPort != 0)) {
				log.Fatalf("It's necessary provide a value distinct from default (0) to " +
					"destination id or destination port")
			}
		} else if systemTypeIsIpAddress(systemType) {
			re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
			if !re.MatchString(ipAddress) {
				log.Fatalf("The IP Address provided, %s, it's not a valid IP Address\n", ipAddress)
			}
		}
	}
}

// getNameOfAction returns the name of the action to be performed according to the value obtained
// for the action parameter
func getNameOfAction(actionOptionName string) string {
	if ActionIsCreate(actionOptionName) {
		return "create"
	} else if ActionIsDelete(actionOptionName) {
		return "delete"
	}
	return "obtain"
}

func ActionIsCreate(actionOptionName string) bool {
	if actionOptionName == "c" || actionOptionName == "create" {
		return true
	}
	return false
}

func ActionIsDelete(actionOptionName string) bool {
	if actionOptionName == "d" || actionOptionName == "delete" {
		return true
	}
	return false
}

func ActionIsObtain(actionOptionName string) bool {
	if actionOptionName == "o" || actionOptionName == "obtain" {
		return true
	}
	return false
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func checkElementsCreatedOrRemoved(papertrailToAddItems []Item) bool {
	createdOrRemoved := false
	for _, item := range papertrailToAddItems {
		if item.Deleted || item.Created {
			createdOrRemoved = true
			break
		}
	}
	return createdOrRemoved
}