package papertrail

import "time"

// Options contains all the app possible options.
type Options struct {

	// Group name defined or to be defined in papertrail
	GroupName string

	// Wildcard to be applied on the systems defined in papertrail
	SystemWildcard string

	// Destination port for sending the logs of the indicated system/s
	DestinationPort int

	// Destination id for sending the logs of the indicated system/s
	DestinationId int

	// Source ip address from sending the logs of the indicated system/s
	IpAddress string

	// System type, can be hostname or IPAddress
	SystemType string

	// Name of saved search to be performed on logs or to be created on a group
	Search string

	// Query to be performed on the group of logs or applied on the search to be created
	Query string

	// Action to be performed with the information provided for papertrail, possible values only c(create) and o(obtain)
	Action string

	// Indicates if all searches in a group or a specific search are to be deleted
	DeleteAllSearches bool

	// Filter only from a specific date
	StartDate string

	// Filter only until a specific date
	EndDate string

	// Path where to store the logs
	Path string
}

type Self struct {
	Href string `json:"href"`
}

type HTML struct {
	Href string `json:"href"`
}

type Search struct {
	Href string `json:"href"`
}

type Links struct {
	Self `json:"self"`
	HTML `json:"html"`
	Search `json:"search"`
}

type Syslog struct {
	Hostname    string      `json:"hostname"`
	Port        int         `json:"port"`
	Description interface{} `json:"description"`
}

type Destination struct {
	ID     int         `json:"id"`
	Filter interface{} `json:"filter"`
	Syslog Syslog 		`json:"syslog"`
}

type System struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	LastEventAt time.Time `json:"last_event_at"`
	AutoDelete  bool      `json:"auto_delete"`
	Links				  `json:"_links"`
	IPAddress interface{} `json:"ip_address"`
	Hostname  string      `json:"hostname"`
	Syslog				  `json:"syslog"`
}

func NewSystem(ID int64, name string, lastEventAt time.Time, autoDelete bool, links Links, IPAddress interface{}, hostname string, syslog Syslog) *System {
	return &System{ID: ID, Name: name, LastEventAt: lastEventAt, AutoDelete: autoDelete, Links: links, IPAddress: IPAddress, Hostname: hostname, Syslog: syslog}
}

type GroupCreateObject struct {
	Name           string `json:"name"`
	SystemWildcard string `json:"system_wildcard"`
}

type GroupCreationObject struct {
	Group GroupCreateObject `json:"group"`
}

type GroupObject struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	SystemWildcard string `json:"system_wildcard"`
	Links          `json:"_links"`
	Systems []System `json:"systems"`
}

func NewGroupObject(ID int, name string, systemWildcard string, links Links, systems []System) *GroupObject {
	return &GroupObject{ID: ID, Name: name, SystemWildcard: systemWildcard, Links: links, Systems: systems}
}

type SearchToCreate struct {
	Name    string `json:"name"`
	Query   string `json:"query"`
	GroupID int `json:"group_id"`
}

type SearchToCreateObject struct {
	SearchToCreate `json:"search"`
}

type SearchHref struct {
	Href string `json:"href"`
}

type HTMLSearchHref struct {
	Href string `json:"href"`
}

type SearchGroupLinks struct {
	Self Self `json:"self"`
	HTML HTML `json:"html"`
	Search SearchHref `json:"search"`
	HTMLSearch HTMLSearchHref `json:"html_search"`
}

type SearchGroup struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Links  SearchGroupLinks `json:"_links"`
}

type SearchObject struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Query string `json:"query"`
	Group SearchGroup `json:"group"`
	Links SearchGroupLinks `json:"_links"`
}

func NewSearchObject(ID int, name string, query string, group SearchGroup, links SearchGroupLinks) *SearchObject {
	return &SearchObject{ID: ID, Name: name, Query: query, Group: group, Links: links}
}

