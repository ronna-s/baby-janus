# baby-janus
The HelloFresh baby Janus challenge

### The challenge:
You will build an HTTP API gateway API (router) that can register routes [origin=>destination] at runtime.
The gateway will then listen to requests on the path defined in origin, and reverse proxy the destination.

Have a look at the project - there are 3 components inside the project:    
1. *app* - a test app that will register the endpoints `/hi` (to return hello world) and `/parts` to return a super secret HTML file.
2. *server* - a project to simulate of a cluster of servers, the servers will randomly register endpoints (the endpoints will be randomly spread accross all the servers) to return the different parts of the HTML file required by the test app.
3. *gateway* - API gateway (API) awaiting your implementation.

When you are done we will use docker to run everything (we will start 10 servers instances to simulate a cluster, the test app and the gateway app to allow communication between them all).
    
### Steps
1. Assign an http handler function to the path: `/register-endpoint`
2. The handler will receive post HTTP messages with the json format

```json
{"orig": "/some/path", "dest": "http://some_domain:port/dest"}
```
3. After parsing the request body and retreiving the endpoint data, register another http handler to the origin path, which will upon request reverse proxy the destination. Closures are your friends!!!
4. To test and execute, run: `docker-compose up baby-janus_gateway`, the docker command will run the tests and only if they pass will start the API gateway.
5. As usual, when the test runs properly, you will be able to see the result:
  
      ``` bash
      docker-compose up --scale baby-janus_server=10
      ```
      Now navigate to http://127.0.0.1:8080/parts to see the results.

### Too easy?
You may have noticed that your implementation doesn't support calling register-endpoint with different targets but the same origin.
Add a round robin mechanism that will allow registering mulitple endpoints with the same origin in the API while iterating over the destinations.

### A little docker
```bash
docker-compose build # build the images listed in your docker-compose file
docker-compose up # starts all containers listed in the docker-compose file
docker ps -a # lists all of your containers
docker exec -it <container-hash-id> bash # will start a bash command line for you inside the container (you can execute many commands - not just bash)
docker-compose run <container-name> bash # similar to docker exec except it will create a new container for you to run bash (or any other command).
docker-compose stop # stop all containers listed in the docker-compose file
docker-compose down # stops and removes all the containers (but not the images)
docker-compose down --rmi #also removes the images
```
