# events-go
This is an internal service which is used to save the events of each user.
This service saves the access time and door access by each user .

# Endpoints 
GET /getevents - This endpoint takes username as query param and returns all the events for the user.

POST /updateevent - This endpoint is used to update the successfull access time and door name for the user.

# Database
Dbname - events

# Database Instance Name 
db-events- <env>

  env can be test,perf or prod
  
# Database export 
  events.json
