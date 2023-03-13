import json
import logging


class DataRateCommandSettings(dict):
    rate_in_seconds: int

    def __init__(self, js: str):
        super().__init__()
        try:
            self.__dict__ = json.loads(js)
        except Exception as err:
            logging.error(err)
            raise err
