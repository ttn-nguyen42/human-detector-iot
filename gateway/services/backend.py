"""
Created: nguyen_tran

Service that communicates with the backend through HTTP
"""

from http.client import HTTPConnection
import json
import logging


class IRemoteBackendService:
    # Provide the database with a device_id,
    # and it will return a pass
    #
    # POST /api/v1/register_device
    # { "device_id": "some_id", "model": "YoloBit" }
    def authenticate(self, device_id: str) -> str:
        pass


class RemoteBackendService(IRemoteBackendService):
    http_con: HTTPConnection = None

    def __init__(self, http_con: HTTPConnection) -> None:
        self.http_con = http_con
        return

    def authenticate(self, device_id: str) -> str:
        self.http_con.request("POST", "/api/backend/register_device", json.dumps({
            "device_id": device_id
        }))

        res = self.http_con.getresponse()
        raw = res.read().decode('utf-8')
        logging.debug(raw)
        payload = json.loads(raw)
        if int(res.status) != 201:
            message = payload["message"]
            raise Exception(f"Cannot retrieve password for device ID, error={message}")
        password = payload["password"]
        return password
