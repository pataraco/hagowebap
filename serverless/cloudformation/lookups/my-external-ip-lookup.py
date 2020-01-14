"""Stacker custom lookup to get my current external IP."""

import logging
import requests

TYPE_NAME = 'MyExternalIp'
LOGGER = logging.getLogger(__name__)

def handler(value, provider, **kwargs):  # pylint: disable=W0613
    """ Lookup my external IP by using provided URL endpoint.

    Can specify a URL endpoint to use, otherwise DEFAULT_ENDPOINT is used

    For example:

    [in the stacker yaml (configuration) file]:

      lookups:
        MyExternalIp: lookups.my-external-ip-lookup.handler

      stacks:
        variables:
          MyCidr: ${MyExternalIp https://checkip.amazonaws.com}
    """

    ip = requests.get(value).text.strip()
    LOGGER.debug('external IP: %s', ip)
    return ip
