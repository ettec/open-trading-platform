import * as React from "react";
import Container from "../Container";
import { InputGroup, Button } from "@blueprintjs/core";
import v4 from 'uuid';
import { createContext } from 'react';
import GrpcContextProvider from "../GrpcContextProvider";
import { Metadata } from "grpc-web";
import { Listing } from "../../serverapi/listing_pb";

export interface Props {
    children?: React.ReactNode
}

export interface State {
    loggedIn : boolean
}

export interface GrcpContextData {
    appInstanceId : string
    serviceUrl : string,
    grpcMetaData: Metadata
}




export default class Login extends React.Component<Props, State> {

    static grpcContext : GrcpContextData

    username: string
    appInstanceId: string

    constructor(props: Props) {
        super(props)

        this.appInstanceId = v4();

        this.username = "bert"

        Login.grpcContext = {
            serviceUrl : 'http://192.168.1.100:32365', 
            grpcMetaData : {"user-name": this.username, "app-instance-id": this.appInstanceId},
            appInstanceId : this.username + "@" + this.appInstanceId
        }


        this.state = {
            loggedIn : true
        }

        this.handleUserNameChange = this.handleUserNameChange.bind(this);
        this.onSubscribe = this.onSubscribe.bind(this);
    }

      handleUserNameChange(e:any) {
        if( e.target && e.target.value) {
            this.username = e.target.value;

            Login.grpcContext = {
                serviceUrl : 'http://192.168.1.100:32365', 
                grpcMetaData : {"user-name": this.username, "app-instance-id": this.appInstanceId},
                appInstanceId : this.username + "@" + this.appInstanceId
            }

          }
      }
    
      onSubscribe() {

        this.setState({loggedIn:true})

      }


    render() {

        if( this.state.loggedIn ) {
            
            return (
                <GrpcContextProvider serviceUrl='http://192.168.1.100:32365' username={this.username} appInstanceId={this.appInstanceId} >
                    <Container ></Container>
                </GrpcContextProvider>
            )
        } else {
            return (
                <div>
                    <InputGroup onChange={this.handleUserNameChange} />
                <Button onClick={this.onSubscribe } >Login</Button>
                </div>
            )
        }

        
    }
}
