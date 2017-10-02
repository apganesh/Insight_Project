
# Spotlite
- A geo-fence can be looked upon as a virtual perimeter that you can draw around any location on a map, and then target customers that enter that location.  A geo-fence could be used for ad targeting for local businesses.  Based on the geo location of the customer we could display the appropriate ads based on the geo-fences he is overlapping, instead of just using personalization recommendations.

- During the three week data engineering project, i attempted to provide a high throughput and low latency service which could provide all the overlapping geo-fences the user/device is overlapping with.  To show case this idea I have a simple use case, where a rider is matched with the taxi-driver's whose geo-fences'(the distance which they are willing to drive) are overlapping with the rider's location.  

## Pipeline



## Pipeline details
- The pipeline consists of ingesting the events, driver's ever changing geo-location and the request from the riders.  
- The events are then processed by the microservice which would match with the the drivers, whose geo-fences are overlapping.
- The matching is saved into a database for further analysis
- ADD A PICTURES HERE
- The pipeline is hosted on AWS cluster with:
 -- 2 node Kafka
 -- 2 node Redis
 -- 3 node Cassandra
 -- 4 node microservice

## Challenges 
- Ever changing location of the drivers' geo location poses a challenge to match with the rider
- Simple cross join of driver's location with riders location is very expensive

## Technologies 
- The events are caputred by Kafka and we have two consumers consuming the driver and rider events and passing them on the microservice to provide the matching.
- All the connectors between Kafka, Redis and Cassandra uses Go bindings.  


## Slides, Video, etc
- Demo: http://spotlite.xyz
- Slides: https://goo.gl/FShSZk
- Video: ![alt text](https://drive.google.com/open?id=0B5DP82LM5Bo7Q2hyYnNuOXlFNEE)
