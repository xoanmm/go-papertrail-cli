package papertrail

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// papertrailApiGroupsEndpoint represents the endpoint for interact with
// groups in papertrail API
const papertrailApiGroupsEndpoint = papertrailApiBaseUrl + "groups.json"

// getGroupInPapertrail obtains a papertrail group, creating it in case it does not exist previously
func getGroupInPapertrail(groupName string, systemWildcard string) (*Item, error) {
	groupExists, groupObject, err := checkGroupExists(groupName)
	if err != nil {
		return nil, err
	}
	var groupItem *Item
	if *groupExists {
		fmt.Printf("Group with name %s already exists with id %d\n", groupName, groupObject.ID)
		groupItem = NewItem(groupObject.ID, "Group", groupObject.Name, false)
	} else {
		fmt.Printf("Group with name %s doesn't exist yet\n", groupName)
		papertrailGroupCreated, err := createPapertrailGroup(groupName, systemWildcard)
		if err != nil {
			return nil, err
		}
		groupItem = NewItem(papertrailGroupCreated.ID, "Group", papertrailGroupCreated.Name, true)
	}
	return groupItem, err
}

// checkGroupExists checks if a group exists in papertrail, returning the information of this one in case it exists
func checkGroupExists(groupName string) (*bool, *GroupObject, error) {
	alreadyExists := false
	var group *GroupObject
	getAllGroupResp, err  := papertrailApiOperation("GET", papertrailApiGroupsEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	if getAllGroupResp.StatusCode == 200 {
		var groups []GroupObject
		json.Unmarshal([]byte(getAllGroupResp.Body), &groups)
		for _, item := range groups {
			if item.Name == groupName {
				alreadyExists = true
				group = NewGroupObject(item.ID, item.Name, item.SystemWildcard, item.Links, item.Systems)
				break
			}
		}
	}
	return &alreadyExists, group, nil
}

// createPapertrailGroup creates a papertrail group using the parameter information
// provided as the group information to be created
func createPapertrailGroup(groupName string, systemWildcard string) (*GroupObject, error){
	papertrailGroupToCreate := GroupCreationObject{Group: GroupCreateObject{
		Name:           groupName,
		SystemWildcard: systemWildcard,
	}}
	b, err := json.Marshal(papertrailGroupToCreate)
	if err != nil {
		return nil, err
	}
	createGroupResp, err  := papertrailApiOperation("POST", papertrailApiGroupsEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	if createGroupResp.StatusCode == 200 {
		var group GroupObject
		json.Unmarshal([]byte(createGroupResp.Body), &group)
		fmt.Printf("Group with name %s and id %d was successfully created\n", group.Name, group.ID)
		return &group, nil
	}
	fmt.Printf("Problems creating group with name %s\n", groupName)
	err = errors.New("Error: Response status code " + strconv.Itoa(createGroupResp.StatusCode))
	return nil, err
}