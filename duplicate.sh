#!/bin/bash

command -v gsed 2>&1 /dev/null && SED_CMD=gsed || SED_CMD=sed

## check sed is GNU sed
${SED_CMD} --version 2>&1 > /dev/null
if [[ $? -ne 0 ]]
then
    echo "sed not compatible. If you are using Mac OS, please install GNU sed: brew install gnu-sed"
    exit 1
fi

# global vars
OLD_PROJECT_NAMESPACE=
OLD_PROJECT_NAME=
OLD_PROJECT_FULL_NAME=
NEW_PROJECT_NAMESPACE=
NEW_PROJECT_NAME=
NEW_PROJECT_FULL_NAME=
DAO_PG=
DAO_MONGO=
DAO_IN_MEMORY=
DELETE_TEMPLATES=1

init()
{
    if [[ -d "$GOPATH" ]]
    then
        BASE_DIR="$GOPATH/src"
    else
        # GOPATH not found, ask
        echo "GOPATH not found, What is your project directory? (give the absolute path)"
        read BASE_DIR
    fi

    OLD_PROJECT_NAMESPACE="github.com/denouche"
    OLD_PROJECT_NAME="go-api-skeleton"
    OLD_PROJECT_FULL_NAME="${OLD_PROJECT_NAMESPACE}/${OLD_PROJECT_NAME}"

    echo "What is the new namespace? (Eg. if you are creating project under 'github.com/foo/bar', enter 'github.com/foo')"
    read NEW_PROJECT_NAMESPACE
    echo "What is the new project name? (Eg. if you are creating project under 'github.com/foo/bar', enter 'bar'))"
    read NEW_PROJECT_NAME

    NEW_PROJECT_FULL_NAME="${NEW_PROJECT_NAMESPACE}/${NEW_PROJECT_NAME}"

    read -r -p "Do you want PostgreSQL DAO? [y/N] " response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])+$ ]]
    then
        DAO_PG=1
    else
        DAO_PG=0
    fi

    read -r -p "Do you want MongoDB DAO? [y/N] " response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])+$ ]]
    then
        DAO_MONGO=1
    else
        DAO_MONGO=0
    fi

    read -r -p "Do you want In Memory DAO? [y/N] " response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])+$ ]]
    then
        DAO_IN_MEMORY=1
    else
        DAO_IN_MEMORY=0
    fi
}

createProject()
{
    # create project
    mkdir -p ${BASE_DIR}/${NEW_PROJECT_NAMESPACE}
    cp -rf . ${BASE_DIR}/${NEW_PROJECT_FULL_NAME}
    cd ${BASE_DIR}/${NEW_PROJECT_FULL_NAME}

    # init git and readme
    rm -rf .git/
    git init
    echo -e "# ${NEW_PROJECT_NAME}\n\n" > README.md

    # change go imports
    find . -iname '*.go' -exec ${SED_CMD} -i "s|${OLD_PROJECT_FULL_NAME}|${NEW_PROJECT_FULL_NAME}|g" {} \;
    ${SED_CMD} -i "s|${OLD_PROJECT_FULL_NAME}|${NEW_PROJECT_FULL_NAME}|g" Makefile Dockerfile go.mod
    ${SED_CMD} -i "s|${OLD_PROJECT_NAME}|${NEW_PROJECT_NAME}|g" Makefile Dockerfile info.yaml cmd/root.go
}

createEntities()
{
    echo
    echo
    echo "Now we will create entities and corresponding CRUD."

    echo
    read -r -p "Do you want to a create a new entity and CRUD? [y/N] " response
    while [[ "$response" =~ ^([yY][eE][sS]|[yY])+$ ]]
    do
        createOneEntity
        echo
        read -r -p "Do you want to a create a new entity and CRUD? [y/N] " response
    done
}

