
## Licensing

```
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

## Description

Service is a project that provides a starter-kit for a REST based web service. It provides best practices around Go web services using POD architecture and design. It contains the following features:

* Minimal application web framework.
* Middleware integration.
* Database support using MongoDB.
* CRUD based pattern.
* Distributed logging and tracing.
* Testing patterns.
* User authentication.
* POD architecture with sidecars for metrics and tracing.
* Use of Docker, Docker Compose, and Makefiles.
* Vendoring dependencies with Modules, requires Go 1.11.
* Deployment to Azure using ACI.

### Go Installation

First of all export some paths, and save them in your .zshrc or .bashrc files for easy use. Use sudo if you get error.

# Go development

export GOPATH="${HOME}/.go"

export GOROOT="$(brew --prefix golang)/libexec"

export PATH="$PATH:${GOPATH}/bin:${GOROOT}/bin"

test -d "${GOPATH}" || mkdir "${GOPATH}"

test -d "${GOPATH}/src/github.com" || mkdir -p "${GOPATH}/src/github.com"

# Then finally install go, with Homebrew.

brew install go

# Dep installation

brew install dep


### Getting the project

You can use the traditional `go get` command to download this project into your configured GOPATH.

```
$ go get -u go_scaffold
```

### Building the project

Navigate to the root of the project and use the `makefile` to build all of the services.

```
$ cd $GOPATH/go_scaffold/cmd/api
$ go build
```

### Running the project

Navigate to the root of the project and use the `makefile` to run all of the services.

```
$ cd $GOPATH/go_scaffold
$ ./api
```

### Stopping the project

You can hit <ctrl>C in the terminal window running.

```
$ <ctrl>C
```

#### Authenticating

Before any authenticated requests can be sent you must acquire an auth token. Make a request using HTTP Basic auth with your email and password to get the token.

```
$ curl --user "admin@example.com:gophers" http://localhost:3000/v1/users/token
```

I suggest putting the resulting token in an environment variable like `$TOKEN`.

#### Authenticated Requests

To make authenticated requests put the token in the `Authorization` header with the `Bearer ` prefix.

```
$ curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users
```

#### Mongo DB installation

Install and Run MongoDB with Homebrew
Open the Terminal app and type brew update.
After updating Homebrew ``` brew install mongodb ```
After downloading Mongo, create the “db” directory. This is where the Mongo data files will live. You can create the directory in the default location by running mkdir -p /data/db
Make sure that the /data/db directory has the right permissions by running

> sudo chown -R `id -un` /data/db
> # Enter your password

Run the Mongo daemon, in one of your terminal windows run mongod. This should start the Mongo server.
Run the Mongo shell, with the Mongo daemon running in one terminal, type mongo in another terminal window. This will run the Mongo shell which is an application to access data in MongoDB.
To exit the Mongo shell run quit()
To stop the Mongo daemon hit ctrl-c