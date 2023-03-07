"""
Created: nguyen_tran

Keeps the business logic for handling settings changes from user
"""

import ast
import logging
from services.serial_data import ISerialService
import models.commands
from paho.mqtt.client import MQTTMessage
from repositories.command import ICommandRepository


class ICommandService:
    # Takes in a device ID then initiates a topic dedicated for
    # listening to activity commands (shutdown, start,...) from backend through AWS
    def receive_command(self, device_id: str) -> None:
        pass
    
    # Takes in a device ID then initiates a topic dedicated for
    # listening to settings changes from backend through AWS
    def receive_settings(self, device_id: str) -> None:
        pass


class CommandService(ICommandService):
    _command_repository: ICommandRepository = None
    _serial: ISerialService = None

    def __init__(self, command_repository: ICommandRepository, serial: ISerialService) -> None:
        self._command_repository = command_repository
        self._serial = serial

    def receive_command(self, device_id: str) -> None:
        try:
            self._command_repository.subscribe_command(
                device_id=device_id,
                callback=self._on_command_received
            )
        except Exception as err:
            raise err
        return
    
    def receive_settings(self, device_id: str) -> None:
        try:
            self._command_repository.subscribe_settings(
                device_id=device_id,
                callback=self._on_settings_received
            )
        except Exception as err:
            raise err
        return
    
    def _on_settings_received(self, client, userdata, message: MQTTMessage):
        logging.debug(f"Settings ID={message.mid} payload={message.payload}")
        if message.dup is True:
            logging.debug(f"Received duplicate settings ID={message.mid}, skip")
            return
        parsed_dict = None
        try:
            parsed_dict = ast.literal_eval(message.payload.decode('utf-8'))
        except Exception as err:
            logging.error(f"Parsing command failed error={err}")
        action = parsed_dict["action"]
        if action == models.commands.UPDATE_SETTINGS:
            # Received update settings command
            # self._serial.write()
            logging.info(f"Received update settings command")
        return

    # Will be executed anytime a command is sent back from the backend
    # through AWS IoT Core
    def _on_command_received(self, client, userdata, message: MQTTMessage):
        logging.debug(f"Command ID={message.mid} payload={message.payload}")
        # If the message is a duplicate
        if message.dup is True:
            logging.debug(f"Received duplicate command ID={message.mid}, skip")
            return
        parsed_dict = None
        try:
            parsed_dict = ast.literal_eval(message.payload.decode('utf-8'))
        except Exception as err:
            logging.error(f"Parsing command failed error={err}")
        action = parsed_dict["action"]
        if action == models.commands.START:
            # Received a start signal
            # self._serial.write()
            logging.info(f"Received start signal")
        if action == models.commands.SHUTDOWN:
            # Received shutdown signal
            # self._serial.write()
            logging.info(f"Received shutdown signal")
        return

