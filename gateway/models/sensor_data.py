"""
Created: nguyen_tran

Keeps the data models (data classes) for the repository layer here
"""


import time


class SensorData(dict):
    # Represents a sensor data packet that will be sent to the back-end

    heat_level: int = 0
    light_level: int = 0
    device_id: str = 0

    timestamp: float = 0.0

    def __init__(self, device_id: str, heat_level: int, light_level: int, timestamp: float) -> None:
        self.device_id = device_id
        self.heat_level = heat_level
        self.light_level = light_level
        self.timestamp = timestamp
        dict.__init__(
            self,
            device_id=device_id,
            heat_level=heat_level,
            light_level=light_level,
            timestamp=timestamp
        )
