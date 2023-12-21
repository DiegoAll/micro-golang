

Open vs Code (Make sure nothing is open)

file > save workspace as (In a specific folder - micro-golang here is front-end unziped)

file > Add Folder to workspace  (micro-golang/front-end)  Add


So if I now browse the Explorer, you'll see that I have front end here. And as we add microservices, they'll all be at the same level as the front end folder, but they'll have their own names and **that just gives us a convenient means of having all the projects open at the same time**. 


### 11. Reviewing the frontend code


### 12. Our first service: the Broker


Now I want to create a new folder over here at the same level as front end. But if **I try to do that using this button, it's always going to add it to the front end folder because this is a workspace**. So instead I'll **go back to the finder and find that folder where front end exists. I'll create a new folder which I'll call broker-service**, and that's right beside front end.
I'll go back to Visual Studio code, go to the file menu and say add folder to workspace. (choose broker-service) click ADD.

QUE GONORREA COMO QUE TRABAJANDO CON WORKSPACE NO SE PUEDEN ABRIR MAS TERMINALES.
TOCA select current working directory

    cd broker-service/

    go mod init broker


I just want to make sure that **I can connect from the front end to the broker service.** So we'll make this dead simple. Now, over time, this broker will become much more complex. So I'm going to do one thing. It's going to take requests and forward them off to some microservice and then send a response back. So why don't we make our lives easier and install a really good third party rooting package? Okay, so I'm going to open my terminal, make sure I'm in the broker service folder. I am. And I'm going to go and get one of the routing packages that I really like go get and it's from GitHub

    go get github.com/go-chi/chi/v5

And I'll just add slash middleware because I'll use some of the middleware that's available to us.

    go get github.com/go-chi/chi/v5/middleware

we need to implement cause protection so we can actually communicate from the front end to the back end.

    go get github.com/go-chi/cors

So there's three packages I'm going to need.

It's called middleware dot heartbeat right there. And make sure you get version five it is and the heartbeat URL is just ping. This will allow me to if I want to at some point in the future really easily make sure that this service is still alive by checking for its heartbeat.

***Okay, so this will be a route, but I want to actually add a receiver here that allows me to share any configuration I might have from my application with routes in case I need them.***

	mux.Use(middleware.Heartbeat("/ping"))
    func (app *Config) Broker(w http.ResponseWriter, r *http.Request)


### 13. Building a docker image for the Broker service

Now there's a couple of ways of doing this, and I'm going to show you two ways. The first way is the multi-stage build using a certified go docker image. And the second way is much faster.

And one of the irritating things about Visual Studio code is when I'm in a workspace like this, I have my front end folder, **I have my broker's service folder I'd like to have right beside these at the same level a Dr. Compose But Visual Studio code only lets me add folders, not files to a workspace,** at least as far as I know.

create folder project (At root level)
Add a folder to Workspace (Choose project folder)

Now inside of that is where I'm going to put my Docker compose. But of course, Docker Compose needs something to run and that means a Docker file. **broker-service.dockerfile** And this will be the Docker file that tells Docker compose how to build the image.

We're going to build our code in this image and then build a smaller image.

Now I just build my go code so I'll run the command and I'm going to add an environment variable. See CGO_ENABLED=0. I'm not using any C libraries, just use the standard library and then it's the standard go build command,
When I run, this should first of all build all of my code on one Docker image and then create a much smaller Docker image and just copy over the executable. And that's exactly what I want.

Once **I'm going to say mode replicated, which I don't have to do, but I want you to get the habit of doing this because we'll be using it later on replicated.** And then one more return and then replicas one. **Now, in this case, we're only ever allowed to have one replica because we can't listen to with two Docker images on Port 88 on the local host. But later on when we implement Service Discovery, we'll be replicating some of our Docker images.**

```
    build:
      #Context: Directorio que se utiliza como origen para construir la imagen del contenedor.
      context: ./../broker-service
      dockerfile: ./../broker-service/broker-service.dockerfile
```
    CURRENT_UID=$(id -u):$(id -g) docker-compose up -d


But once it's going, you should be able to look at your Docker dashboard and see when you actually have this project running.

    (base) CO0C02GD0T7MD6M:project dposada$ docker logs project-broker-service-1
    2023/12/16 03:09:54 Starting broker service on port 80


