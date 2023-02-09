
# Commands available
SHUTDOWN = "shutdown"
START = "start"
UPDATE_SETTINGS = "update_settings"

# This data will be sent to us from the backend through MQTT topic
class StatusSignal(dict):
    action: str
    status_signal: str

    def __init__(self, action: str, status: str) -> None:
        self.status_signal = status
        self.action = action
        dict.__init__(self, action=action, status_signal=status)
        
# Device settings
class DeviceSettings(dict):
    action: str
    data_rate: int
    
    def __init__(self, action: str, data_rate: int) -> None:
        self.action = action
        self.data_rate = data_rate
        dict.__init__(self, action=action, data_rate=data_rate)
