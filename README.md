[![GitHub Release](https://img.shields.io/github/release/xoanmm/go-papertrail-cli.svg?logo=github&labelColor=262b30)](https://github.com/xoanmm/go-papertrail-cli/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/xoanmm/go-papertrail-cli)](https://goreportcard.com/report/github.com/xoanmm/go-papertrail-cli)
[![License](https://img.shields.io/github/license/xoanmm/go-papertrail-cli)](https://github.com/xoanmm/go-papertrail-cli/LICENSE)

# go-papertrail-cli
A simple tool to interacts with [Papertrail](https://papertrailapp.com/) through its [api](https://help.papertrailapp.com/kb/how-it-works/http-api/) to perform both log collection actions and the creation of [systems](https://help.papertrailapp.com/kb/how-it-works/adding-and-removing-senders/), [groups](https://help.papertrailapp.com/kb/how-it-works/groups/) and [searches](https://help.papertrailapp.com/kb/how-it-works/search-syntax).

The tool has been created with the intention of facilitating the creation and/or elimination of the elements mentioned in papertrail, as well as obtaining logs of a given search and storing them in an orderly manner in a file.

## Installation

Go to [release page](https://github.com/xoanmm/go-papertrail-cli/releases) and download the binary you need.

## Example of usage

To use this tool it is necessary to export the variable related to the papertrail token as follows
```
export PAPERTRAIL_API_TOKEN=<user_api_token_papertrail>
```

Examples of implementation for the different actions available are given below:
- Creation:

  - Example of the creation of two systems that will send logs, as well as a group and a search for the registration of the same:
      
      ```bash
      $ ./go-papertrail-cli -a c -g "group-test" -w "15.21.10.1, 3.2.13.90" -S "default search test" -q "*" -p 23633 -t "hostname"
      2020/04/13 19:09:24 Checking conditions for do action 'c' in papertrail params: [group-name group-test] [system-wildcard 15.21.10.1, 3.2.13.90] [search default search test] [query *] [delete-all-searches false] [--start-date 04/13/2020 09:09:24] [--end-date 04/13/2020 17:09:24] [--path /tmp]
      2020/04/13 19:09:25 System with hostname 15.21.10.1 doesn't exist yet
      2020/04/13 19:09:26 System with name 15.21.10.1 based in hostname 15.21.10.1 was successfully created with id 5401098672
      2020/04/13 19:09:26 System with hostname 3.2.13.90 doesn't exist yet
      2020/04/13 19:09:27 System with name 3.2.13.90 based in hostname 3.2.13.90 was successfully created with id 5401098772
      2020/04/13 19:09:28 Group with name group-test and id 19443892 was successfully created
      2020/04/13 19:09:29 Search with name default search test and id 84545002 was successfully created
      2020/04/13 19:09:29 Create actions have been carried out on the following elements
      2020/04/13 19:09:29 - System with ID 5401098672 and name '15.21.10.1'
      2020/04/13 19:09:29 - System with ID 5401098772 and name '3.2.13.90'
      2020/04/13 19:09:29 - Group with ID 19443892 and name 'group-test'
      2020/04/13 19:09:29 - Search with ID 84545002 and name 'default search test'
      ```
    
- Deletion:

  - Example of deleting the resources created in the execution of the previous example on papertrail, deleting only the specified search (`false` value in *delete-all-searches* parameter)
    
       ```bash
       $ ./go-papertrail-cli -a d -g "group-test" -w "15.21.10.1, 3.2.13.90" -S "default search test" -q "*" -p 23633 -t "hostname" -d true
       2020/04/13 19:11:03 Checking conditions for do action 'd' in papertrail params: [group-name group-test] [system-wildcard 15.21.10.1, 3.2.13.90] [search default search test] [query *] [delete-all-searches false] [--start-date 04/13/2020 09:11:03] [--end-date 04/13/2020 17:11:03] [--path /tmp]
       2020/04/13 19:11:04 System with hostname 15.21.10.1 exists with id 5401098672
       2020/04/13 19:11:05 System with id 5401098672 was successfully deleted
       2020/04/13 19:11:05 System with hostname 3.2.13.90 exists with id 5401098772
       2020/04/13 19:11:06 System with id 5401098772 was successfully deleted
       2020/04/13 19:11:06 Group with name group-test exists with id 19443892
       2020/04/13 19:11:06 Search with name default search test exists with id 84545002
       2020/04/13 19:11:07 Search with name default search test and id 84545002 was successfully deleted
       2020/04/13 19:11:07 Group with name group-test exists with id 19443892
       2020/04/13 19:11:08 Group with name group-test and id 19443892 was successfully deleted
       2020/04/13 19:11:08 Delete actions have been carried out on the following elements
       2020/04/13 19:11:08 - System with ID 5401098672 and name '15.21.10.1'
       2020/04/13 19:11:08 - System with ID 5401098772 and name '3.2.13.90'
       2020/04/13 19:11:08 - Search with ID 84545002 and name 'default search test'
       2020/04/13 19:11:08 - Group with ID 19443892 and name 'group-test'
       ```
    
   - Example of deleting the resources created in the execution of the previous example on papertrail, deleting only the specified search (`true` value in *delete-all-searches* parameter)
   
       ```$ ./go-papertrail-cli -a d -g "group-test" -w "15.21.10.1, 3.2.13.90" -S "default search test" -q "*" -p 23633 -t "hostname" -d true
       2020/04/13 19:26:50 Checking conditions for do action 'd' in papertrail params: [group-name group-test] [system-wildcard 15.21.10.1, 3.2.13.90] [delete-all-searches true] [--start-date 04/13/2020 09:26:50] [--end-date 04/13/2020 17:26:50] [--path /tmp]
       2020/04/13 19:26:51 System with hostname 15.21.10.1 exists with id 5401206572
       2020/04/13 19:26:51 System with id 5401206572 was successfully deleted
       2020/04/13 19:26:51 System with hostname 3.2.13.90 exists with id 5401206682
       2020/04/13 19:26:52 System with id 5401206682 was successfully deleted
       2020/04/13 19:26:52 Group with name group-test exists with id 19444202
       2020/04/13 19:26:53 Group with name group-test and id 19444202 was successfully deleted
       2020/04/13 19:26:53 Delete actions have been carried out on the following elements
       2020/04/13 19:26:53 - System with ID 5401206572 and name '15.21.10.1'
       2020/04/13 19:26:53 - System with ID 5401206682 and name '3.2.13.90'
       2020/04/13 19:26:53 - Group with ID 19444202 and name 'group-test'
       ```
     
- Obtain:
  
  - Example of obtaining the logs of a search registered in a group with registered sending systems in a certain time period:
  
      ```bash
      ./go-papertrail-cli -a o -g "group-test" -w "15.21.10.1, 3.2.13.90" -S "default search test" -q "*" -p 23633 -t "hostname"
      2020/04/13 19:31:56 Checking conditions for do action 'o' in papertrail params: [group-name group-test] [system-wildcard 15.21.10.1, 3.2.13.90] [search default search test] [query *] [delete-all-searches false] [--start-date 04/13/2020 09:09:24] [--end-date 04/13/2020 17:09:24] [--path /tmp]
      2020/04/13 19:31:58 System with hostname 15.21.10.1 exists with id 5401206572
      2020/04/13 19:31:58 System with hostname 3.2.13.90 exists with id 5401206682
      2020/04/13 19:31:59 Group with name group-test exists with id 19444202
      2020/04/13 19:32:06 Search with name default search test exists with id 84545002
      2020/04/13 19:32:06 EventsSearch saved in file /tmp/group_test_default_search_04-13-2020_11:20:00_04-13-2020_11:23:00 with 885 events retrieved
      ```
    
## Usage

    NAME:
       go-papertrail-cli - interacts with papertrail through its api to perform both log collection actions and the creation of systems, groups and saved searches
    
    USAGE:
       go-papertrail-cli [--group-name <group-name>] [--system-wildcard <wildcard>] [--search <search-name>] [--query <query>] [--action <action>] [--delete-all-searches <delete-all-searches>] [--startDate <start-date>] [--end-date <end-date>] [--path <path>]
    
    VERSION:
       1.0.0
    
    AUTHOR:
       Xoan Mallon <xoanmallon@gmail.com>
    
    COMMANDS:
       help, h  Shows a list of commands or help for one command
    
    GLOBAL OPTIONS:
       --group-name value, -g value        group defined or to be defined in papertrail (default: "my-log-group")
       --system-wildcard value, -w value   wildcard to be applied on the systems defined in papertrail (default: "*")
       --destination-port value, -p value  destination port for sending the logs of the indicated system/s (default: 0)
       --destination-id value, -I value    destination id for sending the logs of the indicated system/s (default: 0)
       --ip-address value, -i value        source ip address from sending the logs of the indicated system/s
       --system-type value, -t value       Type of system, can be hostname or ip-address (default: "hostname")
       --search value, -S value            name of saved search to be performed on logs or to be created on a group (default: "default search")
       --query value, -q value             query to be performed on the group of logs or applied on the search to be created (default: "*")
       --action value, -a value            Action to be performed with the information provided for papertrail, possible values only c(create), o(obtain) or d(delete) (default: "c")
       --delete-all-searches, -d           Indicates if all searches in a group or a specific search are going to be deleted (default: false)
       --start-date value, -s value        filter only from a date specified ('mm/dd/yyyy hh:mm:ss' format UTC time) (default: $ACTUAL_DATE - 8hours)
       --end-date value, -e value          filter only until a date specified ('mm/dd/yyyy hh:mm:ss' format UTC time) (default: $ACTUAL_DATE)
       --path value, -P value              path where to store the logs (default: "/tmp")
       --help, -h                          show help (default: false)
       --version, -v                       print the version (default: false)
       
### Running the tests

Due to being an application with a single entry point, it does not make sense to perform unit tests, but rather [integration tests](./pkg/papertrail/app_test.go) that check that the expected actions are performed based on the input parameters provided.

#### Tests requirements

A series of variables must be provided in order to carry out the execution of the integration tests mentioned, this variable must be stored in a `.env` file within the `pkg/papertrail` folder, a [template](./pkg/papertrail/.template.env) of the variables that this file must follow is available.

### Dependencies & Refs

  - [urfave/cli](https://github.com/urfave/cli)
  - [joho/godotenv](github.com/joho/godotenv)
  
### LICENSE

 [MIT license](LICENSE)

### Author(s)

- [xoanmm](https://github.com/xoanmm)

### Collaborator(s)

- [boliri](https://github.com/boliri)