So the next step is to leave that running and go back to our front end and add some logic in the HTML template that allows me to hit that broker service just to make sure that the two can talk to one another.


### 14. Adding a button and Javascript to the frontend

So we have our broker service running in the background and now we want to modify our front end to try
to hit that broker service. I just want to make sure I can send a request there and have my broker service respond with some kind of JSON.

I open the templates folder and we're going to open the test page to go HTML.

    let brokerBtn = document.getElementById("brokerBtn");
    let output = document.getElementById("output");
    let sent = document.getElementById("payload");
    let received = document.getElementById("received");

    mux.Post("/", app.Broker)
    fetch("http:\/\/localhost:8080", body)


So I save that and go back to my web browser, reload the page and try that again. Empty post request and I got a response back. That's all I care about at this point. I can, in fact, communicate between my front end and my back end. Perfect.


### 15. Creating some helper functions to deal with JSON and such

(...)

### 16. Simplifyng things with a Makefile (Linux)

    Multistage builkding | without Multistage build (Faster)
    make up_build (docker-compose up --build -d)
    make start (start the frontend)
    make stop (stop the frontend)
    make down (docker-compose down)(brings down my docker images)


### 17. Simplifyng things with a Makefile (Windows)

    && set GOOS=linux&& set GOARCH=amd64&&


## Section3: Building an Authentication Service


### 18. What we'll cover in this section

So in this section of the course, what we're going to do is **add another microservice, we'll add an authentication microservice that exists in Docker beside the broker.** And what **we want to do is have the user try to authenticate through the broker. The broker calls the authentication microservice determines whether or not that user is able to authenticate and sends back the appropriate response.** Now, obviously, **that means that the authentication microservice has to have some kind of persistence.** We need to be able to store user information somewhere.

So what we'll do is **add a Postgres database that's specific to the authentication service**. That means we could use this authentication service for any kind of application.
**If you have multiple services and you have the same users for all of those services, you could use this microservice for authentication** for all of them.
Now you'll notice that I have an arrow from the user. In other words, **the person with the browser going to both the broker and to the authentication service. And that's because we don't necessarily have to use the broker to authenticate. We could contact the authentication service directly, provided we have its port exposed to the clear internet, of course.**

Either way works and chances are we'll do it both ways before the end of this course. But **initially I want to use that broker as the single point of entry and you'll see why as we move on and put rabbit MQ and other kinds of technologies into play.** But in any case, let's get started with our authentication microservice.

***PENDIENTE: HASTA ESTE PUNTO SOLO UTILIZARAN EL BROKER PARA AUTENTICAR, VALIDAR DONDE LO HARAN DIRECTAMENTE MAS ADELANTE***


<img src="/project/brokerservice_auth.jpeg">

<img src="/project/brokerservice_auth.jpeg" align="right" width="300px" height="300px" vspace="20" />



### 19. Setting up a stub Authentication service

**I'll create a new folder here and I'll call nine authentication service** and I suggest you do as well because it'll make working with the Makefile and the Docker files much easier if we have the same naming convention, **and then add it to my workspace.**

    go mod init authentication

Now this is going to be a web service, but it's an API service. And of course we need to listen on a specific port. And just to make things clean up here at the top, I'm going to declare a constant, which I'll call web port. And I'm like, that equal to the string of eight. **We're going to listen on Port 80 and we can do that, of course, even though the broker is listening on Port eight, because Docker lets multiple containers listen on the same port and treats them as individual servers.**

Now **this is not a course in how to create data models for Postgres** and so forth. So to save some time, if you go to the course **resources for this lecture, you'll find a file called models.go,** dot, zip, download that, extract it and copy the contents of that file to this model dog go file and I'll just paste it in here because I have it on my clipboard already and we'll go through it pretty quickly.

How long should I go before it's considered to be a failed operation? And I set that to 3 seconds, which is pretty straightforward.

Create routes.go

    go get -u github.com/go-chi/chi/v5
    go get -u github.com/go-chi/chi/v5/middleware
    go get -u github.com/go-chi/cors


Okay, so **there is a barebones version of our application of our authentication service.**

**Se descargo models.go**


### 20. Creating and connecting to Postgres from the Authentication service

