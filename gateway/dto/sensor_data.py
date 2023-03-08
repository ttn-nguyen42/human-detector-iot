"""
Created: nguyen_tran

Data classes that is being used/passed from the service layer to the repository layer
"""
import json

class SensorDataDto:
    device_id: str = ""
    
    humidity: int = 0
    temp: int = 0
    detected: bool = False

    def __init__(self, device_id: str, humidity: int, temp: int, detected: bool) -> None:
        self.humidity = humidity
        self.temp = temp
        self.detected = detected
        
    
    def __init__(self, js: str):
        self.__dict__ = json.loads(js)
        return
