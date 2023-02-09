"""
Created: nguyen_tran

Listens to command from backend and process
"""

# Register ourselves as a subscriber to MQTT broker
from services.command import ICommandService


def register_command_subscriber(device_id: str, service: ICommandService) -> None:
    try:
        service.receive_command(device_id=device_id)
    except Exception as err:
        exit()
    return