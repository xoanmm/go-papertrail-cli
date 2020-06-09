package papertrail

import (
	"strings"
)

// Header necessary for interact with papertrail
const papertrailTokenName = "X-Papertrail-Token"

// papertrailApiBaseUrl represents the base URL for all API operations in papertrail
const papertrailApiBaseUrl = "https://papertrailapp.com/api/v1/"

// App contains the necessary information to interact with papertrail
type App struct{}

// PapertrailActions interacts with papertrails' API to do the necessary actions
// in function of the values provided for the options
func (a *App) PapertrailActions(options *Options) ([]Item, *string, error) {
	var err error
	printActionsToDoMessage(*options)
	startDateUnix, endDateUnix, err := cnvStDateEndDateToUnixTime(options.StartDate, options.EndDate)
	if err != nil {
		return nil, nil, err
	}
	err = checkNecessaryConditions(options.Action, options.SystemType, options.IpAddress, options.DestinationId,
		options.DestinationPort, startDateUnix, endDateUnix)
	if err != nil {
		return nil, nil, err
	}
	actionName := getNameOfAction(options.Action)
	createdOrDeletedItems, action, err := getItems(*options, actionName, startDateUnix, endDateUnix)
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
func getItems(options Options, actionName string, startDate int64, endDate int64) (*[]Item, *string, error) {
	var papertrailCreatedOrRemovedItems []Item
	var err error
	if !options.DeleteOnlySearches {
		papertrailCreatedOrRemovedItems, err = addSystemElements(options.SystemType, options.SystemWildcard,
			options.DestinationPort, options.DestinationId, options.IpAddress, actionName, options.DeleteAllSystems)
		if err != nil {
			return nil, nil, err
		}
	}
	if !options.DeleteOnlySystems {
		groupAndSearchItems, err := addGroupsAndSearches(options.GroupName, options.SystemWildcard, actionName,
			options.Search, options.Query, options.DeleteAllSearches, options.DeleteAllSystems, startDate, endDate, options.Path)
		if err != nil {
			return &papertrailCreatedOrRemovedItems, &actionName, err
		}
		papertrailCreatedOrRemovedItems = addItemsToCreatedOrDeletedItems(groupAndSearchItems, papertrailCreatedOrRemovedItems)
	}
	return &papertrailCreatedOrRemovedItems, &actionName, nil
}

// addGroupsAndSearches collects the information of items such as
// groups and papertrail searches created or deleted during execution
func addGroupsAndSearches(groupName string, systemWildcard string, actionName string, searchName string,
	searchQuery string, deleteAll bool, deleteAllSystems bool, startDate int64, endDate int64, path string) ([]Item, error) {
	var papertrailCreatedItems []Item
	if ActionIsDelete(actionName) {
		var err error
		papertrailCreatedItems, err = addGroupAndSearchesDeleted(deleteAll, groupName, actionName,
			systemWildcard, searchName, searchQuery, deleteAllSystems)
		if err != nil {
			return nil, err
		}
	} else {
		groupItem, err := doPapertrailGroupNecessaryActions(groupName, actionName, systemWildcard, deleteAllSystems)
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
func addGroupAndSearchesDeleted(deleteAllSearchs bool, groupName string, actionName string,
	systemWildcard string, searchName string, searchQuery string, deleteAllSystems bool) ([]Item, error) {
	var papertrailDeletedItems []Item
	if deleteAllSearchs {
		groupItem, err := doPapertrailGroupNecessaryActions(groupName, actionName, systemWildcard, deleteAllSystems)
		if err != nil {
			return nil, err
		}
		if groupItem != nil {
			papertrailDeletedItems = addItemToCreatedOrDeletedItems(*groupItem, papertrailDeletedItems)
		}
	} else {
		groupItem, err := doPapertrailGroupNecessaryActions(groupName, "obtain", systemWildcard, deleteAllSystems)
		if err != nil {
			return nil, err
		}
		if groupItem != nil {
			searchItem, err := doPapertrailSearchNecessaryActions(searchName, searchQuery, groupItem.ID, actionName)
			if err != nil {
				return nil, err
			}
			papertrailDeletedItems = addItemToCreatedOrDeletedItems(*searchItem, papertrailDeletedItems)
		}
	}
	return papertrailDeletedItems, nil
}

// addSystemElements collects specific system/s details and adds
// them to the list of created/deleted items if they have been created or deleted
func addSystemElements(systemType string, systemWildcard string, destinationPort int,
	destinationId int, ipAddress string, actionName string, deleteAllSystems bool) ([]Item, error) {
	var papertrailCreatedItems []Item
	if systemWildcard != "*" && checkConditionsForDeleteAllSystems(actionName, deleteAllSystems) {
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

// printActionsToDoMessage prints a message at the beginning of the execution
// with the value of the parameters needed to perform the necessary action
func printActionsToDoMessage(options Options) {
	if ActionIsDelete(options.Action) {
		printMessageActionDelete(options)
	} else if ActionIsObtain(options.Action) {
		printMessageActionObtain(options)
	} else {
		printMessageActionCreate(options)
	}
}
