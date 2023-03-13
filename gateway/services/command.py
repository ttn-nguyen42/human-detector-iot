"""
Created: nguyen_tran

Keeps the business logic for handling settings changes from user
"""

from services.serial_data import ISerialService
from models.commands import *
from paho.mqtt.client import MQTTMessage
from repositories.command import ICommandRepository
import threading
from dto.commands import *
import utils.config


class ICommandService:
    # Takes in a device ID then initiates a topic dedicated for
    # listening to activity commands (shutdown, start,...) from backend through AWS
    def receive_command(self, device_id: str) -> None:
        pass

    def send_response(self, res: CommandResponse):
        pass

    def run(self):
        pass


class CommandService(ICommandService):
    _command_repository: ICommandRepository = None
    _serial: ISerialService = None
    _deviceId: str = ""
    message_to_send = []
    messenger_thread = None
    lock = None

    def __init__(self, command_repository: ICommandRepository,
                 serial: ISerialService) -> None:
        self._command_repository = command_repository
        self._serial = serial
        self.lock = threading.Lock()

    def run(self):
        logging.info("Spinning up messenger thread")
        try:
            self.messenger_thread = threading.Thread(
                target=send_messages, args=(self,))
            self.messenger_thread.start()
        except Exception as err:
            logging.error(f"Unable to create messenger thread err={err}")
        return

    def receive_command(self, device_id: str) -> None:
        try:
            self._command_repository.subscribe_command(
                device_id=device_id,
                callback=self._on_command_received
            )
            self._deviceId = device_id
        except Exception as err:
            logging.error(f"Unable to setup command receiver err={err}")
            raise err
        return

    def send_response(self, res: CommandResponse):
        try:
            self._command_repository.send_command_response(
                device_id=self._deviceId, res=res)
            logging.info("Sent activity response")
        except Exception as err:
            logging.error(
                f"Unable to send activity check response err={err}")
        return

    # Will be executed anytime a command is sent back from the backend
    # through AWS IoT Core
    def _on_command_received(self, _, __, message: MQTTMessage):
        logging.debug(f"Command ID={message.mid} payload={message.payload}")
        # If the message is a duplicate
        if message.dup is True:
            logging.debug(f"Received duplicate command ID={message.mid}, skip")
            return
        try:
            payload_str = message.payload.decode('utf-8')
            req_dict = json.loads(payload_str)
            req = CommandRequest(req_dict)
        except Exception as err:
            logging.error(f"Parsing command failed error={err}")
            return
        action = req.action
        action_id = req.action_id
        payload = req.payload
        if action == START:
            # Received a start signal
            # self._serial.write()
            logging.info(f"Received start signal")
            return
        if action == SHUTDOWN:
            # Received shutdown signal
            # self._serial.write()
            logging.info(f"Received shutdown signal")
            return
        if action == ACTIVITY_CHECK:
            # Received activity check
            res: CommandResponse = CommandResponse(action_id=action_id,
                                                   result=SUCCESS,
                                                   payload="")
            logging.info(f"Received activity check")
            try:
                self._serial.write(ACTIVITY_CHECK)
            except Exception as err:
                logging.info(f"Unable to send command to device err={err}")
                res.result = FAILURE
            self.lock.acquire()
            self.message_to_send.append(res)
            logging.info(f"Added response to queue msg={res}")
            self.lock.release()
        if action == CHANGE_RATE:
            # Received data rate update
            res: CommandResponse = CommandResponse(action_id=action_id,
                                                   result=SUCCESS,
                                                   payload="")
            try:
                pl: DataRateCommandSettings = DataRateCommandSettings(payload)
                logging.info(f"Received data rate change to {pl.rate_in_seconds}s")
                utils.config.DATA_RATE = pl.rate_in_seconds
            except Exception as err:
                logging.error(
                    f"Unable to parse data rate command settings err={err}")
            res.result = FAILURE
            self.lock.acquire()
            self.message_to_send.append(res)
            logging.info(f"Added response to queue msg={res}")
            self.lock.release()
        return


def send_messages(service: CommandService):
    logging.info("Waiting for messages")
    while True:
        if len(service.message_to_send) == 0:
            continue
        service.lock.acquire()
        message = service.message_to_send.pop()
        service.lock.release()
        logging.info(
            f"Received response message, ready to send m={message}")
        service.send_response(message)
        continue
    return
