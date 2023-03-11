import json
import logging

DATA_RATE: int = 2

class DataRateCommandSettings(dict):
    rate_in_seconds: int
    
    def __init__(self, rate_in_seconds: int):
        self.rate_in_seconds = rate_in_seconds
        dict.__init__(self,
                      rate_in_seconds=self.rate_in_seconds)
        
    def __init__(self, js: str):
        try:
            dict = json.loads(js)
        except Exception as err:
            logging.error(err)
            raise err