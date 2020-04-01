package papertrail

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/xoanmm/go-papertrail-cli/pkg/papertrail"
	"log"
	"os"
	"testing"
)

var papertrailApiToken string

func setupEnv() error{
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
	code := m.Run()
	os.Exit(code)
}

func TestApp_PapertrailActionsNoProvidedToken(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	os.Setenv("PAPERTRAIL_API_TOKEN", "")
	app := papertrail.App{}
	_, _, err := app.PapertrailActions(&papertrail.Options{
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
	})
	expectedError := errors.New("Error getting value of PAPERTRAIL_API_TOKEN, it's necessary to define this variable with your papertrail's API token ")
	if err.Error() != expectedError.Error()  {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestApp_PapertrailActionsInvalidHostConfiguration(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &papertrail.App{}
	errT := os.Setenv("PAPERTRAIL_API_TOKEN", "tokenPapertrail")
	if errT != nil {
		fmt.Println("errT is", errT)
	}
	_, _, err := app.PapertrailActions(&papertrail.Options{
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
	})
	expectedError := errors.New("It's necessary provide a value distinct from default (0) to destination id or destination port ")
	if err.Error() != expectedError.Error()  {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestApp_PapertrailActionsInvalidActionConfiguration(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &papertrail.App{}
	errT := os.Setenv("PAPERTRAIL_API_TOKEN", "tokenPapertrail")
	if errT != nil {
		fmt.Println("errT is", errT)
	}
	_, _, err := app.PapertrailActions(&papertrail.Options{
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
	})
	expectedError := errors.New("Not valid option provided for action to perform, the only valid values are: \n" +
		"\t'c' or 'create': create new system/s, group and/or search\n" +
		"\t'd' or 'delete': create new system/s, group and/or search\n" +
		"\t'o'or 'obtain': obtain logs in base of parameters provided\n")
	if err.Error() != expectedError.Error()  {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestApp_PapertrailActionsInvalidSystemProvided(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &papertrail.App{}
	errT := os.Setenv("PAPERTRAIL_API_TOKEN", "tokenPapertrail")
	if errT != nil {
		fmt.Println("errT is", errT)
	}
	_, _, err := app.PapertrailActions(&papertrail.Options{
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
	})
	expectedError := errors.New("Not valid option provided for system, the only valid values are: \n" +
		"\t'h' or 'hostname': system based in hostname\n" +
		"\t'i'or 'ip-address': system based in ip-address\n")
	if err.Error() != expectedError.Error()  {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestApp_PapertrailActionsTwoSystemConfigurationProvided(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &papertrail.App{}
	errT := os.Setenv("PAPERTRAIL_API_TOKEN", "tokenPapertrail")
	if errT != nil {
		fmt.Println("errT is", errT)
	}
	_, _, err := app.PapertrailActions(&papertrail.Options{
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
	})
	expectedError := errors.New("If the system is a hostname-type system, only destination " +
		"id or destination port can be specified\n")
	if err.Error() != expectedError.Error()  {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestApp_PapertrailActionsTwoSystemConfigurationIncorrect(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &papertrail.App{}
	errT := os.Setenv("PAPERTRAIL_API_TOKEN", "tokenPapertrail")
	if errT != nil {
		fmt.Println("errT is", errT)
	}
	_, _, err := app.PapertrailActions(&papertrail.Options{
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
	})
	expectedError := errors.New("It's necessary provide a value distinct from default (0) to " +
		"destination id or destination port ")
	if err.Error() != expectedError.Error()  {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestApp_PapertrailActionsInvalidIpAddressConfiguration(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &papertrail.App{}
	errT := os.Setenv("PAPERTRAIL_API_TOKEN", "tokenPapertrail")
	if errT != nil {
		fmt.Println("errT is", errT)
	}
	_, _, err := app.PapertrailActions(&papertrail.Options{
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
	})
	expectedError := errors.New("The IP Address provided, 11111111 it's not a valid IP Address ")
	if err.Error() != expectedError.Error()  {
		t.Fatal("The error obtained is not the expected")
	}
}

func itemEqualWithoutId(expected papertrail.Item, obtained papertrail.Item) bool {
	if expected.Created == obtained.Created && expected.Deleted == obtained.Deleted &&
		expected.ItemName == obtained.ItemName && expected.ItemType == obtained.ItemType {
		return true
	}
	return false
}

func EqualSlices(a, b []papertrail.Item) bool {
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

func testDeleteSystemsHostnameDestinationPortGroupAndSearchs(t *testing.T, options papertrail.Options, createdElements []papertrail.Item) {
	options.Action = "delete"
	app := &papertrail.App{}
	deletedItems, _, err := app.PapertrailActions(&options)
	if err != nil {
		log.Fatal(err)
	}
	expectedDeletedSystem1 := papertrail.NewItem(createdElements[0].ID, createdElements[0].ItemType,
		createdElements[0].ItemName, false, true)
	expectedDeletedSystem2 := papertrail.NewItem(createdElements[1].ID, createdElements[1].ItemType,
		createdElements[1].ItemName, false, true)
	expectedDeletedSearch := papertrail.NewItem(createdElements[3].ID, createdElements[3].ItemType,
		createdElements[3].ItemName, false, true)
	expectedDeletedGroup := papertrail.NewItem(createdElements[2].ID, createdElements[2].ItemType,
		createdElements[2].ItemName, false, true)
	itemsDeletedExpected := []papertrail.Item{
		*expectedDeletedSystem1,
		*expectedDeletedSystem2,
		*expectedDeletedSearch,
		*expectedDeletedGroup,
	}
	if !(EqualSlices(deletedItems, itemsDeletedExpected)){
		log.Fatal("Items deleted are not equal to expected")
	}
}

func TestApp_PapertrailActionsCreateSystemsHostnameDestinationPortGroupAndSearchs(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &papertrail.App{}
	options := &papertrail.Options{
		"group-test",
		"18.211.9.147, 3.222.14.80",
		23633,
		0,
		"",
		"hostname",
		"default search test",
		"*",
		"c",
		false,
	}
	createdItems, _, err := app.PapertrailActions(options)
	defer testDeleteSystemsHostnameDestinationPortGroupAndSearchs(t, *options, createdItems)
	expectedCreatedSystem1 := papertrail.NewItem(0, "System", "18.211.9.147", true, false)
	expectedCreatedSystem2 := papertrail.NewItem(0, "System", "3.222.14.80", true, false)
	expectedCreatedGroup := papertrail.NewItem(0, "Group", "group-test", true, false)
	expectedCreatedSearch := papertrail.NewItem(0, "Search", "default search test", true, false)
	if !(itemEqualWithoutId(createdItems[0], *expectedCreatedSystem1) &&
		itemEqualWithoutId(createdItems[1], *expectedCreatedSystem2) &&
		itemEqualWithoutId(createdItems[2], *expectedCreatedGroup) &&
		itemEqualWithoutId(createdItems[3], *expectedCreatedSearch)){
		log.Fatal("Items created are not equal to expected (without comparing id)")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func TestApp_PapertrailActionsCreateSystemsHostnameDestinationIdGroupAndSearchs(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &papertrail.App{}
	options := &papertrail.Options{
		"group-test",
		"18.211.9.147, 3.222.14.80",
		0,
		17754692,
		"",
		"hostname",
		"default search test",
		"*",
		"c",
		false,
	}
	createdItems, _, err := app.PapertrailActions(options)
	defer testDeleteSystemsHostnameDestinationPortGroupAndSearchs(t, *options, createdItems)
	expectedCreatedSystem1 := papertrail.NewItem(0, "System", "18.211.9.147", true, false)
	expectedCreatedSystem2 := papertrail.NewItem(0, "System", "3.222.14.80", true, false)
	expectedCreatedGroup := papertrail.NewItem(0, "Group", "group-test", true, false)
	expectedCreatedSearch := papertrail.NewItem(0, "Search", "default search test", true, false)
	if !(itemEqualWithoutId(createdItems[0], *expectedCreatedSystem1) &&
		itemEqualWithoutId(createdItems[1], *expectedCreatedSystem2) &&
		itemEqualWithoutId(createdItems[2], *expectedCreatedGroup) &&
		itemEqualWithoutId(createdItems[3], *expectedCreatedSearch)){
		log.Fatal("Items created are not equal to expected (without comparing id)")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func testDeleteSystemIpAddressDestinationPortGroupAndSearchs(t *testing.T, options papertrail.Options,
	createdElements []papertrail.Item) {
	options.Action = "delete"
	app := &papertrail.App{}
	deletedItems, _, err := app.PapertrailActions(&options)
	if err != nil {
		log.Fatal(err)
	}
	expectedDeletedSystem := papertrail.NewItem(createdElements[0].ID, createdElements[0].ItemType,
		createdElements[0].ItemName, false, true)
	expectedDeletedSearch := papertrail.NewItem(createdElements[2].ID, createdElements[2].ItemType,
		createdElements[2].ItemName, false, true)
	expectedDeletedGroup := papertrail.NewItem(createdElements[1].ID, createdElements[1].ItemType,
		createdElements[1].ItemName, false, true)
	itemsDeletedExpected := []papertrail.Item{
		*expectedDeletedSystem,
		*expectedDeletedSearch,
		*expectedDeletedGroup,
	}
	if !(EqualSlices(deletedItems, itemsDeletedExpected)){
		log.Fatal("Items deleted are not equal to expected")
	}
}

func TestApp_PapertrailActionsCreateSystemIpAddressDestinationIdGroupAndSearchs(t *testing.T) {
	defer os.Setenv("PAPERTRAIL_API_TOKEN", papertrailApiToken)
	app := &papertrail.App{}
	options := &papertrail.Options{
		"group-test",
		"18.211.9.147",
		0,
		17754692,
		"18.211.9.147",
		"ip-address",
		"default search test",
		"*",
		"c",
		false,
	}
	createdItems, _, err := app.PapertrailActions(options)
	defer testDeleteSystemIpAddressDestinationPortGroupAndSearchs(t, *options, createdItems)
	expectedCreatedSystem1 := papertrail.NewItem(0, "System", "18.211.9.147", true, false)
	expectedCreatedGroup := papertrail.NewItem(0, "Group", "group-test", true, false)
	expectedCreatedSearch := papertrail.NewItem(0, "Search", "default search test", true, false)
	if !(itemEqualWithoutId(createdItems[0], *expectedCreatedSystem1) &&
		itemEqualWithoutId(createdItems[1], *expectedCreatedGroup) &&
		itemEqualWithoutId(createdItems[2], *expectedCreatedSearch)){
		log.Fatal("Items created are not equal to expected (without comparing id)")
	}
	if err != nil {
		log.Fatal(err)
	}
}