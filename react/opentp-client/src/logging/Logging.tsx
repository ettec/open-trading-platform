import { Error } from "grpc-web"




export function logDebug(msg: string ) {
    console.log("DEBUG:" + msg)
}

export function logError(msg: string ) {
    console.log("ERROR:" + msg)
}

export function logGrpcError(msg: string, err :Error ) {
    console.log("ERROR:" + msg + ":" + err.code + ":" + err.message)
}