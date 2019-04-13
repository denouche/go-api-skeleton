# go-api-skeleton

Get Go 1.12+: https://golang.org/dl/

Set your GOROOT to the install root.

Add go bin in your path.

Create a go root folder (GOPATH env variable). With sub directories src pkg bin.

In the src dir, create the github.com/my-user dirs.

In the my-user dir, clone the go-api-skeleton project under the name of your new project.

```
$GOPATH/
├── bin
├── pkg
└── src
    └── github.com
        └── my-user
            └── my-go-api
```

## Create my project

First, clean and initialize the git:

```bash
$ rm -rf .git
$ git init
```

Then run `./duplicate.sh` and follow the instructions to initialize your project.

This script will replace the template namespace and project name by your namespace and project name.

Then it will ask you for the entities to create, and it will create you the DAO funcs and basic CRUD APIs for this entities.

After creating the basic CRUD and DAO, you can:
- add your entities fields in the corresponding `./storage/model/*` file
- if you use postgres DAOs, you have to modify the SQL requests to take your new entities fields in account

After that you can execute the following commands to initialize dependencies and openapi schema:
```bash
$ make deps
$ make openapi
```

After that you can delete the `duplicate.sh` file and template files.
