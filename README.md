# base-service
A templated starting point for uniform microservices

### Prerequisites

#### NATS Server
Generate required TLS certificates:
```
sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout /etc/ssl/private/uniform-nats.key -out /etc/ssl/certs/uniform-nats.crt
sudo chmod +r /etc/ssl/private/uniform-nats.key
```
Then install NATS server:
```
sudo su
nats_latest_version=$(curl -i https://github.com/nats-io/nats-server/releases/latest | grep location: | sed 's/location: https:\/\/github.com\/nats-io\/nats-server\/releases\/tag\///g' | sed 's/.$//')
nats_latest_zip=$(echo https://github.com/nats-io/nats-server/releases/download/$nats_latest_version/nats-server-$nats_latest_version-linux-amd64.zip)
rm -f nats-server.zip
rm -rf nats-server-$nats_latest_version-linux-amd64
curl -L $nats_latest_zip -o nats-server.zip
unzip -o nats-server.zip
mv nats-server-$nats_latest_version-linux-amd64/nats-server /usr/bin/nats-server
```
Then run the NATS server: 
```
nats-server --tls --tlscert /etc/ssl/certs/uniform-nats.crt --tlskey /etc/ssl/private/uniform-nats.key
```

### Getting Started
First step is to compile resources and metadata into the project's source-code:
```
go generate
```
This will create a file `service/meta.go` which is ignored by the `.gitignore` and contains the project's resources and metadata.
Create `.description` file to override the description as pulled from Github.
Then ensure you set the `AppClient`, `AppProject` and `AppService` constants in the `service/run.go` file before doing anything else.

### CLI Commands
_Note that a command will require at least one running service node in order for command to be executed._

cmd command example `cmd/command.example-one.go`:
```
package cmd

import (
	"github.com/spf13/cobra"
	"service/service"
)

var exampleOneCmd = &cobra.Command{
	Use:   "command:example-one",
	Short: "Request the running " + service.AppName + " to execute the example-one command",
	Long:  "Request the running " + service.AppName + " to execute the example-one command",
	Run: func(cmd *cobra.Command, args []string) {
		service.Command("example-one", natsUri, compileNatsOptions())
	},
}

func init() {
	rootCmd.AddCommand(exampleOneCmd)
}
```

service command example `service/command.example-one.go`:
```
package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local(command("example-one")), exampleOne)
}

func exampleOne(r uniform.IRequest, p diary.IPage) {
	// todo: write logic here
}
```

### Background Worker Process

Use a CLI Command and add it to a scheduled cronjob to avoid the background process from being executed multiple times when scaling service instances.
In other words this will work like a sync.Mutex but across all running instances of the given service, allowing us to add as many service instances as we need.

### Routines

service action example `service/routine.example-two.go`:
```
package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("example-two"), exampleTwo)
}

func exampleTwo(r uniform.IRequest, p diary.IPage) {
	// todo: write logic here
}
```

### Events

service event example `service/event.example-three.go`:
```
package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("example-three"), exampleThree)
}

func exampleThree(r uniform.IRequest, p diary.IPage) {
	// todo: write logic here
}
```

### Maintenance

#### Sync Template Repository
```
git remote add template git@github.com:go-uniform/base-service.git
git fetch template main
git merge template/main --allow-unrelated-histories
```