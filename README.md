# focus

  Focus is an app that lets you control chaos, organise your life and **focus** on what is critical :grimacing:. Focus lets you create lists and add tasks to keep track of your day to day activities. You can also set self-imposed deadlines which will then notify you at the specified time. You will be able to track your progress as you get through the tasks one by one.

  This repo contains the backend of the project that I was working on, with a friend of mine, who handled the frontend. The codebase structure was inspired by a talk that I saw by Mat Ryer, so my focus was mainly on readability/glanceability. I have tried to make everything as explicit as possible to achieve that. Hopefully, even beginners to Golang will be able to easily find their way.

  We had plans to add some more functionality like adding a Trello like board and grouping users to teams and organisations, but as with all projects, this too will forever be unfinished :joy:.

## Running the project

  The server accepts configuration either through a file named config.yaml or through environment variables. By default, the server looks for config.yaml in the current directory, but a path to the yaml file can be specified via the `-config-path` flag. To read from environment variables instead, simply set the flag `-config-env` as true. 

  Examples for the configuration file can be found in the **config-example.yaml** file and the list of environment variables that need to be provided are given in **configenv-example**.

```  
  go run .
```
  If the config.yaml is in the same directory

```
  go run . -config-path /path/to/config.yaml
```
  to specify the path to config file

```
  go run . -config-env true
```
  to read from environment variables

  You can also build the project into a binary and run that similarly.

## Docker

  The project can also be run with Docker. The environment variables for the docker image can be passed in with the `--env-file` which reads in a file of environment variables. To build the docker image,

```
  docker build -t focus .
```
  and to run it

```
  docker run -it --env-file configenv -p 5000:5000 focus
```

## Tests

  I am slowly adding tests to the project. In the current stage, there are a few tests in the root directory to test the handlers and some in the mysql directory. To run the tests, move to the directory and run the command

  ```
    go test -v
  ```

## Kubernetes

  WIP