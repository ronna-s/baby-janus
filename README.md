# baby_janus
The HelloFresh baby Janus challenge

### The challenge:
You will build an HTTP gateway API (router) that can register routes [origin=>target] in runtime.
The gateway will then listen to requests on the path defined in origin, and redirect to the target.

Have a look at the project - there are 3 components inside the project:    
1. *app* - a test app that will register the endpoints `/hi` (to return hello world) and `/parts` to return a super secret HTML file.
2. *server* - a project to simulate of a cluster of servers, the servers will randomly register endpoints (the endpoints will be randomly spread accross all the servers) to return the different parts of the HTML file required by the test app.
3. *gateway* - gateway API awaiting your implementation.

When you are done we will use docker to run everything (we will start 10 servers instances to simulate a cluster, the test app and the gateway app to allow communication between them all).
    
### Steps
1. Assign an http handler function to the path: `/register_endpoint`
2. The handler will receive post HTTP messages with the json format

```javascript
{"origin": "/some/path", "target": "http://some_domain:port/target"}
```
3. After parsing the request body and retreiving the endpoint data, register another http handler to the origin path, which will upon request redirect to the target. Closures are your friends!!!
4. To test and execute, run: `docker-compose up baby_janus_gateway`, the docker command will run the tests and only if they pass will start the gateway API.
5. As usual, when the test runs properly, you will be able to see the result:
  
      ``` bash
      docker-compose up --scale baby_janus_server=10
      ```
      Now navigate to http://127.0.0.1:8080/parts to see the results.

### Too easy?
You may have noticed that your implementation doesn't support calling register_endpoint with different targets but the same origin.
If you remove the call to `t.skip()` inside the round robin test in `main_test.go`, you will find that your tests are failing again. Add a round robin mechanism that will allow registering mulitple endpoints with the same origin in the API and will iterate over the targets. 
