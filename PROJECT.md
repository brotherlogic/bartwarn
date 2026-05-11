# Bart Warn

The bart warner project is a client / server system that:

a. Continually tracks my location against geographic hotspots centered on all bart stations
b. A server system that receives notififications that I've entered one of these stations
c. A process which then determines which train line I should get to get home
d. A process which then sends me an SMS with these details

## Location Tracker

This is a simple android background app - I only need to run this myself, I don't need it to be available
on the App Store. It has hardcoded geogprahical triggers centered around bart station and sends a server
message every time I enter one of these locations.

the URL used for pings is https://bartwarn.brotherlogic-backend.com/

## Server

When the server receives a ping from the location tracker it looks up a real time route from the given station
to my local station (El Cerrito Plaza), and let's me know the best train to get home quickly with the following
constraint:

1. I prefer to take a red line train - if the optimal route only saves 5 minutes over waiting for a red
   line train, I'll wait for the red line train
1. Otherwise I'll take the fastest route suggested to get home.

## SMS integration

The server then sends an SMS message to my phone saying "Take XXX line XX:XX" with details about the colour
line to take and the time it is expected to depart.