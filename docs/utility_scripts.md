# Python utility scripts for easier interaction with the service

| Script | Description | Shortcut |
|--------|-------------|---------- |
| ams-create-users-gocdb.py | Python script that creates ams users and binding.| [Details](#ams-create-users-gocdb) |

<a id="ams-create-users-gocdb"></a>
## AMS Create users from goc db script
Python utility script that takes an xml feed from goc db,creates the respective
ams users under the specified project, assigns to the correct project's topic and
finally creates the binding for each user, using the dn from goc db.

`ams-create-users-gocdb.py -c <ConfigPath> -verify`

`-c : Path to an appropriate config file.If not specified
it will first look at /etc/argo-api-authn/conf.d/ams-users-create-gocdb.cfg
and then will look at projects conf folder`
`-verify: If specified all the requests will check the validity of the ssl certificate`
##### Configuration
Use the `ams-create-users-gocdb.template` to produce your conf file.
The project should exist in AMS in advance.
The service types specified should also be present as topics in ams under the specified project
```buildoutcfg
[AMS]
# under which ams project, the users will be created
ams_project:
# goc db url to pull user data
goc_db_host:
# service types referes to the different service types that will we should keep from the xml and assign them to the respectivew ams topic 
service-types:
# ams use role
users_role: publisher
# token to access ams
ams_token:
# ams url
ams_host:
# ams user email - since we don't get an email from goc db, we can use this field as a wildcard, to identify which users were created with this script
ams_email: goc_db_user

[AUTHN]
# token to access authn
authn_token:
# authn url to create bindings
authn_host: 
# service's uuid where bindings will belong
service_uuid:
# service's host where bindings will belong
service_host:

[LOGS]
syslog_socket:
```

## Requirements
There is a requirements.txt file inside the repo's `bin` folder that specifies which dependencies are needed.