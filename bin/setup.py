from setuptools import setup

NAME = "ARGO API AUTHN SCRIPTS"

setup(
    name = NAME,
    version = "1",
    author = 'GRNET',
    description = 'Collection of useful scripts for interacting with the argo api authn service',
    long_description = 'Collection of useful scripts for interacting with the argo api authn service',
    url = 'https://github.com/ARGOeu/argo-api-authn',
    package_dir = {'argo_api_authn_scripts': 'argo-api-authn-scripts/'},
    packages = ['argo_api_authn_scripts'],
    install_requires=['defusedxml==0.5.0', 'requests==2.18.4']
    )