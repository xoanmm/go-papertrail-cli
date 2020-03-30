package papertrail

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// papertrailApiSearchesEndpoint represents the endpoint for interact with
// searchs in papertrail API
const papertrailApiSearchesEndpoint = papertrailApiBaseUrl + "searches.json"

// getSearchInPapertrailGroup obtains a papertrail search, creating it in case it does not exist previously
func getSearchInPapertrailGroup(searchName string, searchQuery string, groupId int) (*Item, error) {
	searchExists, searchObject, err := checkSearchExists(searchName, searchQuery, groupId)
	if err != nil {
		return nil, err
	}
	searchItem := Item{
		ItemType: "Search",
		Created:  false,
	}
	if *searchExists {
		fmt.Printf("Search with name %s already exists with id %d\n", searchObject.Name, searchObject.ID)
		searchItem.ID = searchObject.ID
		searchItem.ItemName = searchObject.Name
	} else {
		fmt.Printf("Search with name %s doesn't exist yet\n", searchName)
		papertrailSearchCreated, err := createPapertrailSearch(searchName, searchQuery, groupId)
		if err != nil {
			return nil, err
		}
		searchItem.ID = papertrailSearchCreated.ID
		searchItem.ItemName = papertrailSearchCreated.Name
		searchItem.Created = true
	}
	return &searchItem, err
}

// checkSearchExists checks if a search exists in papertrail specific group, returning the information
// of this one in case it exists
func checkSearchExists(searchName string, searchQuery string, groupId int) (*bool, *SearchObject, error) {
	alreadyExists := false
	var search SearchObject
	getAllSearchesResp, err  := papertrailApiOperation("GET", papertrailApiSearchesEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	if getAllSearchesResp.StatusCode == 200 {
		var searches []SearchObject
		json.Unmarshal([]byte(getAllSearchesResp.Body), &searches)
		for _, item := range searches {
			if item.Name == searchName && item.Query == searchQuery && item.Group.ID == groupId {
				alreadyExists = true
				search.ID = item.ID
				search.Name = item.Name
				search.Query = item.Query
				search.Links = item.Links
				search.Group = item.Group
				break
			}
		}
	}
	return &alreadyExists, &search, nil
}

// createPapertrailGroup creates a papertrail search using the parameter information
// provided as the group information to be created in a specific group
func createPapertrailSearch(searchName string, searchQuery string, groupId int) (*SearchObject, error){
	var search SearchObject
	papertrailSearchToCreate := SearchToCreateObject{SearchToCreate: SearchToCreate{
		Name:           searchName,
		Query: 			searchQuery,
		GroupID:		groupId,
	}}
	b, err := json.Marshal(papertrailSearchToCreate)
	if err != nil {
		return nil, err
	}
	createSearchResp, err  := papertrailApiOperation("POST", papertrailApiSearchesEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	if createSearchResp.StatusCode == 200 {
		json.Unmarshal([]byte(createSearchResp.Body), &search)
		fmt.Printf("Search with name %s and id %d was successfully created\n", search.Name, search.ID)
		return &search, nil
	}
	fmt.Printf("Problems creating search with name %s in group with id %d\n", searchName, groupId)
	return &search, nil
}