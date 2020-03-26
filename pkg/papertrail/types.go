package papertrail

import "time"

// Options contains all the app possible options.
type Options struct {

	// Group name defined or to be defined in papertrail
	GroupName string

	// Wildcard to be applied on the systems defined in papertrail
	SystemWildcard string

	// Name of saved search to be performed on logs or to be created on a group
	Search string

	// Query to be performed on the group of logs or applied on the search to be created
	Query string

	// Action to be performed with the information provided for papertrail, possible values only c(create) and o(obtain)
	Action string
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

type Syslog    struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
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
}