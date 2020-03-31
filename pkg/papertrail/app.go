package papertrail

import (
	"log"
	"os"
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
func (a *App) PapertrailNecessaryActions(options *Options) ([]Item, *string, error) {
	checkNecessaryPapertrailConditions(options.Action, options.SystemType, options.IpAddress, options.DestinationId,
		options.DestinationPort)
	actionName := getNameOfAction(options.Action)
	log.Printf("Checking conditions for %s in papertrail params: " +
		"[group-name %s] [system-wildcard %s] [search %s] [query %s]\n",
		actionName, options.GroupName, options.SystemWildcard, options.Search, options.Query)
	createdOrDeletedItems, action, err := getItems(options.GroupName, options.SystemWildcard, options.DestinationPort,
		options.DestinationId, options.IpAddress, options.SystemType, options.Search, options.Query, actionName)
	if err != nil {
		return nil, nil, err
	}
	return *createdOrDeletedItems, action, err
}

// getItems collects specific group and/or search details and adds
// them to the list of created items if they have been created
func getItems(groupName string, systemWildcard string, destinationPort int, destinationId int,
	ipAddress string, systemType string, searchName string, searchQuery string, actionName string) (*[]Item, *string, error) {
	var papertrailCreatedItems []Item
	var err error
	papertrailCreatedItems, err = addSystemElements(systemType, systemWildcard,
		destinationPort, destinationId, ipAddress, actionName)
	if err != nil {
		return nil, nil, err
	}
	groupItem, err := doPapertrailGroupNecessaryActions(groupName, actionName, systemWildcard)
	if err != nil {
		return nil, nil, err
	}
	papertrailCreatedItems = addItemToCreatedOrDeletedItems(*groupItem, papertrailCreatedItems)
	if !ActionIsDelete(actionName) {
		searchItem, err := getSearchInPapertrailGroup(searchName, searchQuery, groupItem.ID)
		if err != nil {
			return nil, nil, err
		}
		papertrailCreatedItems = addItemToCreatedOrDeletedItems(*searchItem, papertrailCreatedItems)
	}
	return &papertrailCreatedItems, &actionName, nil
}

// TODO: function addGroupsAndSearches need distinction between create and delete for recover necessary information and delete respecting relationships between elements

//func addGroupsAndSearches(groupName string, systemWildcard string, actionName string, searchName string,
//	searchQuery string) ([]Item, error) {
//	var papertrailCreatedItems []Item
//	if ActionIsDelete(actionName) {
//		searchItem, err := getSearchInPapertrailGroup(searchName, searchQuery, groupItem.ID)
//		if err != nil {
//			return nil, err
//		}
//		papertrailCreatedItems = addItemToCreatedOrDeletedItems(*searchItem, papertrailCreatedItems)
//		groupItem, err := doPapertrailGroupNecessaryActions(groupName, systemWildcard, actionName)
//		if err != nil {
//			return nil, err
//		}
//		papertrailCreatedItems = addItemToCreatedOrDeletedItems(*groupItem, papertrailCreatedItems)
//	} else if ActionIsCreate(actionName) {
//		groupItem, err := doPapertrailGroupNecessaryActions(groupName, systemWildcard, actionName)
//		if err != nil {
//			return nil, err
//		}
//		papertrailCreatedItems = addItemToCreatedOrDeletedItems(*groupItem, papertrailCreatedItems)
//		if !ActionIsDelete(actionName) {
//			searchItem, err := getSearchInPapertrailGroup(searchName, searchQuery, groupItem.ID)
//			if err != nil {
//				return nil, err
//			}
//			papertrailCreatedItems = addItemToCreatedOrDeletedItems(*searchItem, papertrailCreatedItems)
//		}
//	}
//	return papertrailCreatedItems, nil
//}

// addSystemElements collects specific system/s details and adds
// them to the list of created/deleted items if they have been created or deleted
func addSystemElements(systemType string, systemWildcard string, destinationPort int,
	destinationId int, ipAddress string, actionName string) ([]Item, error) {
	var papertrailCreatedItems []Item
	if systemWildcard != "*" {
		systems := strings.Split(systemWildcard, ", ")
		for _, item := range systems {
			if systemTypeIsHostname(systemType) {
				systemItem, err := getSystemInPapertrailBasedInHostname(item, destinationPort, destinationId, actionName)
				if err != nil {
					return nil, err
				}
				papertrailCreatedItems = addItemToCreatedOrDeletedItems(*systemItem, papertrailCreatedItems)
			} else if systemTypeIsIpAddress(systemType) {
				systemItem, err := getSystemInPapertrailBasedInAddressIp(ipAddress)
				if err != nil {
					return nil, err
				}
				papertrailCreatedItems = addItemToCreatedOrDeletedItems(*systemItem, papertrailCreatedItems)
			}
		}
	}
	return papertrailCreatedItems, nil
}

// addItemToCreatedOrDeletedItems checks whether the papertrail item has been created or deleted during
// execution or not, if it has been created/deleted it is added to the list of created items
func addItemToCreatedOrDeletedItems(papertrailItem Item, papertrailItemsCreatedOrDeleted []Item) []Item {
	var newItems []Item
	if papertrailItem.Created || papertrailItem.Deleted{
		newItems = append(papertrailItemsCreatedOrDeleted, papertrailItem)
	} else {
		return papertrailItemsCreatedOrDeleted
	}
	return newItems
}