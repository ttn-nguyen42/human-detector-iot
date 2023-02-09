"""
Created: nguyen_tran

Repository layer for listening to backend's command
"""


from network.common import MQTTBroker


class ICommandRepository:
    def subscribe_command(self, device_id: str, callback: callable) -> any:
        pass


class CommandRepository(ICommandRepository):
    # Implements ICommandRepository
    _broker: MQTTBroker = None
    _topic: str = "yolobit/command"

    def __init__(self, broker: MQTTBroker) -> None:
        self._broker = broker
        return

    def subscribe_command(self, device_id: str, callback: callable) -> any:
        topic = f"{self._topic}/{device_id}"
        try:
            self._broker.subscribe(topic=topic, func=callback)
        except Exception as err:
            raise err
        
