"""
Created: nguyen_tran

Service that communicates with the backend through HTTP
"""

from http.client import HTTPConnection
import logging
from dto.settings import *


class IRemoteBackendService:
    # Provide the database with a device_id,
    # and it will return a pass
    #
    # POST /api/backend/register_device
    # { "device_id": "some_id", "model": "YoloBit" }
    def authenticate(self, device_id: str) -> str:
        pass
    
    # Login and retrieve a JWT token
    #
    # POST /api/backend/login
    def get_token(self, device_id: str, password: str) -> str:
        pass
    
    # Gets device settings, create new default one if not found
    #
    # GET /api/backend/settings
    def get_settings(self) -> DeviceSettingsResponse:
        pass


class RemoteBackendService(IRemoteBackendService):
    http_con: HTTPConnection = None
    token: str = ""

    def __init__(self, http_con: HTTPConnection) -> None:
        self.http_con = http_con
        return
    
    def get_token(self, device_id: str, password: str) -> str:
        self.http_con.request("POST", '/api/backend/login', json.dumps({
            "device_id": device_id,
            "password": password
        }))
        res = self.http_con.getresponse()
        raw = res.read().decode('utf-8')
        logging.debug(raw)
        payload = json.loads(raw)
        if int(res.status) != 200:
            message = payload["message"]
            raise Exception(
                f"Cannot retrieve token for device ID, error={message}")
        token = payload["token"]
        self.token = token
        return token
    
    def get_settings(self) -> DeviceSettingsResponse:
        self.http_con.request("GET", "/api/backend/settings", headers={
            "token": self.token,
        })
        res = self.http_con.getresponse()
        raw = res.read().decode('utf-8')
        logging.debug(raw)
        dto = DeviceSettingsResponse().from_json(raw)
        if int(res.status) != 200:
            message = dto["message"]
            raise Exception(
                f"Cannot retrieve token for device ID, error={message}")
        return dto

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
