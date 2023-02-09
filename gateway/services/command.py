"""
Created: nguyen_tran

Keeps the business logic for handling settings changes from user
"""

import ast
import models.commands
from repositories.command import ICommandRepository
from awscrt import io, mqtt


class ICommandService:
    def receive_command(self, device_id: str) -> None:
        pass

# `topic` (str): Topic receiving message.
#
# `payload` (bytes): Payload of message.
#
# `dup` (bool): DUP flag. If True, this might be re-delivery
#                  of an earlier attempt to send the message.
# `qos` (: class: `QoS`): Quality of Service used to deliver the message.
#
# `retain` (bool): Retain flag. If True, the message was sent
#    as a result of a new subscription being made by the client.
#
# `**kwargs` (dict): Forward-compatibility kwargs.


def on_command_received(topic: str, payload: bytes, dup: bool, qos: mqtt.QoS, retain: bool, **kwargs: dict) -> None:
    # if dup == True:
    #     # Found a duplicate command
    #     return
    print("Received message")
    try:
        to_dict = ast.literal_eval(payload.decode('utf-8'))
    except Exception as err:
        print("Skipped, payload not a valid JSON")
        return
    print(to_dict)
    action = to_dict["action"]
    if action == models.commands.SHUTDOWN:
        # Received shutdown signal
        print("Shutdown")
        return
    if action == models.commands.START:
        # Received turn on signal
        print("Start")
        return
    if action == models.commands.UPDATE_SETTINGS:
        # Received action to update settings
        print("Update settings")
        return
    return


class CommandService(ICommandService):
    _commandRepository: ICommandRepository = None

    def __init__(self, commandRepository: ICommandRepository) -> None:
        self._commandRepository = commandRepository

    def receive_command(self, device_id: str) -> None:
        try:
            self._commandRepository.subscribe_command(
                device_id=device_id,
                callback=on_command_received
            )
        except Exception as err:
            raise err
        return
