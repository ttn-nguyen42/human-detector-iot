"""
Created: nguyen_tran

Listens to command from backend and process
"""

# Register ourselves as a subscriber to MQTT broker
import logging
from services.command import ICommandService


def register_command_subscriber(device_id: str, service: ICommandService) -> None:
    try:
        service.run()
        service.receive_command(device_id=device_id)
    except Exception as err:
        logging.error(f"Cannot register command subscription: {err}")
    return
