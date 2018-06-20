#!/usr/bin/env python

import os
import sys
import requests
import json
import defusedxml.ElementTree as ET
import ConfigParser
import logging
import logging.handlers
import argparse


def create_users(config, logger, verify):

    # retrieve ams info
    ams_host = config.get("AMS", "ams_host")
    ams_project = config.get("AMS", "ams_project")
    ams_token = config.get("AMS", "ams_token")
    ams_email = config.get("AMS", "ams_email")
    users_role = config.get("AMS", "users_role")
    goc_db_url_arch = config.get("AMS", "goc_db_host")

    # retrieve authn info
    authn_host = config.get("AUTHN", "authn_host")
    authn_service_uuid = config.get("AUTHN", "service_uuid")
    authn_token = config.get("AUTHN", "authn_token")
    authn_service_host = config.get("AUTHN", "service_host")

    # services holds all different services that the users  might belong to(which translates to ams topics)
    # each service will have a list of users associated with it
    services = {}
    conf_services = config.get("AMS", "service-types").split(",")
    for srv_type in conf_services:

        # strip any whitespaces
        srv_type = srv_type.replace(" ", "")

        # user count
        user_count = 0

        # form the goc db url
        goc_db_url = goc_db_url_arch.replace("{{service-type}}", srv_type)
        logger.info("\nAccessing url: " + goc_db_url)
        logger.info("\nStarted the process for service-type:" + srv_type)

        # grab the xml data from goc db
        goc_request = requests.get(goc_db_url, verify=False)
        logger.info(goc_request.text)

        # users from goc db that don't have a dn registered
        missing_dns = []

        # srv_type
        srv_type = srv_type.replace(".", "-")
        services[srv_type] = []

        # build the xml object
        root = ET.fromstring(goc_request.text)
        # iterate through the xml object's service_endpoints
        for service_endpoint in root.findall("SERVICE_ENDPOINT"):
            service_type = service_endpoint.find("SERVICE_TYPE").text.replace(".", "-")

            # grab the dn
            service_dn = service_endpoint.find("HOSTDN")
            if service_dn is None:
                missing_dns.append(service_endpoint.find("HOSTNAME").text)
                continue

            # Create AMS user
            hostname = service_endpoint.find("HOSTNAME").text.replace(".", "-")
            sitename = service_endpoint.find("SITENAME").text.replace(".", "-")
            user_binding_name = service_type + "---" + hostname + "---" + sitename

            # reverse the dn and exclude the last slash
            service_dn = ",".join(x for x in service_dn.text.split("/")[::-1][:-1])
            project = {'project': ams_project, 'roles': [users_role]}
            usr_create = {'projects': [project], 'email': ams_email}

            # create the user
            ams_usr_crt_req = requests.post("https://"+ams_host+"/v1/users/" + user_binding_name + "?key=" + ams_token, data=json.dumps(usr_create), verify=verify)
            logger.info(ams_usr_crt_req.text)

            # if the response doesn't contain the field uuid, it means the user was not created
            req_data = json.loads(ams_usr_crt_req.text)
            if ams_usr_crt_req.status_code != 200:
                logger.critical("\nUser: " + user_binding_name)
                logger.critical("\nSomething went wrong while creating ams user.\nBody data: " + str(usr_create) + "\nResponse Body: " + ams_usr_crt_req.text)
            if "uuid" not in req_data:
                logger.critical("uuid field not found in response from ams")
                continue

            # Create the respective AUTH binding
            bd_data = {'name': user_binding_name, 'service_uuid': authn_service_uuid, 'host': authn_service_host, 'dn': service_dn, 'unique_key': req_data["uuid"]}
            authn_binding_crt_req = requests.post("https://"+authn_host+"/v1/bindings?key="+authn_token, data=json.dumps(bd_data), verify=verify)
            logger.info(authn_binding_crt_req.text)

            if authn_binding_crt_req.status_code != 201:
                logger.critical("Something went wrong while creating a binding.\nBody data: " + str(bd_data) + "\nResponse: " + authn_binding_crt_req.text)
                # delete the ams user, since binding could not be created
                requests.delete("https://"+ams_host+"/v1/users/" + user_binding_name + "?key=" + ams_token, verify=verify)
                continue

            # if both the user and binding have been created, assign the user to the acl of the topic
            services[service_type].append(user_binding_name)

            # count how many users have been created
            user_count += 1

        # modify the acl for each topic , to add all associated users
        authorized_users = services[srv_type]
        requests.post("https://"+ams_host+"/v1/projects/"+ams_project+"/topics/"+srv_type+":modifyAcl?key="+ams_token, data=json.dumps({'authorized_users': authorized_users}), verify=verify)

        logger.critical("\nService Type: " + srv_type)

        logger.critical("\nMissing DNS: "+str(missing_dns))

        logger.critical("\nTotal Users Created: " + str(user_count))

        logger.critical("\n-----------------------------------------")


def main(args=None):

    # set up the config parser
    config = ConfigParser.ConfigParser()

    # check if config file has been given as cli argument else
    # check if config file resides in /etc/argo-api-authn/ folder else
    # check if config file resides in local folder
    if args.ConfigPath is None:
        if os.path.isfile("/etc/argo-api-authn/conf.d/ams-create-users-gocdb.cfg"):
            config.read("/etc/argo-api-authn/conf.d/ams-create-users-gocdb.cfg")
        else:
            config.read("../conf/ams-create-users-gocdb.cfg")
    else:
        config.read(args.ConfigPath)

    # set up logging
    logger = logging.getLogger("AMS User create script")

    # stream(console) handler
    console_handler = logging.StreamHandler()
    logger.addHandler(console_handler)
    logger.setLevel(logging.INFO)

    # sys log handler
    syslog_handler = logging.handlers.SysLogHandler(config.get("LOGS", "syslog_socket"))
    syslog_handler.setFormatter(logging.Formatter('%(name)s[%(process)d]: %(levelname)s %(message)s'))
    syslog_handler.setLevel(logging.WARNING)
    logger.addHandler(syslog_handler)

    # start the process of creating users
    create_users(config, logger, args.Verify)


if __name__ == "__main__":

    parser = argparse.ArgumentParser(description="Create ams users and their respective bindings using data imported from goc db")
    parser.add_argument(
        "-c", "--ConfigPath", type=str, help="Path for the config file")
    parser.add_argument(
        "-verify", "--Verify", help="SSL verification for requests",  action="store_true")

    sys.exit(main(parser.parse_args()))
