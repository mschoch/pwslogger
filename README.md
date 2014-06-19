# pwslogger

A utility which takes a JSON file as input, and builds an HTTP query to update a Weather Underground Personal Weather Station.

## Usage:

This example retrieves sensor JSON from a [nowd](https://github.com/mschoch/nowd) server, pipes the JSON to this utility.

The paths specified are JSON pointers to sensor values.

    curl -s http://192.168.1.11:4793/ | ./pwslogger -id=<mystationid> -passwd=<mypassword> -tempCPath=/wx/t -humidityPath=/wx/h -pressurePath=/basement/p