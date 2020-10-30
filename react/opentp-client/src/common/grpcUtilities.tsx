import { Error, StatusCode } from "grpc-web"


export function getGrpcErrorMessage( error : Error, prepend?: string) : string {
    if( prepend ) {
        prepend = prepend + ": "
    }

    let grpErrorCodeAsStr = error.message
    switch (error.code) {
        case StatusCode.ABORTED:
            grpErrorCodeAsStr = "Aborted" 
            break
        case StatusCode.ALREADY_EXISTS:
            grpErrorCodeAsStr = "Already Exists"
            break 
        case StatusCode.CANCELLED: 
        grpErrorCodeAsStr = "Cancelled"
            break
        case StatusCode.DATA_LOSS:
            grpErrorCodeAsStr = "Data Loss"
            break 
        case StatusCode.DEADLINE_EXCEEDED:
            grpErrorCodeAsStr = "Deadline Exceeded"
            break 
        case StatusCode.FAILED_PRECONDITION: 
        grpErrorCodeAsStr = "Failed Precondition"
            break
        case StatusCode.INTERNAL: 
        grpErrorCodeAsStr = "Internal"
            break
        case StatusCode.INVALID_ARGUMENT:
            grpErrorCodeAsStr = "Invalid Argument"
            break 
        case StatusCode.NOT_FOUND: 
        grpErrorCodeAsStr = "Not Found"
            break
        case StatusCode.OK: 
        grpErrorCodeAsStr = "OK"
            break
        case StatusCode.OUT_OF_RANGE: 
        grpErrorCodeAsStr = "Out Of Range"
            break
        case StatusCode.PERMISSION_DENIED: 
        grpErrorCodeAsStr = "Permission Denied"
            break
        case StatusCode.RESOURCE_EXHAUSTED:
            grpErrorCodeAsStr = "Resource Exhausted"
            break 
        case StatusCode.UNAUTHENTICATED: 
        grpErrorCodeAsStr = "Unauthenticated"
            break
        case StatusCode.UNAVAILABLE: 
        grpErrorCodeAsStr = "Unavailable"
            break
        case StatusCode.UNIMPLEMENTED:
            grpErrorCodeAsStr = "Unimplemented"
            break 
        case StatusCode.UNKNOWN: 
        grpErrorCodeAsStr = "Unknown"
            break
               
    }    

    return prepend + grpErrorCodeAsStr
}