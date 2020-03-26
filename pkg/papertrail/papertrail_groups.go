package papertrail

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func getGroupInPapertrail(groupName string, systemWildcard string) (*Item, error) {
	groupExists, groupObject, err := checkGroupExists(groupName)
	if err != nil {
		return nil, err
	}
	groupItem := Item{
		ItemType: "Group",
	}
	if *groupExists {
		fmt.Printf("Group with name %s already exists with id %d\n", groupName, groupObject.ID)
		groupItem.ID = groupObject.ID
		groupItem.ItemName = groupObject.Name
		groupItem.Created = false
	} else {
		fmt.Printf("Group with name %s doesn't exist yet\n", groupName)
		papertrailGroupCreated, err := createPapertrailGroup(groupName, systemWildcard)
		if err != nil {
			return nil, err
		}
		groupItem.ID = papertrailGroupCreated.ID
		groupItem.ItemName = papertrailGroupCreated.Name
		groupItem.Created = true
	}
	return &groupItem, err
}

func checkGroupExists(groupName string) (*bool, *GroupObject, error) {
	alreadyExists := false
	var group GroupObject
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
				group.ID = item.ID
				group.Name = item.Name
				group.SystemWildcard = item.SystemWildcard
				group.Links = item.Links
				group.Systems = item.Systems
				break
			}
		}
	}
	return &alreadyExists, &group, nil
}

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
	return nil, nil
}