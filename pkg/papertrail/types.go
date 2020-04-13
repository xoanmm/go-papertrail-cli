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

// Self object used by papertrail to identify a Self object
type Self struct {
	Href string `json:"href"`
}

// HTML object used by papertrail to identify a HTML object
type HTML struct {
	Href string `json:"href"`
}

// Search object used by papertrail to identify a Search object
type Search struct {
	Href string `json:"href"`
}

// Links object used by papertrail to group Self, HTML and Search objects
type Links struct {
	Self   `json:"self"`
	HTML   `json:"html"`
	Search `json:"search"`
}

// Syslog object used by papertrail to identify a Syslog object
type Syslog struct {
	Hostname    string      `json:"hostname"`
	Port        int         `json:"port"`
	Description interface{} `json:"description"`
}

// Destination object used by papertrail to identify a Destination object
type Destination struct {
	ID     int         `json:"id"`
	Filter interface{} `json:"filter"`
	Syslog Syslog      `json:"syslog"`
}

// System object used by papertrail to identify a System object
type System struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	LastEventAt time.Time `json:"last_event_at"`
	AutoDelete  bool      `json:"auto_delete"`
	Links       `json:"_links"`
	IPAddress   interface{} `json:"ip_address"`
	Hostname    string      `json:"hostname"`
	Syslog      `json:"syslog"`
}

// NewSystem allows to create a System type struct providing all the information for it
func NewSystem(ID int64, name string, lastEventAt time.Time, autoDelete bool, links Links, IPAddress interface{}, hostname string, syslog Syslog) *System {
	return &System{ID: ID, Name: name, LastEventAt: lastEventAt, AutoDelete: autoDelete, Links: links, IPAddress: IPAddress, Hostname: hostname, Syslog: syslog}
}

// GroupCreateObject is the structure used to collect information about a Group created in papertrail
type GroupCreateObject struct {
	Name           string `json:"name"`
	SystemWildcard string `json:"system_wildcard"`
}

// GroupCreationObject is the structure used to send information about a Group to be created on papertrail
type GroupCreationObject struct {
	Group GroupCreateObject `json:"group"`
}

// GroupObject is the representation of a papertrail's Group
type GroupObject struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	SystemWildcard string `json:"system_wildcard"`
	Links          `json:"_links"`
	Systems        []System `json:"systems"`
}

// NewGroupObject allows to create a Group type struct providing all the information for it
func NewGroupObject(ID int, name string, systemWildcard string, links Links, systems []System) *GroupObject {
	return &GroupObject{ID: ID, Name: name, SystemWildcard: systemWildcard, Links: links, Systems: systems}
}

// SearchToCreate is the information used to send information about a Search to be created on papertrail
type SearchToCreate struct {
	Name    string `json:"name"`
	Query   string `json:"query"`
	GroupID int    `json:"group_id"`
}

// SearchToCreateObject is the structure used to send information about a Search to be created on papertrail
type SearchToCreateObject struct {
	SearchToCreate `json:"search"`
}

// SearchHref is the object used by papertrail to identify a Href information of a Search object
type SearchHref struct {
	Href string `json:"href"`
}

// HTMLSearchHref is the object used by papertrail to identify a Href information of a Search's HTML object
type HTMLSearchHref struct {
	Href string `json:"href"`
}

// SearchGroupLinks is the object used by papertrail to identify a Link of a Search in a Group
type SearchGroupLinks struct {
	Self       Self           `json:"self"`
	HTML       HTML           `json:"html"`
	Search     SearchHref     `json:"search"`
	HTMLSearch HTMLSearchHref `json:"html_search"`
}

// SearchGroup is the object used by papertrail to identify a Link of a Search in a Group
type SearchGroup struct {
	ID    int              `json:"id"`
	Name  string           `json:"name"`
	Links SearchGroupLinks `json:"_links"`
}

// SearchObject object used by papertrail to identify a Search object
type SearchObject struct {
	ID    int              `json:"id"`
	Name  string           `json:"name"`
	Query string           `json:"query"`
	Group SearchGroup      `json:"group"`
	Links SearchGroupLinks `json:"_links"`
}

// NewSearchObject allows to create a Search type struct providing all the information for it
func NewSearchObject(ID int, name string, query string, group SearchGroup, links SearchGroupLinks) *SearchObject {
	return &SearchObject{ID: ID, Name: name, Query: query, Group: group, Links: links}
}

// ApiResponse represents the information collected about a response of papertrail's API requests
type ApiResponse struct {
	Body       []byte
	StatusCode int
	err        error
}

// Item is the structure used to represent the different papertrail elements
// on which some action is performed during the execution of the cli
type Item struct {
	ID       int
	ItemType string
	ItemName string
	Created  bool
	Deleted  bool
}

