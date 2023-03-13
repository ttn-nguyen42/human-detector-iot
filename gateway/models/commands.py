import json 

# Commands available
SHUTDOWN = "shutdown"
START = "start"
UPDATE_SETTINGS = "update-settings"
ACTIVITY_CHECK = "is-active"
CHANGE_RATE = "change-rate"

# This data will be sent to us from the backend through MQTT topic


class CommandRequest(dict):
    action_id: str
    action: str
    payload: str

    def __init__(self, js: str):
        super().__init__()
        self.__dict__ = json.loads(js)
        return


# Responses
SUCCESS = "success"
FAILURE = "failure"

# This data will be sent from us to MQTT response topic


class CommandResponse(dict):
    action_id: str
    result: str
    payload: str

    def __init__(self, action_id: str, result: str, payload: str) -> None:
        self.action_id = action_id
        self.result = result
        self.payload = payload
        dict.__init__(self,
                      action_id=action_id,
                      result=result,
                      payload=payload)