We're going to be connecting to Postgres. So that means we need to get some drivers and we'll use the ones that are currently the best recommended ones. And that's Jack C slash PDX.

    go get github.com/jackc/pgconn
    go get github.com/jackc/pgx/v4
    go get github.com/jackc/pgx/v4/stdlib

It'll be at the end and those are necessary drivers. Okay, so let's close our terminal window and let's connect to the database.
We'll create a new function and I'll call this one openDB and it takes a DSN, which is a string.

The problem is we're going to add Postgres to our Docker composed dot y email file and we need to make sure that it's available before we actually return the database connection because this service might start out before the database does. So what we'll do is **we'll create another function here below this one called Func ConnectToDB()**.

Get on right there and we're going to look for DSM. Now I just create an infinite for loop and I stay in there until I successfully connect to the database. And it's pretty straightforward.

Let's say that the database crashes. The database refuses to come up for some reason. So I actually want to try this only a certain number of times, so maybe I'll back off for 2 seconds. So what I'll do is I'll create a new variable right at the top. Right up here. I'll create a variable called counts and I'll make it of type in 64. Okay, now back down here in this connect to DB function or connect.

So now we'll try to connect the database, but we'll only try for a total of, say, 20 seconds. So we have ten counts here and we're waiting for 2 seconds each time and that should get us connected.

Now let's actually connect the database so we'll get rid of the TB word here. We don't have any routes, so we can't do anything yet, but **we should be able to start it and connect to the database. So the last thing we have to do before we can actually connect to the database is come up here, including the necessary imports. So I'll just put them right here. We're going to use the branch identifier because we're not going to directly call this, but we have to have access to it.

We should be able to add Postgres to our Docker compose to set up an entry for our authentication service and set the appropriate environment variable DSN and actually ping the database when this runs.


### 21. A note about PostgreSQL

An important note about Postgres Version In the next lecture, I'm going to ask you to add Postgres to your docker-compose.yml, and I suggest that you use something like this:


### 22. Updating our docker-compose.yml for Postgres and the Authentication service

```
postgres:
    image: 'postgres:14.2'
    ports:
      - "54322:5432"
    deploy:
      mode: replicated
      replicas: 1
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: users
    volumes:
      - ./db-data/postgres/:/var/lib/postgresql/data/
```
And inside of that, inside the DB data folder or create another folder called Postgres. And this is where the Docker container will actually store my Postgres database. So I have a local copy of my computer and I know where it is and that's just convenient for me.

Now the other one we want, of course, is going to be the authentication service. So I'll put that right here between these two and I'll backspace to the correct level. The next step is to go over to our **authentication service folder and create a Docker file**. Copy the one from broker service, but we'll just modify it. And all I have to remember is when I actually build this application, I build it with the name of the app. So let's go to our Makefile and I'll do the Mac and Linux version first

Now let's go specify the auth binary, which I called off app.

```
AUTH_BINARY=authApp

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"
```

    make up_build

```
chown: changing ownership of '/var/lib/postgresql/data': Operation not permitted
```
    chmod -R 777 /db-data/postgres (no funciono como de costumbre)

Funciono harcodeando el UID+GID en el makefile (ver README.md)

Postgres. It's still initializing itself and it might take more than 20 seconds, so we might have to stop it and started again, but only the first time. After this, it will have the database all initialized and it should be much faster to start up.

Auth service
```
(base) CO0C02GD0T7MD6M:mercadolibre dposada$ docker logs project-authentication-service-1
2023/12/20 04:36:47 Starting authentication service
2023/12/20 04:36:47 Postgres not yet ready ...
2023/12/20 04:36:47 Backing off for two seconds ...
2023/12/20 04:36:49 Postgres not yet ready ...
2023/12/20 04:36:49 Backing off for two seconds ...
2023/12/20 04:36:51 Postgres not yet ready ...
2023/12/20 04:36:51 Backing off for two seconds ...
2023/12/20 04:36:53 Postgres not yet ready ...
2023/12/20 04:36:53 Backing off for two seconds ...
2023/12/20 04:36:55 Postgres not yet ready ...
2023/12/20 04:36:55 Backing off for two seconds ...
2023/12/20 04:36:57 Postgres not yet ready ...
2023/12/20 04:36:57 Backing off for two seconds ...
2023/12/20 04:36:59 Postgres not yet ready ...
2023/12/20 04:36:59 Backing off for two seconds ...
2023/12/20 04:37:01 Connected to Postgres!
```

