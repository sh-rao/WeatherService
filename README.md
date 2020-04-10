# WeatherService
# Prelude
A simple http weather service that reports current temperature and wind speed.
There are two sources used by this service:
- http://api.weatherstack.com/current
- http://api.openweathermap.org/data/2.5/weather

These sources can be configured as PRIMARY and SECONDARY sources in config.yml by setting the Primary attribute.
If the service fails to obtain weather details from the priamry source it fails over to the secondary source.

## Design
`Handler` is responsible for calling the `WeatherService` to get current weather information and sending appropriate success or failure response.
`httputil` handles the actual writing of the response via `http.ResponseWriter`.
`WeatherService` does bulk of the work of making request to the weather providers and parsing/transforming the response to desired format represented by `WeatherDetails` struct.

### What could have been done better
I hate to use time as an excuse but given an extra few hours or so, I would imrovise the current design and implementation with the following:

- The configuration can be agnostic of the number of weather sources and they can be specifed as an array, without having to name them as primary and seconday.
- `main.go` has been kind of polluted by having to instantiate `config` and pass it to the `handler`. A better design would be to implement a `config service` which can be responsible for initialising the config and then injecting it into the `weather service`. Also an `Application Context` could have added value in defining and passing context to the services and also could have helped in customising the `http.Client` (this service uses `http.DefaultClient`).
- API Keys and APP IDs should never be stored in config. If this was hosted in cloud(assuming AWS as the default service provider) it would be stored in secrets manager (or may be system manager's parameter store) or injected as environment variable for non-cloud deployments.
- Separate out the implementation of extracting the weather details (`WeatherDetails` struct) into a separate service and get that service to ensure that the response from the weather providers/sources conform to a JSON schema. I have been reading about PACT a lot, so this would be perfect place to use PACT - contract based testing.
- Add some integration tests to make sure that, what's sent out by the service conforms to our JSON schema and also what comes in conforms to weather provider(s) schemas.
- Logging and error handling can be improved by passing context and using context logger whereever necessary.
  The downstream API errors can be handled better and then transformed to more meaningful errors but having said that, this is   not an API, it's a service.

# Prerequisites
- Make sure you have installed the latest version of Golang from https://golang.org/
  This service has been built and tested with go1.13 darwin/amd64 (on MacOS Mojave v10.14.6)
- Github account. You can sign-up for one here - https://github.com/join

# How to run the service
- Clone this repository into your local folder using git.
  e.g.
  ~~~
  mkdir WeatherService
  cd WeatherService
  git clone https://github.com/sh-rao/WeatherService.git
  ~~~
  
- From the project root folder (e.g. WeatherService), run this command to download all the dependencies
  ~~~
  go get -u ./...
  ~~~
  
- Start the service from the project root folder (e.g. WeatherService) by running the following command
  ~~~
  go run main.go
  ~~~

- You can get the current weather details for Melbourne(Australia) city from the command line by running
  ~~~
  curl "http://localhost:8080/v1/weather/"
  ~~~
  OR
  you can type or copy+paste http://localhost:8080/v1/weather/ in the browser of your choice.
  
  Note: The city has been made configurable in the config file.
  
 - Sample output would look like this
   ~~~
   {"wind_speed": 11, "temperature_degrees": 15}
   ~~~
  
  # Running unit tests
  Units tests can be run by the following command from the project's root folder (e.g. WeatherService)
  ~~~
  go test ./... -v
  ~~~
  
