package papertrail

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Header necessary for interact with papertrail
const papertrailTokenName = "X-Papertrail-Token"
// Token necessary for interact with papertrail is obtained
// from environment variable with name PAPERTRAIL_API_TOKEN
var papertrailToken string = os.Getenv("PAPERTRAIL_API_TOKEN")
// papertrailApiBaseUrl represents the base URL for all API operations in papertrail
const papertrailApiBaseUrl = "https://papertrailapp.com/api/v1/"

// App contains the necessary information to interact with papertrail
type App struct{}

// PapertrailNecessaryActions interacts with papertrails' API to do the necessary actions
// in function of the values provided for the options
func (a *App) PapertrailNecessaryActions(options *Options) ([]Item, error) {
	checkNecessaryPapertrailConditions(options.Action, options.SystemType, options.IpAddress, options.DestinationId,
		options.DestinationPort)
	actionName := getNameOfAction(options.Action)
	fmt.Printf("Checking conditions for %s in papertrail params: " +
		"[group-name %s] [system-wildcard %s] [search %s] [query %s]\n",
		actionName, options.GroupName, options.SystemWildcard, options.Search, options.Query)
	createdItems, err := getItems(options.GroupName, options.SystemWildcard, options.DestinationPort,
		options.DestinationId, options.IpAddress, options.SystemType, options.Search, options.Query)
	if err != nil {
		return nil, err
	}
	return *createdItems, err
}

// getItems collects specific group and/or search details and adds
// them to the list of created items if they have been created
func getItems(groupName string, systemWildcard string, destinationPort int, destinationId int,
	ipAddress string, systemType string, searchName string, searchQuery string) (*[]Item, error) {
	var papertrailCreatedItems []Item
	systems := strings.Split(systemWildcard, ", ")
	for _, item := range systems {
		if systemTypeIsHostname(systemType) {
			systemItem, err := getSystemInPapertrailBasedInHostname(item, destinationPort, destinationId)
			if err != nil {
				return nil, err
			}
			addItemToCreatedItems(*systemItem, &papertrailCreatedItems)
		} else if systemTypeIsIpAddress(systemType) {
			systemItem, err := getSystemInPapertrailBasedInAddressIp(ipAddress)
			if err != nil {
				return nil, err
			}
			addItemToCreatedItems(*systemItem, &papertrailCreatedItems)
		}
	}
	groupItem, err := getGroupInPapertrail(groupName, systemWildcard)
	if err != nil {
		return nil, err
	}
	addItemToCreatedItems(*groupItem, &papertrailCreatedItems)
	searchItem, err := getSearchInPapertrailGroup(searchName, searchQuery, groupItem.ID)
	if err != nil {
		return nil, err
	}
	addItemToCreatedItems(*searchItem, &papertrailCreatedItems)
	return &papertrailCreatedItems, nil
}

// addItemToCreatedItems checks whether the papertrail item has been created during
// execution or not, if it has been created it is added to the list of created items
func addItemToCreatedItems(papertrailItem Item, papertrailItemsCreated *[]Item) *[]Item {
	if papertrailItem.Created {
		*papertrailItemsCreated = append(*papertrailItemsCreated, papertrailItem)
	}
	return papertrailItemsCreated
}

// papertrailApiOperation is generic function to interact with the papertrail API, in which
// a series of headers necessary for the interaction with this API are established.
// Through the parameters it is possible to indicate the type of operation, the body to be sent
// and the specific URL of the API
func papertrailApiOperation(method string, url string, bodyToSend io.Reader) (*ApiResponse, error) {
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
	checkValidSystemTypeConditions(systemType, ipAddress, destinationId, destinationPort)
}

func checkValidActionsConditions(action string) {
	validActions := []string{"c", "create", "o", "obtain"}
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

func checkValidSystemTypeConditions(systemType  string, ipAddress string, destinationId int, destinationPort int) {
	validSystemTypes := []string{"h", "hostname", "i", "ip-address"}
	_, found := Find(validSystemTypes, systemType)
	if !found {
		log.Fatalf("Not valid option provided for system, the only valid values are: \n" +
			"\t'h' or 'hostname': system based in hostname\n" +
			"\t'i'or 'ip-address': system based in ip-address\n")
	} else {
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
	if actionOptionName == "c" || actionOptionName == "create" {
		return "create"
	}
	return "obtain"
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

// CheckErr checks if given error is not nil and exit program with signal 1
func CheckErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}