Broker service
```
(base) CO0C02GD0T7MD6M:mercadolibre dposada$ docker logs project-broker-service-1
2023/12/20 04:36:47 Starting broker service on port 80
```

Postgres
```
2023-12-20 04:37:01.082 UTC [1] LOG:  database system is ready to accept connections
2023-12-20 06:33:25.711 UTC [67] LOG:  invalid length of startup packet
2023-12-20 06:33:26.742 UTC [68] LOG:  invalid length of startup packet
```

INTERESANTE! (ASI CONECTAN Y ESTAN EN LA MISMA RED) DOCKER COMPOSE CREA LA NETWORK.
LA CONEXION ES ENTRE CONTENEDORES, EL PUERTO DE LA IZQUIERDA ES PARA LA MAQUINA LOCAL.

    environment:
      DSN: "host=postgres port=5432 user=postgres password=password dbname=users sslmode=disable timezone=UTC connect_timeout=5"

    ports:
      - "54322:5432"

EN RESUMEN: (TOCO HARDCODEAR EN EL MAKEFILE)

	CURRENT_UID=164865804:1010544492 docker-compose up --build -d


**2023/12/20 04:37:01 Connected to Postgres!**

There it connected to Postgres. So now we know we can actually bring up the authentication service and we can connect to the database. And we do have an empty database in Postgres named user users plural. So now we need to go create a table in there and put some data in it so we can actually test this authentication service and we'll get started


### 23.Populating the Postgres database

Paste that in the query window and then we'll select all and run the selection. And at this point I now have a user's table and if I double click on it, we have one entry in there. Okay, so now we have something we can actually work with and we can start working on actually trying to authenticate through the broker and back to the end user


### 24. Adding a route and handlers to accept JSON

So one of the things we're going to do here is this microservice authentication right now, all it's going to do is listen for a post request that has a Jason body with a username and password. Then it will use our data models, the user type in there to check to see if the password and username supplied match what we have in our database and we'll send some response back. And of course, that means we're going to be reading and writing.

Create a new file called Helpers Dot. Go and paste the contents in there. Okay. And that will work just fine.

It's dead simple. It's very, very easy to make this a reusable package. And all I would do is, as you see here under installation, go get this file and use it like this. Create a variable called tools of type toolbox type tools. And then I have access to all of the methods in there. And if you look at the tools go, I actually have one for I have the type for JSON response that I can import where I need to use it.
**https://github.com/tsawler/toolbox**

And the only difference of course is that each of these three functions related to Jason, they're exported and they start with a capital letter. But for the sake of simplicity in this course, I'll just copy and paste for now.

Let's go back to our roots file and set up a root to a handler that doesn't exist yet.

So first of all, put in that heartbeat, I might as well mark start use and I'm going to use middleware dot heartbeat right there and I'll just pass it slash ping. I may never use this, but it's there if I need it.

We'll **create a new route, It'll be a post request and it's going to be /authenticate and it will go to a handler that doesn't exist yet. app.authenticate seems appropriate**. Okay, so we need to go make that handler. So I'll create a new file in the cmd/api folder of the authentication service folder. New file called Handler Go.

***Authenticate() Implementation***

**If we did everything right and I think we did, we should be able to now go and modify our front end application to make a request to the broker. And the next step, of course, is to have that broker receive that request, fire off the request to this service and then send a response back to the user.  (PASAMANOS)**


### 25. Update the Broker for a standar JSON format, and connect to our Auth service


**So we have a first version of our authentication microservice up and running in Docker and we have a database with some information in it. And the next step before we can actually try this out is to modify the broker application to the broker microservice to listen for a request from the front end, then fire a request off to the authentication microservice, receive the response from the microservice and send some kind of response back to the end user.**




### 26. Updating the front end to authenticate theought the Broker and trying things out

TENER VARIOS SERVICIOS ABIERTOS SIN TENER INCONVENIENTES POR LOS MODULOS.
RECORDAR TEMA GOPLS, SI NO HAY MODULO


https://ron-liu.medium.com/what-canonical-http-header-mean-in-golang-2e97f854316d 

In Golang, the header in http.Request is key canonicalised by default.
If you want to set the non-canonicalised key in header, you can directly assign to the map.