import { Error, StatusCode } from "grpc-web"


export function getGrpcErrorMessage( error : Error, prepend?: string) : string {
    if( prepend ) {
        prepend = prepend + ": "
    }

    switch (error.code) {
        case StatusCode.PERMISSION_DENIED:
            return prepend + "Permission Denied"
            
            default:
                return prepend + error.message
    }    
}