package papertrail

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
)

// papertrailApiSearchesEndpoint represents the endpoint for interact with
// searches in papertrail API
const papertrailApiSearchesEndpoint = papertrailApiBaseUrl + "searches.json"

// doPapertrailSearchesNecessaryActions is in charge of carrying out the indicated actions
// on the indicated papertrail search, as well as checking if it exists
func doPapertrailSearchNecessaryActions(searchName string, searchQuery string, groupId int,
	actionName string) (*Item, error) {
	var searchItem *Item
	searchObject, err := checkSearchExists(searchName, groupId)
	if err != nil {
		return nil, err
	}
	if searchObject != nil {
		log.Printf("Search with name %s exists with id %d\n", searchName, searchObject.ID)
		if ActionIsObtain(actionName) || ActionIsCreate(actionName) {
			return NewItem(searchObject.ID, "Search", searchName, false, false), nil
		} else if ActionIsDelete(actionName) {
			searchItem, err = deleteSearch(searchName, searchObject.ID)
			if err != nil {
				return nil, err
			}
		}
	} else {
		if ActionIsCreate(actionName) {
			searchItem, err = createSearch(searchName, searchQuery, groupId)
			if err != nil {
				return nil, err
			}
		} else {
			err := errors.New("Error: Search with name " + searchName + " doesn't exist")
			return nil, err
		}
	}
	return searchItem, nil
}

// createSearch attempts to create a search in papertrail using the parameters provided as search information
func createSearch(searchName string, searchQuery string, groupId int) (*Item, error) {
	papertrailSearchCreated, err := createPapertrailSearchOperation(searchName, searchQuery, groupId)
	if err != nil {
		return nil, err
	}
	return NewItem(papertrailSearchCreated.ID, "Search", searchName, true, false), nil
}

// deleteSearch attempts to delete a search using the parameters provided as search information
func deleteSearch(searchName string, searchId int) (*Item, error) {
	papertrailSearchDeleted, err := deletePapertrailSearchOperation(searchName, searchId)
	if err != nil {
		return nil, err
	}
	if *papertrailSearchDeleted {
		return NewItem(searchId, "Search", searchName, false, true), nil
	}
	return nil, errors.New("Error: Search with " + searchName + " doesn't exist")
}

// checkSearchExists checks if a search exists in papertrail specific group, returning the information
// of this one in case it exists
func checkSearchExists(searchName string, groupId int) (*SearchObject, error) {
	var search *SearchObject
	getAllSearchesResp, err := apiOperation("GET", papertrailApiSearchesEndpoint, nil)
	if err != nil {
		return search, err
	}
	if getAllSearchesResp.StatusCode == 200 {
		var searches []SearchObject
		json.Unmarshal([]byte(getAllSearchesResp.Body), &searches)
		for _, item := range searches {
			if item.Name == searchName && item.Group.ID == groupId {
				search = NewSearchObject(item.ID, item.Name, item.Query, item.Group, item.Links)
				break
			}
		}
	}
	return search, nil
}

// createPapertrailSearchOperation creates a papertrail search using the parameter information
// provided as the search information to be created in a specific group
func createPapertrailSearchOperation(searchName string, searchQuery string, groupId int) (*SearchObject, error) {
	var search SearchObject
	papertrailSearchToCreate := SearchToCreateObject{SearchToCreate: SearchToCreate{
		Name:    searchName,
		Query:   searchQuery,
		GroupID: groupId,
	}}
	b, err := json.Marshal(papertrailSearchToCreate)
	if err != nil {
		return nil, err
	}
	createSearchResp, err := apiOperation("POST", papertrailApiSearchesEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	if createSearchResp.StatusCode == 200 {
		json.Unmarshal([]byte(createSearchResp.Body), &search)
		log.Printf("Search with name %s and id %d was successfully created\n", search.Name, search.ID)
		return &search, nil
	}
	log.Printf("Problems creating search with name %s in group with id %d\n", searchName, groupId)
	err = convertStatusCodeToError(createSearchResp.StatusCode, "Search", "Creating")
	return nil, err
}

// deletePapertrailGroupOperation do the necessary calls in papertrail
// to delete a search using the parameter information provided as the search information to be deleted
func deletePapertrailSearchOperation(searchName string, searchId int) (*bool, error) {
	deleted := false
	searchIdUrl := strings.SplitAfter(papertrailApiSearchesEndpoint, "searches")[0] +
		"/" + strconv.Itoa(searchId) + strings.SplitAfter(papertrailApiSearchesEndpoint, "searches")[1]
	deleteSearchResp, err := apiOperation("DELETE", searchIdUrl, nil)
	if err != nil {
		return nil, err
	}
	if deleteSearchResp.StatusCode == 200 {
		deleted = true
		log.Printf("Search with name %s and id %d was successfully deleted\n", searchName, searchId)
		return &deleted, nil
	}
	log.Printf("Problems deleting group with id %d\n", searchId)
	err = convertStatusCodeToError(deleteSearchResp.StatusCode, "Search", "Deleting")
	return &deleted, err
}
