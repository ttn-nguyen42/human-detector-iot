#
"""
Created: nguyen_tran

Sets up AWS credentials and utility functions
for connecting to AWS services
"""

# Environment variables key
import os


AWS_IOT_CORE_CERT = "AWS_IOT_CORE_CERT"
AWS_IOT_CORE_PRIVATE = "AWS_IOT_CORE_PRIVATE"
AWS_IOT_CORE_PUBLIC = "AWS_IOT_CORE_PUBLIC"
AWS_IOT_CORE_ENDPOINT = "AWS_IOT_CORE_ENDPOINT"
AWS_IOT_CORE_ROOT_CA = "AWS_IOT_CORE_ROOT_CA"


class Certs:
    ca: str = ""
    cert: str = ""
    private_key: str = ""
    public_key: str = ""

    def __init__(self, cert: str, private_key: str, public_key: str, ca: str) -> None:
        self.cert = cert
        self.private_key = private_key
        self.public_key = public_key
        self.ca = ca


def get_certs_path() -> Certs:
    # Retrieve the certificates path
    # from environment variables
    cert_path = os.environ.get(AWS_IOT_CORE_CERT)
    if len(cert_path) == 0:
        raise Exception(
            "Missing certificate (.cert.pem) file path for AWS IoT Core")
    priv_key_path = os.environ.get(AWS_IOT_CORE_PRIVATE)
    if len(cert_path) == 0:
        raise Exception(
            "Missing private key (.private.key) file path for AWS IoT Core")
    public_key_path = os.environ.get(AWS_IOT_CORE_PUBLIC)
    if len(public_key_path) == 0:
        raise Exception(
            "Missing public key (.public.key) file path for AWS IoT Core")
    ca_cert = os.environ.get(AWS_IOT_CORE_ROOT_CA)
    if len(ca_cert) == 0:
        raise Exception("Missing root CA (root-CA.crt) for AWS IoT Core")
    return Certs(
        cert=cert_path,
        private_key=priv_key_path,
        public_key=public_key_path,
        ca=ca_cert
    )


def get_url_endpoint() -> str:
    # Retrieve AWS IoT Core URL endpoint
    # from environment variables
    url_endpoint = os.environ.get(AWS_IOT_CORE_ENDPOINT)
    if len(url_endpoint) == 0:
        raise Exception("Missing AWS IoT Core endpoint URL")
    return url_endpoint
