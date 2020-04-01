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
		options.DestinationId, options.IpAddress, options.SystemType,
		options.Search, options.Query, actionName, options.DeleteAllSearches)
	if err != nil {
		if createdOrDeletedItems != nil {
			return *createdOrDeletedItems, action, err
		}
		return nil, nil, err
	}
	return *createdOrDeletedItems, action, err
}

// getItems collects specific group and/or search details and adds
// them to the list of created items if they have been created
func getItems(groupName string, systemWildcard string, destinationPort int, destinationId int,
	ipAddress string, systemType string, searchName string, searchQuery string, actionName string,
	deleteAll bool) (*[]Item, *string, error) {
	var papertrailCreatedOrRemovedItems []Item
	var err error
	papertrailCreatedOrRemovedItems, err = addSystemElements(systemType, systemWildcard,
		destinationPort, destinationId, ipAddress, actionName)
	if err != nil {
		return nil, nil, err
	}
	groupAndSearchItems, err := addGroupsAndSearches(groupName, systemWildcard, actionName, searchName, searchQuery, deleteAll)
	if err != nil {
		return &papertrailCreatedOrRemovedItems, &actionName, err
	}
	papertrailCreatedOrRemovedItems = addItemsToCreatedOrDeletedItems(groupAndSearchItems, papertrailCreatedOrRemovedItems)
	return &papertrailCreatedOrRemovedItems, &actionName, nil
}

func addGroupsAndSearches(groupName string, systemWildcard string, actionName string, searchName string,
	searchQuery string, deleteAll bool) ([]Item, error) {
	var papertrailCreatedItems []Item
	if ActionIsDelete(actionName) {
		var err error
		papertrailCreatedItems, err = addGroupAndSearchesDeleted(deleteAll, groupName, actionName,
		systemWildcard, searchName, searchQuery)
		if err != nil {
			return nil, err
		}
	} else {
		groupItem, err := doPapertrailGroupNecessaryActions(groupName, actionName, systemWildcard)
		if err != nil {
			return nil, err
		}
		papertrailCreatedItems = addItemToCreatedOrDeletedItems(*groupItem, papertrailCreatedItems)
		searchItem, err := doPapertrailSearchesNecessaryActions(searchName, searchQuery, groupItem.ID, actionName)
		if err != nil {
			return nil, err
		}
		papertrailCreatedItems = addItemToCreatedOrDeletedItems(*searchItem, papertrailCreatedItems)
	}
	return papertrailCreatedItems, nil
}

func addGroupAndSearchesDeleted(deleteAll bool, groupName string, actionName string,
	systemWildcard string, searchName string, searchQuery string) ([]Item, error) {
	var papertrailCreatedItems []Item
	if deleteAll {
		groupItem, err := doPapertrailGroupNecessaryActions(groupName, actionName, systemWildcard)
		if err != nil {
			return nil, err
		}
		papertrailCreatedItems = addItemToCreatedOrDeletedItems(*groupItem, papertrailCreatedItems)
	} else {
		groupItem, err := doPapertrailGroupNecessaryActions(groupName, "obtain", systemWildcard)
		if err != nil {
			return nil, err
		}
		searchItem, err := doPapertrailSearchesNecessaryActions(searchName, searchQuery, groupItem.ID, actionName)
		if err != nil {
			return nil, err
		}
		papertrailCreatedItems = addItemToCreatedOrDeletedItems(*searchItem, papertrailCreatedItems)
	}
	return papertrailCreatedItems, nil
}

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

// addItemsToCreatedOrDeletedItems checks whether the papertrail item has been created or deleted during
// execution or not, if it has been created/deleted it is added to the list of created items
func addItemsToCreatedOrDeletedItems(papertrailToAddItems []Item, papertrailItemsCreatedOrDeleted []Item) []Item {
	if checkElementsCreatedOrRemoved(papertrailToAddItems) {
		return append(papertrailItemsCreatedOrDeleted, papertrailToAddItems...)
	}
	return papertrailItemsCreatedOrDeleted

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