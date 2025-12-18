from fastapi import HTTPException
import grpc

def handle_grpc_error(error: grpc.RpcError) -> HTTPException:
    if error.code() == grpc.StatusCode.UNAVAILABLE:
        return HTTPException(status_code=503, detail="IDM service unavailable")
    elif error.code() == grpc.StatusCode.INVALID_ARGUMENT:
        return HTTPException(status_code=400, detail="Invalid request")
    elif error.code() == grpc.StatusCode.NOT_FOUND:
        return HTTPException(status_code=404, detail="Resource not found")
    elif error.code() == grpc.StatusCode.PERMISSION_DENIED:
        return HTTPException(status_code=403, detail="Permission denied")
    elif error.code() == grpc.StatusCode.UNAUTHENTICATED:
        return HTTPException(status_code=401, detail="Unauthenticated")
    else:
        return HTTPException(status_code=500, detail="Internal server error")
