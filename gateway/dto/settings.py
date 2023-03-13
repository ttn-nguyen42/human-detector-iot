# From backend
import json


class DeviceSettingsResponse(dict):
    # In seconds
    data_rate: int
    device_id: str

    def __init__(self, data_rate: int = 3, device_id: str = "") -> None:
        self.data_rate = data_rate
        self.device_id = device_id
        dict.__init__(self, data_rate=data_rate, device_id=device_id)
        return

    def from_json(self, js: str) -> any:
        self.__dict__ = json.loads(js)
        return self
