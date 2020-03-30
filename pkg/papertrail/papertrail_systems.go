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
const papertrailApiSystemsEndpoint = papertrailApiBaseUrl + "systems.json"

// getSystemInPapertrail obtains a papertrail system, creating it in case it does not exist previously
func getSystemInPapertrailBasedInHostname(hostname string, destinationPort int, destinationId int) (*Item, error) {
	systemExists, systemObject, err := checkSystemExistsBasedInHostname(hostname, destinationPort, destinationId)
	if err != nil {
		return nil, err
	}
	var systemItem *Item
	if (systemExists != nil) && *systemExists {
		fmt.Printf("System with hostname %s already exists with id %d\n", hostname, systemObject.ID)
		systemItem = NewItem(int(systemObject.ID), "System", systemObject.Name, false)
	} else {
		fmt.Printf("System with hostname %s doesn't exist yet\n", hostname)
		var papertrailSystemCreated *System
		if destinationPort != 0 {
			papertrailSystemCreated, err = createPapertrailSystemBasedInHostnameAndDestinationPort(hostname, destinationPort)
		} else {
			papertrailSystemCreated, err = createPapertrailSystemBasedInHostnameAndDestinationId(hostname, destinationId)
		}
		if err != nil {
			return nil, err
		}
		systemItem = NewItem(int(papertrailSystemCreated.ID), "System", papertrailSystemCreated.Name, true)
	}
	return systemItem, err
}

// getSystemInPapertrail obtains a papertrail system, creating it in case it does not exist previously
func getSystemInPapertrailBasedInAddressIp(addressIP string) (*Item, error) {
	systemExists, systemObject, err := checkSystemExistsBasedInAddressIP(addressIP)
	if err != nil {
		return nil, err
	}
	var systemItem *Item
	if (systemExists != nil) && *systemExists {
		fmt.Printf("System with IPAddress %s already exists with id %d\n", addressIP, systemObject.ID)
		systemItem = NewItem(int(systemObject.ID), "System", systemObject.Name, false)
	} else {
		fmt.Printf("System with IPAddress %s doesn't exist yet\n", addressIP)
		papertrailSystemCreated, err := createPapertrailSystemBasedInIPAddress(addressIP)
		if err != nil {
			return nil, err
		}
		systemItem = NewItem(int(papertrailSystemCreated.ID), "System", papertrailSystemCreated.Name, true)
	}
	return systemItem, err
}

