# Device settings
class DeviceSettings(dict):
    action: str
    # In seconds
    data_rate: int

    def __init__(self, action: str = "get", data_rate: int = 3) -> None:
        self.action = action
        self.data_rate = data_rate
        dict.__init__(self, action=action, data_rate=data_rate)
        return