// NewItem allows to create a Item type struct providing all the information for it
func NewItem(ID int, itemType string, itemName string, created bool, deleted bool) *Item {
	return &Item{ID: ID, ItemType: itemType, ItemName: itemName, Created: created, Deleted: deleted}
}

// SystemBasedInHostname is the structure used to represent the information
// of a hostname based papertrail system
type SystemBasedInHostname struct {
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
}

// SystemToCreateBasedInHostnameToDestinationID is the structure used to represent a system to be
// created based on hostname and a destination identifier to be used
type SystemToCreateBasedInHostnameToDestinationID struct {
	System        SystemBasedInHostname `json:"system"`
	DestinationID int                   `json:"destination_id"`
}

// NewSystemToCreateBasedInHostnameToDestinationID allows to create a SystemToCreateBasedInHostnameToDestinationID type
// struct providing all the information for it
func NewSystemToCreateBasedInHostnameToDestinationID(system SystemBasedInHostname, destinationID int) *SystemToCreateBasedInHostnameToDestinationID {
	return &SystemToCreateBasedInHostnameToDestinationID{System: system, DestinationID: destinationID}
}

// SystemToCreateBasedInHostnameToDestinationPort is the structure used to represent a system to be
// created based on hostname and a destination port to be used
type SystemToCreateBasedInHostnameToDestinationPort struct {
	System          SystemBasedInHostname `json:"system"`
	DestinationPort int                   `json:"destination_port"`
}

// NewSystemToCreateBasedInHostnameToDestinationPort allows to create a SystemToCreateBasedInHostnameToDestinationPort type
// struct providing all the information for it
func NewSystemToCreateBasedInHostnameToDestinationPort(system SystemBasedInHostname, destinationPort int) *SystemToCreateBasedInHostnameToDestinationPort {
	return &SystemToCreateBasedInHostnameToDestinationPort{System: system, DestinationPort: destinationPort}
}

// SystemBasedInIPAddress is the structure used to represent a system to be
// created based on hostname and an IP Address to be used
type SystemBasedInIPAddress struct {
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
}

// SystemToCreateBasedInIpAddress is the structure used to represent a system to be
// created based on a SystemBasedInIPAddress object
type SystemToCreateBasedInIpAddress struct {
	System SystemBasedInIPAddress `json:"system"`
}

// NewSystemToCreateBasedInIpAddress allows to create a SystemToCreateBasedInIpAddress type
// struct providing all the information for it
func NewSystemToCreateBasedInIpAddress(system SystemBasedInIPAddress) *SystemToCreateBasedInIpAddress {
	return &SystemToCreateBasedInIpAddress{System: system}
}

// Events is the structure used to represent the information of the events obtained in a papertrail search
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

// EventsSearch represents the event information obtained from a papertrail search
type EventsSearch struct {
	MinID              string    `json:"min_id"`
	MaxID              string    `json:"max_id"`
	Events             []Events  `json:"events"`
	Sawmill            bool      `json:"sawmill"`
	ReachedBeginning   bool      `json:"reached_beginning"`
	MinTimeAt          time.Time `json:"min_time_at"`
	ReachedRecordLimit bool      `json:"reached_record_limit"`
}

// EventsSearchRequestWithMinAndMaxTime represents the information used to request events
// from a search without providing the parameters for minimum and maximum time
type EventsSearchRequestWithMinAndMaxTime struct {
	GroupID int    `json:"group_id"`
	Q       string `json:"q"`
	MinTime string `json:"min_time"`
	MaxTime string `json:"max_time"`
}

// NewEventsSearchRequestWithMinAndMaxTime allows to create a EventsSearchRequestWithMinAndMaxTime type
// struct providing all the information for it
func NewEventsSearchRequestWithMinAndMaxTime(groupID int, q string, minTime string, maxTime string) *EventsSearchRequestWithMinAndMaxTime {
	return &EventsSearchRequestWithMinAndMaxTime{GroupID: groupID, Q: q, MinTime: minTime, MaxTime: maxTime}
}

// EventsSearchRequestWithMinTimeMaxId represents the information used to request events
// from a search providing the parameters for minimum time and the maximum id
type EventsSearchRequestWithMinTimeMaxId struct {
	GroupID int    `json:"group_id"`
	Q       string `json:"q"`
	MinTime string `json:"min_time"`
	MaxId   string `json:"max_id"`
}

// NewEventsSearchRequestWithMinTimeMaxId allows to create a EventsSearchRequestWithMinTimeMaxId type
// struct providing all the information for it
func NewEventsSearchRequestWithMinTimeMaxId(groupID int, q string, minTime string, maxId string) *EventsSearchRequestWithMinTimeMaxId {
	return &EventsSearchRequestWithMinTimeMaxId{GroupID: groupID, Q: q, MinTime: minTime, MaxId: maxId}
}
