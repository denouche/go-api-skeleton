#!/bin/bash

OLD_PROJECT_NAME="github.com/denouche/go-api-skeleton"

echo "What is the namespace and project name? (Eg. github.com/foo/bar)"
read NEW_PROJECT_NAME

find . -iname '*.go' -exec sed -i "s|${OLD_PROJECT_NAME}|${NEW_PROJECT_NAME}|g" {} \;
sed -i "s|${OLD_PROJECT_NAME}|${NEW_PROJECT_NAME}|g" Makefile Dockerfile

main()
{
    while true
    do
        echo "What is the entity name you want to create? (name to be used in URL path, write it lower case, plural)"
        read ENTITY_NAME
        ENTITY_NAME_UP="${ENTITY_NAME^}"

        echo "What is the postgresql schema to use for this entity?"
        read ENTITY_SCHEMA

        cp handlers/template_handler.go handlers/${ENTITY_NAME}_handler.go
        sed -i -r "s/template/${ENTITY_NAME}/g" handlers/${ENTITY_NAME}_handler.go
        sed -i -r "s/Template/${ENTITY_NAME_UP}/g" handlers/${ENTITY_NAME}_handler.go

        cp storage/dao/postgresql/database_postgresql_template.go storage/dao/postgresql/database_postgresql_${ENTITY_NAME}.go
        sed -i -r "s/template/${ENTITY_NAME}/g" storage/dao/postgresql/database_postgresql_${ENTITY_NAME}.go
        sed -i -r "s/Template/${ENTITY_NAME_UP}/g" storage/dao/postgresql/database_postgresql_${ENTITY_NAME}.go
        sed -i -r "s/schema/${ENTITY_SCHEMA}/g" storage/dao/postgresql/database_postgresql_${ENTITY_NAME}.go

        cp storage/dao/fake/database_fake_template.go storage/dao/fake/database_fake_${ENTITY_NAME}.go
        sed -i -r "s/template/${ENTITY_NAME}/g" storage/dao/fake/database_fake_${ENTITY_NAME}.go
        sed -i -r "s/Template/${ENTITY_NAME_UP}/g" storage/dao/fake/database_fake_${ENTITY_NAME}.go

        cp storage/model/template.go storage/model/${ENTITY_NAME}.go
        sed -i -r "s/Template/${ENTITY_NAME_UP}/g" storage/model/${ENTITY_NAME}.go

        sed -i -r "/\/\/ start: template routes/{:next;N;/\/\/ end: template routes/{bend};bnext;:end;p;s|template|${ENTITY_NAME}|g;s|Template|${ENTITY_NAME_UP}|g}" handlers/handler.go

        sed -i -r "/\/\/ start: template dao funcs/{:next;N;/\/\/ end: template dao funcs/{bend};bnext;:end;p;s|template|${ENTITY_NAME}|g;s|Template|${ENTITY_NAME_UP}|g}" storage/dao/database.go

        echo "Done"
        echo "If you want to stop here, hit Ctrl+C"
    done
}

main
