package papertrail

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// papertrailApiDestinationsEndpoint represents the endpoint for interact with
// groups in papertrail API
const papertrailApiEventsSearchEndpoint = papertrailApiBaseUrl + "events/search.json"

// doPapertrailGroupNecessaryActions is in charge of get the logs
// on the indicated papertrail search and save it in a file
func doPapertrailEventsSearch(groupName string, groupId int, searchName string, searchQuery string,
	startDateUnix int64, endDateUnix int64, path string) (*Item, error) {
	var getEventsSearchResp *ApiResponse
	var eventsSearchItem *Item
	numOfEvents := 0
	pathFileName := CreateFilenameForEventsSearch(path, groupName, searchName, startDateUnix, endDateUnix)
	papertrailEventsSearch := NewEventsSearchRequestWithMinAndMaxTime(groupId, searchQuery, strconv.Itoa(int(startDateUnix)), strconv.Itoa(int(endDateUnix)))
	b, err := json.Marshal(papertrailEventsSearch)
	if err != nil {
		return nil, err
	}
	getEventsSearchResp, err = apiOperation("GET", papertrailApiEventsSearchEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	if getEventsSearchResp != nil && getEventsSearchResp.StatusCode == 200 {
		var eventsMessages []string
		file, err := os.OpenFile(pathFileName, os.O_CREATE|os.O_WRONLY, 0644)
		defer file.Close()
		if err != nil {
			return nil, err
		}
		var eventsSearch EventsSearch
		json.Unmarshal([]byte(getEventsSearchResp.Body), &eventsSearch)
		if len(eventsSearch.Events) > 0 {
			for _, event := range eventsSearch.Events {
				eventsMessages = append(eventsMessages, event.Message)
			}
			if eventsSearch.MinTimeAt.Unix() > startDateUnix {
				maxId := eventsSearch.MinID
				eventsMessages, err = getPapertrailEventsSearchIterations(groupId, searchQuery, maxId, startDateUnix, eventsMessages)
				if err != nil {
					return nil, err
				}
			}
			saveLogsToFile(*file, eventsMessages)
		}
		if eventsMessages != nil {
			numOfEvents = len(eventsMessages)
		}
		eventsSearchItem = NewItem(0, "EventsSearch", getNameOfFileLogsSaved(pathFileName)+
			" with "+strconv.Itoa(numOfEvents)+" events retrieved", false, false)
	} else {
		err := convertStatusCodeToError(getEventsSearchResp.StatusCode, "EventsSearch", "Obtaining")
		return nil, err
	}
	return eventsSearchItem, nil
}

// getPapertrailEventsSearchIterations takes care of obtaining events in the specified search
// when more than one iteration is necessary, since it changes the struct used to send as body
func getPapertrailEventsSearchIterations(groupId int, searchQuery string, maxId string,
	startDateUnix int64, eventsMessages []string) ([]string, error) {
	var eventsSearchIt EventsSearch
	for {
		papertrailEventsSearchIt := NewEventsSearchRequestWithMinTimeMaxId(groupId, searchQuery, maxId, strconv.Itoa(int(startDateUnix)))
		b, err := json.Marshal(papertrailEventsSearchIt)
		if err != nil {
			return nil, err
		}
		getEventsSearchItResp, err := apiOperation("GET", papertrailApiEventsSearchEndpoint, bytes.NewBuffer(b))
		if err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(getEventsSearchItResp.Body), &eventsSearchIt)
		for index := len(eventsSearchIt.Events) - 1; index >= 1; index-- {
			eventsMessages = append([]string{eventsSearchIt.Events[index].Message}, eventsMessages...)
		}
		if eventsSearchIt.MinTimeAt.Unix() <= startDateUnix {
			break
		}
		maxId = eventsSearchIt.MinID
	}
	return eventsMessages, nil
}

// saveLogsToFile takes care of saving all events in the specified
// string slice in the received file as a parameter
func saveLogsToFile(file os.File, eventsSearch []string) {
	for _, event := range eventsSearch {
		file.WriteString(event + "\n")
	}
}

// CreateFilenameForEventsSearch creates the name of the file where to save the log events from the received parameters
func CreateFilenameForEventsSearch(path string, groupName string, searchName string, startDateUnix int64, endDateUnix int64) string {
	groupNameFixedChars := strings.Replace(groupName, " ", "_", -1)
	searchNameFixedChars := strings.Replace(searchName, " ", "_", -1)
	startDateUnixFixedChars := strings.Replace(GetTimeInUTCFromUnixTime(startDateUnix), " ", "_", -1)
	startDateUnixFixedChars = strings.Replace(startDateUnixFixedChars, "/", "-", -1)
	endDateUnixFixedChars := strings.Replace(GetTimeInUTCFromUnixTime(endDateUnix), " ", "_", -1)
	endDateUnixFixedChars = strings.Replace(endDateUnixFixedChars, "/", "-", -1)
	return path + string(filepath.Separator) + groupNameFixedChars + "_" + searchNameFixedChars + "_" + startDateUnixFixedChars + "_" + endDateUnixFixedChars
}

// getNameOfFileLogsSaved performs the escape of special characters that
// can cause errors in the file name where the log events will be saved
func getNameOfFileLogsSaved(pathName string) string {
	filesPathName := pathName
	if strings.Contains(filesPathName, "]") {
		filesPathName = strings.ReplaceAll(filesPathName, "]", "\\]")
	}
	if strings.Contains(filesPathName, "[") {
		filesPathName = strings.ReplaceAll(filesPathName, "[", "\\[")
	}
	return filesPathName
}
