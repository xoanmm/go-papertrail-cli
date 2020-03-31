package papertrail

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
)

// papertrailApiGroupsEndpoint represents the endpoint for interact with
// groups in papertrail API
const papertrailApiGroupsEndpoint = papertrailApiBaseUrl + "groups.json"

//
func doPapertrailGroupNecessaryActions(groupName string, actionName string, systemWildcard string) (*Item, error) {
	var groupItem *Item
	groupExists, groupObject, err := checkGroupExists(groupName)
	if err != nil {
		return nil, err
	}
	if ActionIsObtain(actionName) {
		groupItem, err = getGroup(*groupExists, groupObject.ID, groupObject.Name)
		if err != nil {
			return nil, err
		}
	} else if ActionIsCreate(actionName) {
		groupItem, err = createGroup(*groupExists, groupObject.ID, groupObject.Name, systemWildcard)
		if err != nil {
			return nil, err
		}
	} else if ActionIsDelete(actionName) {
		groupItem, err = deleteGroup(*groupExists, groupObject.ID, groupObject.Name)
		if err != nil {
			return nil, err
		}
	}
	return groupItem, nil
}

func getGroup(groupExists bool, groupId int, groupName string) (*Item, error) {
	if groupExists {
		return NewItem(groupId, "Group", groupName, false, false), nil
	} else {
		err := errors.New("Error getting group: Group with name " + groupName + " doesn't exist")
		return nil, err
	}
}

func createGroup(groupExists bool, groupId int, groupName string, systemWildcard string) (*Item, error) {
	if groupExists {
		return NewItem(groupId, "Group", groupName, false, false), nil
	} else {
		papertrailGroupCreated, err := createPapertrailGroupAction(groupName, systemWildcard)
		if err != nil {
			return nil, err
		}
		return NewItem(papertrailGroupCreated.ID, "Group", papertrailGroupCreated.Name, true, false), nil
	}
}

func deleteGroup(groupExists bool, groupId int, groupName string) (*Item, error){
	if groupExists {
		papertrailGroupDeleted, err := deletePapertrailGroupOperation(groupName, groupId)
		if err != nil {
			return nil, err
		}
		if *papertrailGroupDeleted {
			return NewItem(groupId, "Group", groupName, false, true), nil
		}
	}
	err := errors.New("Error: Group with " + groupName + " doesn't exist")
	return nil, err
}

// checkGroupExists checks if a group exists in papertrail, returning the information of this one in case it exists
func checkGroupExists(groupName string) (*bool, *GroupObject, error) {
	alreadyExists := false
	var group *GroupObject
	getAllGroupResp, err  := ApiOperation("GET", papertrailApiGroupsEndpoint, nil)
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
func createPapertrailGroupAction(groupName string, systemWildcard string) (*GroupObject, error){
	papertrailGroupToCreate := GroupCreationObject{Group: GroupCreateObject{
		Name:           groupName,
		SystemWildcard: systemWildcard,
	}}
	b, err := json.Marshal(papertrailGroupToCreate)
	if err != nil {
		return nil, err
	}
	createGroupResp, err  := ApiOperation("POST", papertrailApiGroupsEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	if createGroupResp.StatusCode == 200 {
		var group GroupObject
		json.Unmarshal([]byte(createGroupResp.Body), &group)
		log.Printf("Group with name %s and id %d was successfully created\n", group.Name, group.ID)
		return &group, nil
	}
	log.Printf("Problems creating group with name %s\n", groupName)
	err = errors.New("Error: Response status code " + strconv.Itoa(createGroupResp.StatusCode))
	return nil, err
}

// deletePapertrailGroup deletes a papertrail group using the groupId
// provided as the group information to be deleted
func deletePapertrailGroupOperation(groupName string, groupId int) (*bool, error){
	deleted := false
	groupIdUrl := strings.SplitAfter(papertrailApiGroupsEndpoint, "groups")[0] +
		"/" + strconv.Itoa(groupId) + strings.SplitAfter(papertrailApiGroupsEndpoint, "groups")[1]
	deleteGroupResp, err  := ApiOperation("DELETE", groupIdUrl, nil)
	if err != nil {
		return nil, err
	}
	if deleteGroupResp.StatusCode == 200 {
		deleted = true
		log.Printf("Group with name %s and id %d was successfully deleted\n", groupName, groupId)
		return &deleted, nil
	}
	log.Printf("Problems deleting group with id %d\n", groupId)
	err = errors.New("Error: Response status code " + strconv.Itoa(deleteGroupResp.StatusCode))
	return &deleted, err
}