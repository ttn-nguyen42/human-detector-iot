# Device settings
class DeviceSettings(dict):
    # In seconds
    data_rate: int

    def __init__(self, data_rate: int = 3) -> None:
        self.data_rate = data_rate
        dict.__init__(self, data_rate=data_rate)
        return
