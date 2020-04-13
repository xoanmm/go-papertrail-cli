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

// doPapertrailGroupNecessaryActions is in charge of carrying out the indicated actions
// on the indicated papertrail group, as well as checking if it exists
func doPapertrailGroupNecessaryActions(groupName string, actionName string, systemWildcard string) (*Item, error) {
	var groupItem *Item
	groupObject, err := checkGroupExists(groupName)
	if err != nil {
		return nil, err
	}
	if groupObject != nil {
		log.Printf("Group with name %s exists with id %d\n", groupName, groupObject.ID)
		if ActionIsObtain(actionName) || ActionIsCreate(actionName) {
			return NewItem(groupObject.ID, "Group", groupName, false, false), nil
		} else if ActionIsDelete(actionName) {
			groupItem, err = deleteGroup(groupObject.ID, groupObject.Name)
			if err != nil {
				return nil, err
			}
		}
	} else {
		if ActionIsCreate(actionName) {
			groupItem, err = createGroup(groupName, systemWildcard)
			if err != nil {
				return nil, err
			}
		} else {
			err := errors.New("Error: Group with name " + groupName + " doesn't exist ")
			return nil, err
		}
	}
	return groupItem, nil
}

// createGroup attempts to create a group in papertrail using the parameters provided as group information
func createGroup(groupName string, systemWildcard string) (*Item, error) {
	papertrailGroupCreated, err := createPapertrailGroupOperation(groupName, systemWildcard)
	if err != nil {
		return nil, err
	}
	return NewItem(papertrailGroupCreated.ID, "Group", papertrailGroupCreated.Name, true, false), nil
}

// deleteGroup attempts to delete a group using the parameters provided as group information
func deleteGroup(groupId int, groupName string) (*Item, error) {
	papertrailGroupDeleted, err := deletePapertrailGroupOperation(groupName, groupId)
	if err != nil {
		return nil, err
	}
	if *papertrailGroupDeleted {
		return NewItem(groupId, "Group", groupName, false, true), nil
	}
	return nil, errors.New("Error: Group with " + groupName + " doesn't exist")
}

// checkGroupExists checks if a group exists in papertrail, returning the information of this one in case it exists
func checkGroupExists(groupName string) (*GroupObject, error) {
	var group *GroupObject
	getAllGroupResp, err := apiOperation("GET", papertrailApiGroupsEndpoint, nil)
	if err != nil {
		return group, err
	}
	if getAllGroupResp.StatusCode == 200 {
		var groups []GroupObject
		json.Unmarshal([]byte(getAllGroupResp.Body), &groups)
		for _, item := range groups {
			if item.Name == groupName {
				group = NewGroupObject(item.ID, item.Name, item.SystemWildcard, item.Links, item.Systems)
				break
			}
		}
	}
	return group, nil
}

// createPapertrailGroupOperation do the necessary calls in papertrail
// to create a group using the parameter information provided as the group information to be created
func createPapertrailGroupOperation(groupName string, systemWildcard string) (*GroupObject, error) {
	papertrailGroupToCreate := GroupCreationObject{Group: GroupCreateObject{
		Name:           groupName,
		SystemWildcard: systemWildcard,
	}}
	b, err := json.Marshal(papertrailGroupToCreate)
	if err != nil {
		return nil, err
	}
	createGroupResp, err := apiOperation("POST", papertrailApiGroupsEndpoint, bytes.NewBuffer(b))
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
	err = convertStatusCodeToError(createGroupResp.StatusCode, "Group", "Creating")
	return nil, err
}

// deletePapertrailGroupOperation do the necessary calls in papertrail
// to delete a group using the parameter information provided as the group information to be deleted
func deletePapertrailGroupOperation(groupName string, groupId int) (*bool, error) {
	deleted := false
	groupIdUrl := strings.SplitAfter(papertrailApiGroupsEndpoint, "groups")[0] +
		"/" + strconv.Itoa(groupId) + strings.SplitAfter(papertrailApiGroupsEndpoint, "groups")[1]
	deleteGroupResp, err := apiOperation("DELETE", groupIdUrl, nil)
	if err != nil {
		return nil, err
	}
	if deleteGroupResp.StatusCode == 200 {
		deleted = true
		log.Printf("Group with name %s and id %d was successfully deleted\n", groupName, groupId)
		return &deleted, nil
	}
	log.Printf("Problems deleting group with id %d\n", groupId)
	err = convertStatusCodeToError(deleteGroupResp.StatusCode, "Group", "Deleting")
	return &deleted, err
}
