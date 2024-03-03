In app.env, see testing email at lastpass with title simple bank

To create a new API:
- Create new rpc_<new_api_name>.proto and use it to declare request and response param
- Edit service_simple_bank.proto with the new API path and desc
- Run `make proto`, see new pb.go files and updated service_simple_bank.pb.go with the new API interface
- Implement the interface function in `gapi` by creating new rpc_<new_api_name>.go