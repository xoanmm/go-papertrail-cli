package papertrail

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

// checkValidActionsConditions checks if a valid value is being used for the action parameter
func CheckValidActionsConditions(action string) error {
	validActions := []string{"c", "create", "o", "obtain", "d", "delete"}
	_, found := find(validActions, action)
	if !found {
		return errors.New("Not valid option provided for action to perform, the only valid values are: \n" +
			"\t'c' or 'create': create new system/s, group and/or search\n" +
			"\t'd' or 'delete': create new system/s, group and/or search\n" +
			"\t'o'or 'obtain': obtain logs in base of parameters provided\n")
	}
	return nil
}

// checkNecessaryPapertrailConditions checks if the conditions to provide a token to interact
// with papertrail are met, as well as that a valid action is provided (c/create, d/delete or o/obtain)
func checkNecessaryPapertrailConditions(action string, systemType string, ipAddress string,
	destinationId int, destinationPort int) error {
	papertrailToken := os.Getenv("PAPERTRAIL_API_TOKEN")
	if len(papertrailToken) == 0 {
		return errors.New("Error getting value of PAPERTRAIL_API_TOKEN, " +
			"it's necessary to define this variable with your papertrail's API token ")
	}
	err := CheckValidActionsConditions(action)
	if err != nil {
		return err
	}
	err = checkValidSystemTypeConditions(systemType, ipAddress, destinationId, destinationPort, action)
	if err != nil {
		return err
	}
	return nil
}

// apiOperation is a generic function to interact with the papertrail API, in which
// a series of headers necessary for the interaction with this API are established.
// Through the parameters it is possible to indicate the type of operation, the body to be sent
// and the specific URL of the API
func apiOperation(method string, url string, bodyToSend io.Reader) (*ApiResponse, error) {
	req, err := http.NewRequest(method, url, bodyToSend)
	req.Header.Add(papertrailTokenName, os.Getenv("PAPERTRAIL_API_TOKEN"))
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

// systemTypeIsHostname checks if the type of system entered
// is hostname for the system-type parameter
func systemTypeIsHostname(systemType string) bool {
	if systemType == "h" || systemType == "hostname" {
		return true
	}
	return false
}

// systemTypeIsHostname checks if the type of system entered
// is ip-address for the system-type parameter
func systemTypeIsIpAddress(systemType string) bool {
	if systemType == "i" || systemType == "ip-address" {
		return true
	}
	return false
}

// checkValidSystemTypeConditions checks that whether ip-address or hostname and system
// type have been entered, the configuration values on them are valid
func checkValidSystemTypeConditions(systemType  string, ipAddress string, destinationId int,
	destinationPort int, actionType string) error {
	validSystemTypes := []string{"h", "hostname", "i", "ip-address"}
	_, found := find(validSystemTypes, systemType)
	if !found {
		return errors.New("Not valid option provided for system, the only valid values are: \n" +
			"\t'h' or 'hostname': system based in hostname\n" +
			"\t'i'or 'ip-address': system based in ip-address\n")
	} else if !actionIsDelete(actionType) {
		if systemTypeIsHostname(systemType) {
			if destinationId != 0 && destinationPort != 0 {
				return errors.New("If the system is a hostname-type system, only destination " +
					"id or destination port can be specified\n")
			} else if !((destinationId != 0 && destinationPort == 0) ||
				(destinationId == 0 && destinationPort != 0)) {
				return errors.New("It's necessary provide a value distinct from default (0) to " +
					"destination id or destination port ")
			}
		} else if systemTypeIsIpAddress(systemType) {
			re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
			if !re.MatchString(ipAddress) {
				return errors.New("The IP Address provided, " + ipAddress + " it's not a valid IP Address ")
			}
		}
	}
	return nil
}

// actionIsCreate checks if the value entered for the action parameter is to create
func actionIsCreate(actionOptionName string) bool {
	if actionOptionName == "c" || actionOptionName == "create" {
		return true
	}
	return false
}

// actionIsDelete checks if the value entered for the action parameter is to delete
func actionIsDelete(actionOptionName string) bool {
	if actionOptionName == "d" || actionOptionName == "delete" {
		return true
	}
	return false
}

// actionIsObtain checks if the value entered for the action parameter is to obtain
func actionIsObtain(actionOptionName string) bool {
	if actionOptionName == "o" || actionOptionName == "obtain" {
		return true
	}
	return false
}

// getNameOfAction returns the name of the action to be performed according to the value obtained
// for the action parameter
func getNameOfAction(actionOptionName string) string {
	if actionIsCreate(actionOptionName) {
		return "create"
	} else if actionIsDelete(actionOptionName) {
		return "delete"
	}
	return "obtain"
}

// find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// getOnlyElementsCreatedOrRemoved concatenates to the list of elements created or
// removed from papertrail those of the provided list, adding only those that
// fulfill the condition of created or removed to the first list
func getOnlyElementsCreatedOrRemoved(papertrailToAddItems []Item, createdOrRemovedItems []Item) []Item {
	for _, item := range papertrailToAddItems {
		if item.Deleted || item.Created {
			createdOrRemovedItems = append(createdOrRemovedItems, item)
		}
	}
	return createdOrRemovedItems
}

// addItemsToCreatedOrDeletedItems checks whether the list of papertrail item has been created
// or deleted during execution or not, if it has been created/deleted it is added to the list of created items
func addItemsToCreatedOrDeletedItems(papertrailToAddItems []Item, papertrailItemsCreatedOrDeleted []Item) []Item {
	return getOnlyElementsCreatedOrRemoved(papertrailToAddItems, papertrailItemsCreatedOrDeleted)
}

// addItemsToCreatedOrDeletedItems checks whether the papertrail item has been created or deleted during
// execution or not, if it has been created/deleted it is added to the list of created items
func addItemToCreatedOrDeletedItems(papertrailToAddItem Item, papertrailItemsCreatedOrDeleted []Item) []Item {
	var newItems []Item
	if papertrailToAddItem.Deleted || papertrailToAddItem.Created {
		newItems = append(papertrailItemsCreatedOrDeleted, papertrailToAddItem)
	} else {
		return papertrailItemsCreatedOrDeleted
	}
	return newItems
}