type ApiResponse struct {
	Body 		[]byte
	StatusCode	int
	err			error
}

type Item struct {
	ID 			int
	ItemType	string
	ItemName	string
	Created		bool
	Deleted		bool
}

func NewItem(ID int, itemType string, itemName string, created bool, deleted bool) *Item {
	return &Item{ID: ID, ItemType: itemType, ItemName: itemName, Created: created, Deleted: deleted}
}

type SystemBasedInHostname struct {
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
}

type SystemToCreateBasedInHostnameToDestinationID struct {
	System SystemBasedInHostname `json:"system"`
	DestinationID int `json:"destination_id"`
}

func NewSystemToCreateBasedInHostnameToDestinationID(system SystemBasedInHostname, destinationID int) *SystemToCreateBasedInHostnameToDestinationID {
	return &SystemToCreateBasedInHostnameToDestinationID{System: system, DestinationID: destinationID}
}

type SystemToCreateBasedInHostnameToDestinationPort struct {
	System SystemBasedInHostname `json:"system"`
	DestinationPort int `json:"destination_port"`
}

func NewSystemToCreateBasedInHostnameToDestinationPort(system SystemBasedInHostname, destinationPort int) *SystemToCreateBasedInHostnameToDestinationPort {
	return &SystemToCreateBasedInHostnameToDestinationPort{System: system, DestinationPort: destinationPort}
}

type SystemBasedInIPAddress struct {
	Name     	string `json:"name"`
	IPAddress 	string `json:"ip_address"`
}

type SystemToCreateBasedInIpAddress struct {
	System SystemBasedInIPAddress `json:"system"`
}

func NewSystemToCreateBasedInIpAddress(system SystemBasedInIPAddress) *SystemToCreateBasedInIpAddress {
	return &SystemToCreateBasedInIpAddress{System: system}
}

type Events struct {
	ID                string    `json:"id"`
	SourceIP          string    `json:"source_ip"`
	Program           string    `json:"program"`
	Message           string    `json:"message"`
	ReceivedAt        time.Time `json:"received_at"`
	GeneratedAt       time.Time `json:"generated_at"`
	DisplayReceivedAt string    `json:"display_received_at"`
	SourceID          int64     `json:"source_id"`
	SourceName        string    `json:"source_name"`
	Hostname          string    `json:"hostname"`
	Severity          string    `json:"severity"`
	Facility          string    `json:"facility"`
}

type EventsSearch struct {
	MinID  string `json:"min_id"`
	MaxID  string `json:"max_id"`
	Events []Events 			 `json:"events"`
	Sawmill            bool      `json:"sawmill"`
	ReachedBeginning   bool      `json:"reached_beginning"`
	MinTimeAt          time.Time `json:"min_time_at"`
	ReachedRecordLimit bool      `json:"reached_record_limit"`
}

type EventsSearchRequestWithMinAndMaxTime struct {
	GroupID int    `json:"group_id"`
	Q       string `json:"q"`
	MinTime string `json:"min_time"`
	MaxTime string `json:"max_time"`
}

func NewEventsSearchRequestWithMinAndMaxTime(groupID int, q string, minTime string, maxTime string) *EventsSearchRequestWithMinAndMaxTime {
	return &EventsSearchRequestWithMinAndMaxTime{GroupID: groupID, Q: q, MinTime: minTime, MaxTime: maxTime}
}

type EventsSearchRequestWithMinTimeMaxId struct {
	GroupID int    `json:"group_id"`
	Q       string `json:"q"`
	MinTime string `json:"min_time"`
	MaxId	string	`json:"max_id"`
}

func NewEventsSearchRequestWithMinTimeMaxId(groupID int, q string, minTime string, maxId string) *EventsSearchRequestWithMinTimeMaxId {
	return &EventsSearchRequestWithMinTimeMaxId{GroupID: groupID, Q: q, MinTime: minTime, MaxId: maxId}
}