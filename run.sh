#!/bin/bash

go build -o bookings-app cmd/web/*.go && ./bookings-app -dbname=bookings_app -dbuser=postgres -dbpass=postgres -cache=false -production=false