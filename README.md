In app.env, see testing email at lastpass with title simple bank

To create a new API:
- Create new rpc_<new_api_name>.proto and use it to declare request and response param
- Edit service_simple_bank.proto with the new API path and desc
- Run `make proto`, see new pb.go files and updated service_simple_bank.pb.go with the new API interface
- Implement the interface function in `gapi` by creating new rpc_<new_api_name>.go

To create new migration:
- Update schema at doc/db.dbml
- Run make dbschema (ensure dbml2sql installed with make dbschema-install)
- Run make newmigration name=<migration_name>
- New migration will be added to db/migration
- Edit the new migration and check if the update is reflected in the db
- Run make sqlc to update the model in go code