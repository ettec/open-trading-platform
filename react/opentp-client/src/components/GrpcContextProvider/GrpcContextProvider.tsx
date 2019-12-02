import * as React from "react";
import * as grpcWeb from 'grpc-web'
import { createContext } from 'react';
import { string } from "prop-types";


const GrpcContext = createContext({
    serviceUrl : "",
    grpcMetaData: {}
});

export interface Props {
    serviceUrl : string
    username: string
    appInstanceId: string
}

export interface State {
    serviceUrl : string,
    grpcMetaData: {}
}

export default class GrpcContextProvider extends React.Component<Props, State> {

    constructor(props: Props) {
        super(props)

        let grpcMetaDataMap = new Map();
        grpcMetaDataMap.set("username", props.username)
        grpcMetaDataMap.set("appInstanceId", props.appInstanceId)

        this.state = {
            serviceUrl: props.serviceUrl,
            grpcMetaData: grpcMetaDataMap
        }
    }

    render() {
        return (
          <GrpcContext.Provider value={this.state}>
            {this.props.children}
          </GrpcContext.Provider>
        );
      }

}

export const GrpcConsumer = GrpcContext.Consumer;