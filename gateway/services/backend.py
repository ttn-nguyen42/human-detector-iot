"""
Created: nguyen_tran

Service that communicates with the backend through HTTP
"""

class IRemoteBackendService:
    # Provide the database with a device_id
    # and it will returns a pass
    #
    # POST /api/v1/register_device
    # { "device_id": "some_id", "model": "YoloBit" }
    def authenticate(self, device_id: str, model: str) -> str:
        pass
    
class RemoteBackendService(IRemoteBackendService):
    def __init__(self) -> None:
        return
        
    def authenticate(self, device_id: str, model: str) -> str:
        return "remote_password"