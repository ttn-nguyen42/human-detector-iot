"""
Created: nguyen_tran

Data classes that is being used/passed from the service layer to the repository layer
"""


class SensorDataDto:
    heat_level: int = 0
    light_level: int = 0
    device_id: str = ""

    def __init__(self, device_id: str, heat_level: int, light_level: int) -> None:
        self.device_id = device_id
        self.heat_level = heat_level
        self.light_level = light_level