// checkSystemExists checks if a system exists in papertrail, returning the information of this one in case it exists
func checkSystemExistsBasedInHostname(hostname string, destinationPort int, destinationId int) (*bool, *System, error) {
	getAllSystems, err  := papertrailApiOperation("GET", papertrailApiSystemsEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	var system *System
	var alreadyExists *bool
	if getAllSystems.StatusCode == 200 {
		var systems []System
		json.Unmarshal([]byte(getAllSystems.Body), &systems)
		if destinationPort != 0 {
			alreadyExists, system = checkSystemExistsBasedInHostnameAndDestinationPort(systems, hostname, destinationPort)

		} else if destinationId != 0 {
			alreadyExists, system = checkSystemExistsBasedInHostnameAndDestinationId(systems, hostname, destinationId)
		}
	}
	return alreadyExists, system, nil
}

// checkSystemExists checks if a system exists in papertrail, returning the information of this one in case it exists
func checkSystemExistsBasedInAddressIP(addressIP string) (*bool, *System, error) {
	getAllSystems, err  := papertrailApiOperation("GET", papertrailApiSystemsEndpoint, nil)
	alreadyExists := false
	if err != nil {
		return nil, nil, err
	}
	var system *System
	if getAllSystems.StatusCode == 200 {
		var systems []System
		json.Unmarshal([]byte(getAllSystems.Body), &systems)
		for _, item := range systems {
			if item.IPAddress == addressIP {
				alreadyExists = true
				system = NewSystem(item.ID, item.Name, item.LastEventAt,
					item.AutoDelete, item.Links, item.IPAddress, item.Hostname, item.Syslog)
				break
			}
		}
	}
	return &alreadyExists, system, nil
}

func checkSystemExistsBasedInHostnameAndDestinationPort(systems []System, hostname string, destinationPort int) (*bool, *System){
	var system *System
	alreadyExists := false
	for _, item := range systems {
		if item.Hostname == hostname && item.Port == destinationPort {
			alreadyExists = true
			system = NewSystem(
							item.ID,
							item.Name,
							item.LastEventAt,
							item.AutoDelete,
							item.Links,
							item.IPAddress,
							item.Hostname,
							item.Syslog)
			break
		}
	}
	return &alreadyExists, system
}

func checkSystemExistsBasedInHostnameAndDestinationId(systems []System, hostname string, destinationId int) (*bool, *System){
	var system *System
	alreadyExists := false
	for _, item := range systems {
		if item.Hostname == hostname && int(item.ID) == destinationId {
			system = NewSystem(
				item.ID,
				item.Name,
				item.LastEventAt,
				item.AutoDelete,
				item.Links,
				item.IPAddress,
				item.Hostname,
				item.Syslog)
			break
		}
	}
	return &alreadyExists, system
}

func SystemToCreateBasedInHostnameAndDestinationId(hostname string, destinationId int) *SystemToCreateBasedInHostnameToDestinationID {
	return NewSystemToCreateBasedInHostnameToDestinationID(SystemBasedInHostname{
		Name:     hostname,
		Hostname: hostname,
	}, destinationId)
}

func SystemToCreateBasedInHostnameAndDestinationPort(hostname string, destinationPort int) *SystemToCreateBasedInHostnameToDestinationPort {
	return NewSystemToCreateBasedInHostnameToDestinationPort(SystemBasedInHostname{
		Name:     hostname,
		Hostname: hostname,
	}, destinationPort)
}



// createPapertrailSystemBasedInHostnameAndDestinationId creates a papertrail
// system using the parameter information provided as the group information to be created
func createPapertrailSystemBasedInHostnameAndDestinationId(hostname string, destinationId int) (*System, error){
	papertrailSystemToCreate := SystemToCreateBasedInHostnameAndDestinationId(hostname, destinationId)
	b, err := json.Marshal(papertrailSystemToCreate)
	fmt.Printf("Body is %s\n", b)
	if err != nil {
		return nil, err
	}
	createSystemResp, err  := papertrailApiOperation("POST", papertrailApiSystemsEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	if createSystemResp.StatusCode == 200 {
		var system System
		json.Unmarshal([]byte(createSystemResp.Body), &system)
		fmt.Printf("System with name %s based in hostname %s was successfully " +
			"created with id %d\n", system.Name, system.Hostname, system.ID)
		return &system, nil
	}
	fmt.Printf("Problems creating system with name %s and hostname%s\n", hostname, hostname)
	err = errors.New("Error: Response status code " + strconv.Itoa(createSystemResp.StatusCode))
	return nil, err
}

// createPapertrailSystemBasedInHostnameAndDestinationPort creates a papertrail group using the parameter information
// provided as the system information to be created
func createPapertrailSystemBasedInHostnameAndDestinationPort(hostname string, destinationPort int) (*System, error){
	papertrailSystemToCreate := SystemToCreateBasedInHostnameAndDestinationPort(hostname, destinationPort)
	b, err := json.Marshal(papertrailSystemToCreate)
	if err != nil {
		return nil, err
	}
	createSystemResp, err  := papertrailApiOperation("POST", papertrailApiSystemsEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	if createSystemResp.StatusCode == 200 {
		var system System
		json.Unmarshal([]byte(createSystemResp.Body), &system)
		fmt.Printf("System with name %s based in hostname %s was successfully " +
			"created with id %d\n", system.Name, system.Hostname, system.ID)
		return &system, nil
	}
	fmt.Printf("Problems creating system with name %s and hostname%s\n", hostname, hostname)
	err = errors.New("Error: Response status code " + strconv.Itoa(createSystemResp.StatusCode))
	return nil, err
}

// createPapertrailSystemBasedInIPAddress creates a papertrail system using the parameter information
// provided as the system information to be created
func createPapertrailSystemBasedInIPAddress(ipAddress string) (*System, error){
	papertrailSystemToCreate := NewSystemToCreateBasedInIpAddress(SystemBasedInIPAddress{
		Name:      ipAddress,
		IPAddress: ipAddress,
	})
	b, err := json.Marshal(papertrailSystemToCreate)
	if err != nil {
		return nil, err
	}
	createSystemResp, err  := papertrailApiOperation("POST", papertrailApiSystemsEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	if createSystemResp.StatusCode == 200 {
		var system System
		json.Unmarshal([]byte(createSystemResp.Body), &system)
		fmt.Printf("System with name %s and IPAddress %s " +
			"was successfully created with id %d\n", system.Name, system.IPAddress, system.ID)
		return &system, nil
	}
	fmt.Printf("Problems creating system with name %s and IPAddress %s\n", ipAddress, ipAddress)
	err = errors.New("Error: Response status code " + strconv.Itoa(createSystemResp.StatusCode))
	return nil, err
}