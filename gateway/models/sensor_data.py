"""
Created: nguyen_tran

Keeps the data models (data classes) for the repository layer here
"""


import time


class SensorData(dict):
    # Represents a sensor data packet that will be sent to the back-end

    device_id: str = ""

    humidity: int = 0
    temp: int = 0
    detected: bool = False
    timestamp: str = ""

    def __init__(self, device_id: str, humidity: int, temp: int, detected: bool, timestamp: str) -> None:
        self.humidity = humidity
        self.temp = temp
        self.detected = detected
        self.timestamp = timestamp
        self.device_id = device_id
        dict.__init__(
            self,
            device_id=device_id,
            temp=temp,
            humidity=humidity,
            detected=detected,
            timestamp=timestamp
        )
        

