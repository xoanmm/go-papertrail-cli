/*
NAME:
   go-papertrail-cli - interacts with papertrail through its api to perform both log collection actions and the creation/deletion of systems, groups and saved searches

USAGE:
   go-papertrail-cli [--group-name <group-name>] [--system-wildcard <wildcard>] [--search <search-name>] [--query <query>] [--action <action>] [--delete-all-searches <delete-all-searches>] [--delete-only-searches <delete-only-searches>] [--delete-all-systems <delete-all-systems>]  [--delete-only-systems <delete-only-systems>][--start-date <start-date>] [--end-date <end-date>] [--path <path>]

VERSION:
   1.3.0

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
   --delete-only-searches              Indicates if only searches specified are going to be deleted (default: false)
   --delete-all-systems, -D            Indicates if all systems specified are going to be deleted (default: true)
   --delete-only-systems               Indicates if only systems specified are going to be deleted (default: false)
   --start-date value, -s value        filter only from a date specified ('mm/dd/yyyy hh:mm:ss' format UTC time) (default: $ACTUAL_DATE - 8hours)
   --end-date value, -e value          filter only until a date specified ('mm/dd/yyyy hh:mm:ss' format UTC time) (default: $ACTUAL_DATE)
   --path value, -P value              path where to store the logs (default: "/tmp")
   --help, -h                          show help (default: false)
   --version, -v                       print the version (default: false)
*/
package main