createOneEntity()
{
    echo "Creating a new entity:"
    echo "What is the entity name you want to create? (name to be used in URL path, write it lower case, singular)"
    read ENTITY_NAME
    ENTITY_NAME_UP="$(tr '[:lower:]' '[:upper:]' <<< ${ENTITY_NAME:0:1})${ENTITY_NAME:1}" # because ENTITY_NAME_UP="${ENTITY_NAME^}" is for bash 4 only

    if [[ ${DAO_PG} -eq 1 ]]
    then
        echo "What is the postgresql schema to use for this entity? (if you plan to use MongoDB only, just type Enter)"
        read ENTITY_SCHEMA
    fi

    if [[ "$ENTITY_NAME" = "template" ]]
    then
        # in this case everything is ok, just change the SQL schema
        DELETE_TEMPLATES=0
        ${SED_CMD} -i -r "s/schema/${ENTITY_SCHEMA}/g" storage/dao/postgresql/database_postgresql_${ENTITY_NAME}.go
    else
        cp handlers/template_handler.go handlers/${ENTITY_NAME}_handler.go
        ${SED_CMD} -i -r "s/template/${ENTITY_NAME}/g" handlers/${ENTITY_NAME}_handler.go
        ${SED_CMD} -i -r "s/Template/${ENTITY_NAME_UP}/g" handlers/${ENTITY_NAME}_handler.go

        cp storage/dao/postgresql/database_postgresql_template.go storage/dao/postgresql/database_postgresql_${ENTITY_NAME}.go
        ${SED_CMD} -i -r "s/template/${ENTITY_NAME}/g" storage/dao/postgresql/database_postgresql_${ENTITY_NAME}.go
        ${SED_CMD} -i -r "s/Template/${ENTITY_NAME_UP}/g" storage/dao/postgresql/database_postgresql_${ENTITY_NAME}.go
        ${SED_CMD} -i -r "s/schema/${ENTITY_SCHEMA}/g" storage/dao/postgresql/database_postgresql_${ENTITY_NAME}.go

        cp storage/dao/mongodb/database_mongodb_template.go storage/dao/mongodb/database_mongodb_${ENTITY_NAME}.go
        ${SED_CMD} -i -r "s/template/${ENTITY_NAME}/g" storage/dao/mongodb/database_mongodb_${ENTITY_NAME}.go
        ${SED_CMD} -i -r "s/Template/${ENTITY_NAME_UP}/g" storage/dao/mongodb/database_mongodb_${ENTITY_NAME}.go

        cp storage/dao/mock/database_mock_template.go storage/dao/mock/database_mock_${ENTITY_NAME}.go
        ${SED_CMD} -i -r "s/template/${ENTITY_NAME}/g" storage/dao/mock/database_mock_${ENTITY_NAME}.go
        ${SED_CMD} -i -r "s/Template/${ENTITY_NAME_UP}/g" storage/dao/mock/database_mock_${ENTITY_NAME}.go

        cp storage/dao/fake/database_fake_template.go storage/dao/fake/database_fake_${ENTITY_NAME}.go
        ${SED_CMD} -i -r "s/template/${ENTITY_NAME}/g" storage/dao/fake/database_fake_${ENTITY_NAME}.go
        ${SED_CMD} -i -r "s/Template/${ENTITY_NAME_UP}/g" storage/dao/fake/database_fake_${ENTITY_NAME}.go

        cp client/model/template.go client/model/${ENTITY_NAME}.go
        ${SED_CMD} -i -r "s/Template/${ENTITY_NAME_UP}/g" client/model/${ENTITY_NAME}.go

        ${SED_CMD} -i -r "/\/\/ start: template routes/{:next;N;/\/\/ end: template routes/{bend};bnext;:end;p;s|template|${ENTITY_NAME}|g;s|Template|${ENTITY_NAME_UP}|g}" handlers/handler.go

        ${SED_CMD} -i -r "/\/\/ start: template dao funcs/{:next;N;/\/\/ end: template dao funcs/{bend};bnext;:end;p;s|template|${ENTITY_NAME}|g;s|Template|${ENTITY_NAME_UP}|g}" storage/dao/database.go

        ${SED_CMD} -i -r "/\/\/ Template export/{p;s/Template/${ENTITY_NAME_UP}/g}" storage/dao/fake/database_fake.go

        ${SED_CMD} -i -r "/\/\/ Template index/{p;s/Template/${ENTITY_NAME_UP}/g}" storage/dao/mongodb/database_mongodb.go
    fi
}

build()
{
    git add .
    git commit -am'first commit'
    make deps openapi
}

clean()
{
    if [[ ${DELETE_TEMPLATES} -eq 1 ]]
    then
        # remove template data
        ${SED_CMD} -i -r "/\/\/ start: template routes/{:next;N;/\/\/ end: template routes/{bend};bnext;:end;d}" handlers/handler.go
        ${SED_CMD} -i -r "/\/\/ start: template dao funcs/{:next;N;/\/\/ end: template dao funcs/{bend};bnext;:end;d}" storage/dao/database.go
        ${SED_CMD} -i -r "/\/\/ Template export/d" storage/dao/fake/database_fake.go
        ${SED_CMD} -i -r "/\/\/ Template index/d" storage/dao/mongodb/database_mongodb.go

        find . -iname '*template*' -exec rm {} \;
    fi

    # remove unwanted DAO
    if [[ ${DAO_MONGO} -eq 0 ]]
    then
        ${SED_CMD} -i -r '/\/\/ DAO MONGO/d' handlers/handler.go
        rm -rf ./storage/dao/mongodb
    fi

    if [[ ${DAO_PG} -eq 0 ]]
    then
        ${SED_CMD} -i -r '/\/\/ DAO PG/d' handlers/handler.go
        rm -rf ./storage/dao/postgresql
    fi

    if [[ ${DAO_IN_MEMORY} -eq 0 ]]
    then
        ${SED_CMD} -i -r '/\/\/ DAO IN MEMORY/d' handlers/handler.go cmd/root.go
        ${SED_CMD} -i -r '/(start-offline|db-in-memory)/d' Makefile
        rm -rf ./storage/dao/fake
    fi

    rm duplicate.sh
}

init
createProject
createEntities
clean
build
