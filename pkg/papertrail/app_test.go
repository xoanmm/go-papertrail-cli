package papertrail

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"strconv"
	"time"

	"log"
	"os"
	"testing"
)

// dateLayout define the layout to use with format
// mm/dd/yyyy HH:MM:SS
const dateLayout = "01/02/2006 15:04:05"

var now = time.Now().UTC()
var nowDate = now.Format(dateLayout)
var nowDateLessEightHours = now.Add(-8 * time.Hour).Format(dateLayout)
var papertrailApiToken string
var destinationDefaultId int
var destinationDefaultPort int

// setupEnv checks if there is an `.env' file where a series of variables
// used to perform the integration tests are defined
func setupEnv() error {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		return errors.New("Error loading .env file ")
	}
	return nil
}

func TestMain(m *testing.M) {
	// load .env file
	err := setupEnv()
	if err != nil {
		log.Fatal(err)
	}
	papertrailApiToken = os.Getenv("PAPERTRAIL_API_TOKEN")
	destinationDefaultId, err = strconv.Atoi(os.Getenv("DESTINATION_DEFAULT_ID"))
	if err != nil {
		log.Fatal(err)
	}
	destinationDefaultPort, err = strconv.Atoi(os.Getenv("DESTINATION_DEFAULT_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestApp_PapertrailActionsNoProvidedToken(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	os.Setenv("PAPERTRAIL_API_TOKEN", "")
	app := App{}
	_, _, err := app.PapertrailActions(&Options{
		"group-name",
		"*",
		0,
		7777,
		"",
		"hostname",
		"default search",
		"*",
		"c",
		false,
		false,
		true,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	})
	expectedError := errors.New("Error getting value of PAPERTRAIL_API_TOKEN, it's necessary to define this variable with your papertrail's API token ")
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestApp_PapertrailActionsInvalidHostConfiguration(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	errT := os.Setenv("PAPERTRAIL_API_TOKEN", "tokenPapertrail")
	if errT != nil {
		log.Fatal(errT)
	}
	_, _, err := app.PapertrailActions(&Options{
		"group-name",
		"*",
		0,
		0,
		"",
		"hostname",
		"default search",
		"*",
		"c",
		false,
		false,
		true,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	})
	expectedError := errors.New("It's necessary provide a value distinct from default (0) to destination id or destination port ")
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestApp_PapertrailActionsInvalidActionConfiguration(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	errT := os.Setenv("PAPERTRAIL_API_TOKEN", "tokenPapertrail")
	if errT != nil {
		log.Fatal(errT)
	}
	_, _, err := app.PapertrailActions(&Options{
		"group-name",
		"*",
		0,
		0,
		"",
		"hostname",
		"default search",
		"*",
		"ddd",
		false,
		false,
		true,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	})
	expectedError := errors.New("Not valid option provided for action to perform, the only valid values are: \n" +
		"\t'c' or 'create': create new system/s, group and/or search\n" +
		"\t'd' or 'delete': create new system/s, group and/or search\n" +
		"\t'o'or 'obtain': obtain logs in base of parameters provided\n")
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestApp_PapertrailActionsInvalidSystemProvided(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	errT := os.Setenv("PAPERTRAIL_API_TOKEN", "tokenPapertrail")
	if errT != nil {
		log.Fatal(errT)
	}
	_, _, err := app.PapertrailActions(&Options{
		"group-name",
		"*",
		0,
		7777,
		"",
		"hostnameee",
		"default search",
		"*",
		"c",
		false,
		false,
		true,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	})
	expectedError := errors.New("Not valid option provided for system, the only valid values are: \n" +
		"\t'h' or 'hostname': system based in hostname\n" +
		"\t'i'or 'ip-address': system based in ip-address\n")
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestApp_PapertrailActionsTwoSystemConfigurationProvided(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	errT := os.Setenv("PAPERTRAIL_API_TOKEN", "tokenPapertrail")
	if errT != nil {
		log.Fatal(errT)
	}
	_, _, err := app.PapertrailActions(&Options{
		"group-name",
		"*",
		7777,
		7777,
		"",
		"hostname",
		"default search",
		"*",
		"c",
		false,
		false,
		true,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	})
	expectedError := errors.New("If the system is a hostname-type system, only destination " +
		"id or destination port can be specified\n")
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestApp_PapertrailActionsTwoSystemConfigurationIncorrect(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	errT := os.Setenv("PAPERTRAIL_API_TOKEN", "tokenPapertrail")
	if errT != nil {
		log.Fatal(errT)
	}
	_, _, err := app.PapertrailActions(&Options{
		"group-name",
		"*",
		0,
		0,
		"",
		"hostname",
		"default search",
		"*",
		"c",
		false,
		false,
		true,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	})
	expectedError := errors.New("It's necessary provide a value distinct from default (0) to " +
		"destination id or destination port ")
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestApp_PapertrailActionsInvalidIpAddressConfiguration(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	errT := os.Setenv("PAPERTRAIL_API_TOKEN", "tokenPapertrail")
	if errT != nil {
		log.Fatal(errT)
	}
	_, _, err := app.PapertrailActions(&Options{
		"group-name",
		"*",
		0,
		0,
		"11111111",
		"ip-address",
		"default search",
		"*",
		"c",
		false,
		false,
		true,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	})
	expectedError := errors.New("The IP Address provided, 11111111 it's not a valid IP Address ")
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func itemEqualWithoutId(expected Item, obtained Item) bool {
	if expected.Created == obtained.Created && expected.Deleted == obtained.Deleted &&
		expected.ItemName == obtained.ItemName && expected.ItemType == obtained.ItemType {
		return true
	}
	return false
}

func EqualSlices(a, b []Item) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func testDeleteOnlySystemsHostnameDestinationPort(t *testing.T, options Options, createdElements []Item) {
	options.DeleteAllSystems = true
	options.DeleteAllSearches = true
	options.Action = "delete"
	app := &App{}
	deletedItems, _, err := app.PapertrailActions(&options)
	if err != nil {
		log.Fatal(err)
	}
	expectedDeletedSystem1 := NewItem(createdElements[0].ID, createdElements[0].ItemType,
		createdElements[0].ItemName, false, true)
	expectedDeletedSystem2 := NewItem(createdElements[1].ID, createdElements[1].ItemType,
		createdElements[1].ItemName, false, true)
	itemsDeletedExpected := []Item{
		*expectedDeletedSystem1,
		*expectedDeletedSystem2,
	}
	if !(EqualSlices(deletedItems, itemsDeletedExpected)) {
		log.Fatal("Items deleted are not equal to expected")
	}
}

func testDeleteSystemsHostnameDestinationPortGroupAndAllSearchs(t *testing.T, options Options, createdElements []Item) {
	options.DeleteAllSystems = true
	options.DeleteAllSearches = true
	options.Action = "delete"
	app := &App{}
	deletedItems, _, err := app.PapertrailActions(&options)
	if err != nil {
		log.Fatal(err)
	}
	expectedDeletedSystem1 := NewItem(createdElements[0].ID, createdElements[0].ItemType,
		createdElements[0].ItemName, false, true)
	expectedDeletedSystem2 := NewItem(createdElements[1].ID, createdElements[1].ItemType,
		createdElements[1].ItemName, false, true)
	expectedDeletedGroup := NewItem(createdElements[2].ID, createdElements[2].ItemType,
		createdElements[2].ItemName, false, true)
	itemsDeletedExpected := []Item{
		*expectedDeletedSystem1,
		*expectedDeletedSystem2,
		*expectedDeletedGroup,
	}
	if !(EqualSlices(deletedItems, itemsDeletedExpected)) {
		log.Fatal("Items deleted are not equal to expected")
	}
}

func testDeleteSystemsHostnameDestinationPortGroupAndSearchsOnlySearchs(t *testing.T, options Options, createdElements []Item) {
	options.Action = "delete"
	options.DeleteAllSystems = false
	app := &App{}
	deletedItems, _, err := app.PapertrailActions(&options)
	if err != nil {
		log.Fatal(err)
	}
	expectedDeletedSearch := NewItem(createdElements[3].ID, createdElements[3].ItemType,
		createdElements[3].ItemName, false, true)
	itemsDeletedExpected := []Item{
		*expectedDeletedSearch,
	}
	if !(EqualSlices(deletedItems, itemsDeletedExpected)) {
		log.Fatal("Items deleted are not equal to expected")
	}
}

func TestApp_PapertrailActionsCreateSystemsHostnameDestinationPortGroupAndSearchs(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	options := &Options{
		"group-test",
		"15.21.10.1, 3.2.13.90",
		destinationDefaultPort,
		0,
		"",
		"hostname",
		"default search test",
		"*",
		"c",
		false,
		false,
		true,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	}
	createdItems, _, err := app.PapertrailActions(options)
	defer testDeleteSystemsHostnameDestinationPortGroupAndAllSearchs(t, *options, createdItems)
	defer testDeleteSystemsHostnameDestinationPortGroupAndSearchsOnlySearchs(t, *options, createdItems)
	expectedCreatedSystem1 := NewItem(0, "System", "15.21.10.1", true, false)
	expectedCreatedSystem2 := NewItem(0, "System", "3.2.13.90", true, false)
	expectedCreatedGroup := NewItem(0, "Group", "group-test", true, false)
	expectedCreatedSearch := NewItem(0, "Search", "default search test", true, false)
	if !(itemEqualWithoutId(createdItems[0], *expectedCreatedSystem1) &&
		itemEqualWithoutId(createdItems[1], *expectedCreatedSystem2) &&
		itemEqualWithoutId(createdItems[2], *expectedCreatedGroup) &&
		itemEqualWithoutId(createdItems[3], *expectedCreatedSearch)) {
		log.Fatal("Items created are not equal to expected (without comparing id)")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func TestApp_PapertrailActionsCreateSystemsHostnameDestinationIdGroupAndSearchs(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	options := &Options{
		"group-test",
		"15.21.10.1, 3.2.13.90",
		0,
		destinationDefaultId,
		"",
		"hostname",
		"default search test",
		"*",
		"c",
		true,
		false,
		true,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	}
	createdItems, _, err := app.PapertrailActions(options)
	defer testDeleteSystemsHostnameDestinationPortGroupAndAllSearchs(t, *options, createdItems)
	expectedCreatedSystem1 := NewItem(0, "System", "15.21.10.1", true, false)
	expectedCreatedSystem2 := NewItem(0, "System", "3.2.13.90", true, false)
	expectedCreatedGroup := NewItem(0, "Group", "group-test", true, false)
	expectedCreatedSearch := NewItem(0, "Search", "default search test", true, false)
	if !(itemEqualWithoutId(createdItems[0], *expectedCreatedSystem1) &&
		itemEqualWithoutId(createdItems[1], *expectedCreatedSystem2) &&
		itemEqualWithoutId(createdItems[2], *expectedCreatedGroup) &&
		itemEqualWithoutId(createdItems[3], *expectedCreatedSearch)) {
		log.Fatal("Items created are not equal to expected (without comparing id)")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func testDeleteOnlySystemIpAddressDestinationPort(t *testing.T, options Options,
	createdElements []Item) {
	options.Action = "delete"
	options.DeleteAllSystems = true
	app := &App{}
	deletedItems, _, err := app.PapertrailActions(&options)
	if err != nil {
		log.Fatal(err)
	}
	expectedDeletedSystem1 := NewItem(createdElements[0].ID, createdElements[0].ItemType,
		createdElements[0].ItemName, false, true)
	itemsDeletedExpected := []Item{
		*expectedDeletedSystem1,
	}
	if !(EqualSlices(deletedItems, itemsDeletedExpected)) {
		log.Fatal("Items deleted are not equal to expected")
	}
}

func TestApp_PapertrailActionsCreateSystemIpAddressDestinationIdGroupAndSearchsWithoutSystems(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	options := &Options{
		"group-test",
		"15.21.10.1",
		0,
		0,
		"15.21.10.1",
		"ip-address",
		"default search test",
		"*",
		"c",
		false,
		false,
		false,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	}
	createdItems, _, err := app.PapertrailActions(options)
	defer testDeleteOnlySystemIpAddressDestinationPort(t, *options, createdItems)
	defer testDeleteGroupAndSearchsWithSystemBasedInDestinationPort(t, *options, createdItems)
	defer testDeleteOnlySearchsWithSystemBasedInDestinationPort(t, *options, createdItems)
	expectedCreatedSystem1 := NewItem(0, "System", "15.21.10.1", true, false)
	expectedCreatedGroup := NewItem(0, "Group", "group-test", true, false)
	expectedCreatedSearch := NewItem(0, "Search", "default search test", true, false)
	if !(itemEqualWithoutId(createdItems[0], *expectedCreatedSystem1) &&
		itemEqualWithoutId(createdItems[1], *expectedCreatedGroup) &&
		itemEqualWithoutId(createdItems[2], *expectedCreatedSearch)) {
		log.Fatal("Items created are not equal to expected (without comparing id)")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func testDeleteSystemIpAddressDestinationPortGroupSearchsAndSystems(t *testing.T, options Options,
	createdElements []Item) {
	options.Action = "delete"
	options.DeleteAllSearches = true
	options.DeleteAllSystems = true
	app := &App{}
	deletedItems, _, err := app.PapertrailActions(&options)
	if err != nil {
		log.Fatal(err)
	}
	expectedDeletedSystem := NewItem(createdElements[0].ID, createdElements[0].ItemType,
		createdElements[0].ItemName, false, true)
	expectedDeletedGroup := NewItem(createdElements[1].ID, createdElements[1].ItemType,
		createdElements[1].ItemName, false, true)
	itemsDeletedExpected := []Item{
		*expectedDeletedSystem,
		*expectedDeletedGroup,
	}
	if !(EqualSlices(deletedItems, itemsDeletedExpected)) {
		log.Fatal("Items deleted are not equal to expected")
	}
}

func TestApp_PapertrailActionsCreateSystemIpAddressDestinationIdGroupAndSearchs(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	options := &Options{
		"group-test",
		"15.21.10.1",
		0,
		0,
		"15.21.10.1",
		"ip-address",
		"default search test",
		"*",
		"c",
		true,
		false,
		true,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	}
	createdItems, _, err := app.PapertrailActions(options)
	defer testDeleteSystemIpAddressDestinationPortGroupSearchsAndSystems(t, *options, createdItems)
	expectedCreatedSystem1 := NewItem(0, "System", "15.21.10.1", true, false)
	expectedCreatedGroup := NewItem(0, "Group", "group-test", true, false)
	expectedCreatedSearch := NewItem(0, "Search", "default search test", true, false)
	if !(itemEqualWithoutId(createdItems[0], *expectedCreatedSystem1) &&
		itemEqualWithoutId(createdItems[1], *expectedCreatedGroup) &&
		itemEqualWithoutId(createdItems[2], *expectedCreatedSearch)) {
		log.Fatal("Items created are not equal to expected (without comparing id)")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func TestApp_PapertrailActionsCreateSystemIpAddressInvalid(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	options := &Options{
		"group-test",
		"10.1.2.11",
		0,
		0,
		"192.168.0.1",
		"ip-address",
		"default search test",
		"*",
		"c",
		false,
		false,
		true,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	}
	_, _, err := app.PapertrailActions(options)
	expectedError := convertStatusCodeToError(400, "System", "Creating")
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestApp_PapertrailActionsInvalidadDestinationId(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	options := &Options{
		"group-test",
		"15.21.10.1, 3.2.13.90",
		0,
		177547777692,
		"",
		"hostname",
		"default search test",
		"*",
		"c",
		false,
		false,
		true,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	}
	_, _, err := app.PapertrailActions(options)
	expectedError := errors.New("Error: Destination not found ")
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func testDeleteSystemsHostnameDestinationPortGroupAndSearchsDeleteAll(t *testing.T,
	options Options, createdElements []Item) {
	options.Action = "delete"
	options.DeleteAllSearches = true
	app := &App{}
	deletedItems, _, err := app.PapertrailActions(&options)
	if err != nil {
		log.Fatal(err)
	}
	expectedDeletedSystem := NewItem(createdElements[0].ID, createdElements[0].ItemType,
		createdElements[0].ItemName, false, true)
	expectedDeletedGroup := NewItem(createdElements[1].ID, createdElements[1].ItemType,
		createdElements[1].ItemName, false, true)
	itemsDeletedExpected := []Item{
		*expectedDeletedSystem,
		*expectedDeletedGroup,
	}
	if !(EqualSlices(deletedItems, itemsDeletedExpected)) {
		log.Fatal("Items deleted are not equal to expected")
	}
}

func TestApp_PapertrailActionsCreateRepeatedSystemsHostnameDestinationPortGroupAndSearchs(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	options := &Options{
		"group-test",
		"10.1.2.11, 10.1.2.11",
		destinationDefaultPort,
		0,
		"",
		"hostname",
		"default search test",
		"*",
		"c",
		false,
		false,
		true,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	}
	createdItems, _, err := app.PapertrailActions(options)
	defer testDeleteSystemsHostnameDestinationPortGroupAndSearchsDeleteAll(t, *options, createdItems)
	expectedCreatedSystem := NewItem(0, "System", "10.1.2.11", true, false)
	expectedCreatedGroup := NewItem(0, "Group", "group-test", true, false)
	expectedCreatedSearch := NewItem(0, "Search", "default search test", true, false)
	if !(itemEqualWithoutId(createdItems[0], *expectedCreatedSystem) &&
		itemEqualWithoutId(createdItems[1], *expectedCreatedGroup) &&
		itemEqualWithoutId(createdItems[2], *expectedCreatedSearch)) {
		log.Fatal("Items created are not equal to expected (without comparing id)")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func testDeleteOnlySearchsWithSystemBasedInDestinationPort(t *testing.T, options Options, createdElements []Item) {
	options.DeleteAllSystems = false
	options.Action = "delete"
	app := &App{}
	deletedItems, _, err := app.PapertrailActions(&options)
	if err != nil {
		log.Fatal(err)
	}
	expectedDeletedSearch := NewItem(createdElements[2].ID, createdElements[2].ItemType,
		createdElements[2].ItemName, false, true)
	itemsDeletedExpected := []Item{
		*expectedDeletedSearch,
	}
	if !(EqualSlices(deletedItems, itemsDeletedExpected)) {
		log.Fatal("Items deleted are not equal to expected")
	}
}

func testDeleteGroupAndSearchsWithSystemBasedInDestinationPort(t *testing.T, options Options, createdElements []Item) {
	options.DeleteAllSystems = false
	options.DeleteAllSearches = true
	options.Action = "delete"
	app := &App{}
	deletedItems, _, err := app.PapertrailActions(&options)
	if err != nil {
		log.Fatal(err)
	}
	expectedDeletedSearch := NewItem(createdElements[1].ID, createdElements[1].ItemType,
		createdElements[1].ItemName, false, true)
	itemsDeletedExpected := []Item{
		*expectedDeletedSearch,
	}
	if !(EqualSlices(deletedItems, itemsDeletedExpected)) {
		log.Fatal("Items deleted are not equal to expected")
	}
}

func testDeleteGroupAndSearchsWithoutSystems(t *testing.T, options Options, createdElements []Item) {
	options.DeleteAllSystems = false
	options.DeleteAllSearches = true
	options.Action = "delete"
	app := &App{}
	deletedItems, _, err := app.PapertrailActions(&options)
	if err != nil {
		log.Fatal(err)
	}
	expectedDeletedGroup := NewItem(createdElements[2].ID, createdElements[2].ItemType,
		createdElements[2].ItemName, false, true)
	itemsDeletedExpected := []Item{
		*expectedDeletedGroup,
	}
	if !(EqualSlices(deletedItems, itemsDeletedExpected)) {
		log.Fatal("Items deleted are not equal to expected")
	}
}

func TestApp_PapertrailActionsDeleteInvalidGroup(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	options := &Options{
		"group-test",
		"15.21.10.1, 3.2.13.90",
		destinationDefaultPort,
		0,
		"",
		"hostname",
		"default search test",
		"*",
		"c",
		false,
		false,
		false,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	}
	createdItems, _, err := app.PapertrailActions(options)
	expectedCreatedSystem1 := NewItem(0, "System", "15.21.10.1", true, false)
	expectedCreatedSystem2 := NewItem(0, "System", "3.2.13.90", true, false)
	expectedCreatedGroup := NewItem(0, "Group", "group-test", true, false)
	expectedCreatedSearch := NewItem(0, "Search", "default search test", true, false)
	if !(itemEqualWithoutId(createdItems[0], *expectedCreatedSystem1) &&
		itemEqualWithoutId(createdItems[1], *expectedCreatedSystem2) &&
		itemEqualWithoutId(createdItems[2], *expectedCreatedGroup) &&
		itemEqualWithoutId(createdItems[3], *expectedCreatedSearch)) {
		log.Fatal("Items created are not equal to expected (without comparing id)")
	}
	if err != nil {
		log.Fatal(err)
	}
	defer testDeleteOnlySystemsHostnameDestinationPort(t, *options, createdItems)
	defer testDeleteGroupAndSearchsWithoutSystems(t, *options, createdItems)
	options.DeleteAllSystems = false
	options.Action = "delete"
	options.GroupName = "group-test invalid"
	deletedItems, _, err := app.PapertrailActions(options)
	errExpected := errors.New("Error: Group with name " + options.GroupName + " doesn't exist ")
	if err.Error() != errExpected.Error() {
		t.Fatal("The error obtained is not the expected")
	}
	if len(deletedItems) > 0 {
		t.Fatal("The list of deleted items should not be greater than 0")
	}
}

func TestApp_PapertrailActionsDeleteOnlySystems(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	options := &Options{
		"group-test",
		"15.21.10.1, 3.2.13.90",
		destinationDefaultPort,
		0,
		"",
		"hostname",
		"default search test",
		"*",
		"c",
		false,
		false,
		true,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	}
	createdItems, _, err := app.PapertrailActions(options)
	expectedCreatedSystem1 := NewItem(0, "System", "15.21.10.1", true, false)
	expectedCreatedSystem2 := NewItem(0, "System", "3.2.13.90", true, false)
	expectedCreatedGroup := NewItem(0, "Group", "group-test", true, false)
	expectedCreatedSearch := NewItem(0, "Search", "default search test", true, false)
	if !(itemEqualWithoutId(createdItems[0], *expectedCreatedSystem1) &&
		itemEqualWithoutId(createdItems[1], *expectedCreatedSystem2) &&
		itemEqualWithoutId(createdItems[2], *expectedCreatedGroup) &&
		itemEqualWithoutId(createdItems[3], *expectedCreatedSearch)) {
		log.Fatal("Items created are not equal to expected (without comparing id)")
	}
	if err != nil {
		log.Fatal(err)
	}
	defer testDeleteGroupAndSearchsWithoutSystems(t, *options, createdItems)
	options.Action = "delete"
	options.DeleteOnlySystems = true
	deletedItems, _, _ := app.PapertrailActions(options)
	expectedDeletedSystem1 := NewItem(createdItems[0].ID, createdItems[0].ItemType,
		createdItems[0].ItemName, false, true)
	expectedDeletedSystem2 := NewItem(createdItems[1].ID, createdItems[1].ItemType,
		createdItems[1].ItemName, false, true)
	itemsDeletedExpected := []Item{
		*expectedDeletedSystem1,
		*expectedDeletedSystem2,
	}
	if !(EqualSlices(deletedItems, itemsDeletedExpected)) {
		log.Fatal("Items deleted are not equal to expected")
	}
}

func TestApp_PapertrailActionsDeleteInvalidSearch(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	options := &Options{
		"group-test",
		"15.21.10.1, 3.2.13.90",
		destinationDefaultPort,
		0,
		"",
		"hostname",
		"default search test",
		"*",
		"c",
		false,
		false,
		true,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	}
	createdItems, _, err := app.PapertrailActions(options)
	expectedCreatedSystem1 := NewItem(0, "System", "15.21.10.1", true, false)
	expectedCreatedSystem2 := NewItem(0, "System", "3.2.13.90", true, false)
	expectedCreatedGroup := NewItem(0, "Group", "group-test", true, false)
	expectedCreatedSearch := NewItem(0, "Search", "default search test", true, false)
	if !(itemEqualWithoutId(createdItems[0], *expectedCreatedSystem1) &&
		itemEqualWithoutId(createdItems[1], *expectedCreatedSystem2) &&
		itemEqualWithoutId(createdItems[2], *expectedCreatedGroup) &&
		itemEqualWithoutId(createdItems[3], *expectedCreatedSearch)) {
		log.Fatal("Items created are not equal to expected (without comparing id)")
	}
	if err != nil {
		log.Fatal(err)
	}
	defer testDeleteGroupAndSearchsWithoutSystems(t, *options, createdItems)
	options.Action = "delete"
	options.Search = "default search invalid"
	deletedItems, _, err := app.PapertrailActions(options)
	errExpected := errors.New("Error: Search with name " + options.Search + " doesn't exist")
	if err.Error() != errExpected.Error() {
		t.Fatal("The error obtained is not the expected")
	}
	expectedDeletedSystem1 := NewItem(createdItems[0].ID, createdItems[0].ItemType,
		createdItems[0].ItemName, false, true)
	expectedDeletedSystem2 := NewItem(createdItems[1].ID, createdItems[1].ItemType,
		createdItems[1].ItemName, false, true)
	itemsDeletedExpected := []Item{
		*expectedDeletedSystem1,
		*expectedDeletedSystem2,
	}
	if !(EqualSlices(deletedItems, itemsDeletedExpected)) {
		log.Fatal("Items deleted are not equal to expected")
	}
}

func EqualItems(a, b []Item) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestApp_PapertrailActionsObtainLogsSystemsHostnameDestinationPortGroupAndSearchs(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	options := &Options{
		"group-test",
		"15.21.10.1, 3.2.13.90",
		destinationDefaultPort,
		0,
		"",
		"hostname",
		"default search test",
		"*",
		"c",
		false,
		false,
		true,
		false,
		nowDateLessEightHours,
		nowDate,
		"/tmp/",
	}
	createdItems, _, err := app.PapertrailActions(options)
	defer testDeleteSystemsHostnameDestinationPortGroupAndAllSearchs(t, *options, createdItems)
	expectedCreatedSystem1 := NewItem(0, "System", "15.21.10.1", true, false)
	expectedCreatedSystem2 := NewItem(0, "System", "3.2.13.90", true, false)
	expectedCreatedGroup := NewItem(0, "Group", "group-test", true, false)
	expectedCreatedSearch := NewItem(0, "Search", "default search test", true, false)
	if !(itemEqualWithoutId(createdItems[0], *expectedCreatedSystem1) &&
		itemEqualWithoutId(createdItems[1], *expectedCreatedSystem2) &&
		itemEqualWithoutId(createdItems[2], *expectedCreatedGroup) &&
		itemEqualWithoutId(createdItems[3], *expectedCreatedSearch)) {
		log.Fatal("Items created are not equal to expected (without comparing id)")
	}
	if err != nil {
		log.Fatal(err)
	}
	options.Action = "obtain"
	obtainedItems, _, err := app.PapertrailActions(options)
	unixStartDate, _ := GetTimeStampUnixFromDate(options.StartDate)
	unixEndDate, _ := GetTimeStampUnixFromDate(options.EndDate)
	itemExpectedName := CreateFilenameForEventsSearch(options.Path, options.GroupName, options.Search, unixStartDate, unixEndDate) + " with 0 events retrieved"
	obtainedItemExpected := Item{
		ID:       0,
		ItemType: "EventsSearch",
		ItemName: itemExpectedName,
		Created:  false,
		Deleted:  false,
	}
	obtainedItemsExpected := []Item{obtainedItemExpected}
	if err != nil {
		log.Fatal(err)
	}
	if !EqualItems(obtainedItemsExpected, obtainedItems) {
		log.Fatal("Items obtained are not equal to expected")
	}
}

func TestApp_PapertrailActionsObtainIncorrectStartDate(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	options := &Options{
		"group-test",
		"15.21.10.1, 3.2.13.90",
		destinationDefaultPort,
		0,
		"",
		"hostname",
		"default search test",
		"*",
		"o",
		false,
		false,
		true,
		false,
		"14/08/2020 10:20:00",
		"04/08/2020 10:40:00",
		"/tmp/",
	}
	_, _, err := app.PapertrailActions(options)
	expectedError := fmt.Errorf("cannot parse startdate: parsing time \"%v\": month out of range", options.StartDate)
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestApp_PapertrailActionsObtainIncorrectEndDate(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	options := &Options{
		"group-test",
		"15.21.10.1, 3.2.13.90",
		destinationDefaultPort,
		0,
		"",
		"hostname",
		"default search test",
		"*",
		"o",
		false,
		false,
		true,
		false,
		"04/08/2020 10:20:00",
		"14/08/2020 10:40:00",
		"/tmp/",
	}
	_, _, err := app.PapertrailActions(options)
	expectedError := fmt.Errorf("cannot parse enddate: parsing time \"%v\": month out of range", options.EndDate)
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestApp_PapertrailActionsObtainIncorrectDates(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &App{}
	options := &Options{
		"group-test",
		"15.21.10.1, 3.2.13.90",
		destinationDefaultPort,
		0,
		"",
		"hostname",
		"default search test",
		"*",
		"o",
		false,
		false,
		true,
		false,
		nowDate,
		nowDateLessEightHours,
		"/tmp/",
	}
	_, _, err := app.PapertrailActions(options)
	expectedError := errors.New("startdate > enddate - please set proper data boundaries")
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}
