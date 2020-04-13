package papertrail

import (
	"log"
	"strings"
)

// Header necessary for interact with papertrail
const papertrailTokenName = "X-Papertrail-Token"
// Token necessary for interact with papertrail is obtained
// from environment variable with name PAPERTRAIL_API_TOKEN

// papertrailApiBaseUrl represents the base URL for all API operations in papertrail
const papertrailApiBaseUrl = "https://papertrailapp.com/api/v1/"

// App contains the necessary information to interact with papertrail
type App struct{}

// PapertrailNecessaryActions interacts with papertrails' API to do the necessary actions
// in function of the values provided for the options
func (a *App) PapertrailActions(options *Options) ([]Item, *string, error) {
	printActionsToDoMessage(options.Action, options.GroupName, options.SystemWildcard, options.Search, options.Query,
		options.DeleteAllSearches, options.StartDate, options.EndDate, options.Path)
	startDateUnix, endDateUnix, err := checkNecessaryPapertrailConditions(options.Action, options.SystemType, options.IpAddress, options.DestinationId,
		options.DestinationPort, options.StartDate, options.EndDate)
	if err != nil {
		return nil, nil, err
	}
	actionName := getNameOfAction(options.Action)
	createdOrDeletedItems, action, err := getItems(options.GroupName, options.SystemWildcard, options.DestinationPort,
		options.DestinationId, options.IpAddress, options.SystemType,
		options.Search, options.Query, actionName, options.DeleteAllSearches,
		startDateUnix, endDateUnix, options.Path)
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
	deleteAll bool, startDate int64, endDate int64, path string) (*[]Item, *string, error) {
	var papertrailCreatedOrRemovedItems []Item
	var err error
	papertrailCreatedOrRemovedItems, err = addSystemElements(systemType, systemWildcard,
		destinationPort, destinationId, ipAddress, actionName)
	if err != nil {
		return nil, nil, err
	}
	groupAndSearchItems, err := addGroupsAndSearches(groupName, systemWildcard, actionName,
		searchName, searchQuery, deleteAll, startDate, endDate, path)
	if err != nil {
		return &papertrailCreatedOrRemovedItems, &actionName, err
	}
	papertrailCreatedOrRemovedItems = addItemsToCreatedOrDeletedItems(groupAndSearchItems, papertrailCreatedOrRemovedItems)
	return &papertrailCreatedOrRemovedItems, &actionName, nil
}

// addGroupsAndSearches collects the information of items such as
// groups and papertrail searches created or deleted during execution
func addGroupsAndSearches(groupName string, systemWildcard string, actionName string, searchName string,
	searchQuery string, deleteAll bool, startDate int64, endDate int64, path string) ([]Item, error) {
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
		searchItem, err := doPapertrailSearchNecessaryActions(searchName, searchQuery, groupItem.ID, actionName)
		if err != nil {
			return nil, err
		}
		if ActionIsObtain(actionName) {
			eventSearchItem, err := doPapertrailEventsSearch(groupName, groupItem.ID, searchName,
				searchQuery, startDate, endDate, path)
			if err != nil {
				return nil, err
			}
			papertrailCreatedItems = addItemToCreatedOrDeletedItems(*eventSearchItem, papertrailCreatedItems)
		}
		papertrailCreatedItems = addItemToCreatedOrDeletedItems(*searchItem, papertrailCreatedItems)
	}
	return papertrailCreatedItems, nil
}

// addGroupAndSearchesDeleted collects the information of items such as
// groups and papertrail searches deleted during execution
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
		var err error
		papertrailCreatedItems, err = addSearchesAndGroupsDeleted(groupName, actionName, systemWildcard, searchName, searchQuery)
		if err != nil {
			return nil, err
		}
	}
	return papertrailCreatedItems, nil
}

// addSearchesAndGroupsDeleted collects the information of items such as
// searches and groups deleted during execution
func addSearchesAndGroupsDeleted(groupName string, actionName string,
	systemWildcard string, searchName string, searchQuery string) ([]Item, error) {
	var papertrailDeletedItems []Item
	groupItem, err := doPapertrailGroupNecessaryActions(groupName, "obtain", systemWildcard)
	if err != nil {
		return nil, err
	}
	searchItem, err := doPapertrailSearchNecessaryActions(searchName, searchQuery, groupItem.ID, actionName)
	if err != nil {
		return nil, err
	}
	papertrailDeletedItems = addItemToCreatedOrDeletedItems(*searchItem, papertrailDeletedItems)
	groupDeletedItem, err := doPapertrailGroupNecessaryActions(groupName, actionName, systemWildcard)
	if err != nil {
		return nil, err
	}
	papertrailDeletedItems = addItemToCreatedOrDeletedItems(*groupDeletedItem, papertrailDeletedItems)
	return papertrailDeletedItems, nil
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
				if systemItem != nil {
					papertrailCreatedItems = addItemToCreatedOrDeletedItems(*systemItem, papertrailCreatedItems)
				}
			} else if systemTypeIsIpAddress(systemType) {
				systemItem, err := getSystemInPapertrailBasedInAddressIp(ipAddress, actionName)
				if err != nil {
					return nil, err
				}
				if systemItem != nil {
					papertrailCreatedItems = addItemToCreatedOrDeletedItems(*systemItem, papertrailCreatedItems)
				}
			}
		}
	}
	return papertrailCreatedItems, nil
}

func printActionsToDoMessage(actionName string, groupName string, systemWildcard string, search string, query string,
	deleteAllSearches bool, startDate string, endDate string, path string) {
	if deleteAllSearches {
		log.Printf("Checking conditions for do action '%s' in papertrail params: " +
			"[group-name %s] [system-wildcard %s] [delete-all-searches %t] [--start-date %s] [--end-date %s] [--path %s]\n",
			actionName, groupName, systemWildcard, deleteAllSearches, startDate, endDate, path)
	} else {
		log.Printf("Checking conditions for do action '%s' in papertrail params: " +
			"[group-name %s] [system-wildcard %s] [search %s] [query %s] " +
			"[delete-all-searches %t] [--start-date %s] [--end-date %s] [--path %s]\n",
			actionName, groupName, systemWildcard, search, query,
			deleteAllSearches, startDate, endDate, path)
	}
}