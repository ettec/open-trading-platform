import { AnchorButton, Button, Classes, Dialog, InputGroup, Intent, Tooltip } from "@blueprintjs/core";
import { Error, Metadata } from "grpc-web";
import * as React from "react";
import v4 from 'uuid';
import { ReactComponent as ReactLogo } from '../../opentp.svg';
import { LoginServiceClient } from "../../serverapi/LoginserviceServiceClientPb";
import { LoginParams, Token } from "../../serverapi/loginservice_pb";
import Container from "../Container";
import GrpcContextProvider from "../GrpcContextProvider";


export interface Props {
    children?: React.ReactNode
}

export interface State {
    isOpen: boolean
    usePortal: boolean
    loggedIn : boolean
    showPassword : boolean
}

export interface GrcpContextData {
    appInstanceId : string
    serviceUrl : string,
    grpcMetaData: Metadata
}

export default class Login extends React.Component<Props, State> {

    static grpcContext : GrcpContextData
    static username: string
    static desk: string

  
    appInstanceId: string

    username: string
    password: string
    serverUrl : string

    loginServiceClient : LoginServiceClient

    constructor(props: Props) {
        super(props)


        this.serverUrl = window.location.href

        if (this.serverUrl.endsWith("/")) {
            this.serverUrl = this.serverUrl.substr(0, this.serverUrl.length-1)
        }

        console.log("Connecting to services at:" + this.serverUrl)

        this.loginServiceClient = new LoginServiceClient(this.serverUrl, null, null)

        this.appInstanceId = v4();

        this.username = ""
        this.password = ""

        Login.username = "bert"
        Login.desk = "Delta1"

        
        this.state = {
            isOpen: true,
            usePortal: false,
            loggedIn : false,
            showPassword: false
        }

        this.onPasswordChange = this.onPasswordChange.bind(this);
        this.onUserNameChange = this.onUserNameChange.bind(this);
        this.onLogin = this.onLogin.bind(this);
    }

      onUserNameChange(e:any) {
        if( e.target && e.target.value) {
            this.username = e.target.value;
          }
      }


      onPasswordChange(e:any) {
        if( e.target && e.target.value) {
            this.password = e.target.value;
          }
      }
    
      onLogin() {


        let params = new LoginParams()
        params.setUser(this.username)
        params.setPassword(this.password)
        this.loginServiceClient.login(params, {}, (err: Error,
            response: Token) =>{

            if( err )     {
                window.alert("Failed to login: "+ err.message)
            } else {
                Login.desk = response.getDesk()
                Login.username = this.username
                Login.grpcContext = {
                    serviceUrl : this.serverUrl, 
                    grpcMetaData : {"user-name": this.username, "app-instance-id": this.appInstanceId, "auth-token" : response.getToken()},
                    appInstanceId : this.username + "@" + this.appInstanceId
                }
                this.setState({loggedIn:true})
            }

        })

      }


    render() {

        const lockButton = (
            <Tooltip content={`${this.state.showPassword ? "Hide" : "Show"} Password`} >
                <Button
                    
                    icon={this.state.showPassword ? "unlock" : "lock"}
                    intent={Intent.WARNING}
                    minimal={true}
                    onClick={this.handleLockClick}
                />
            </Tooltip>
        );


        if( this.state.loggedIn ) {
            
            return (
                <GrpcContextProvider serviceUrl={this.serverUrl} username={Login.username} appInstanceId={this.appInstanceId} >
                    <Container ></Container>
                </GrpcContextProvider>
            )
        } else {
            return (

               
                

                     <Dialog isCloseButtonShown={false}
                title="Open Trading Platform" 
                {...this.state}
                className="bp3-dark">
                <div className={Classes.DIALOG_BODY} >
                <ReactLogo  />
                    <div>
                  <InputGroup style={{marginBottom: 30}}  placeholder="Username..." onChange={this.onUserNameChange} />
                  </div>
                  <div>
              
                  <InputGroup 
                    placeholder="Password..."
                    rightElement={lockButton}
                    type={this.state.showPassword ? "text" : "password"}
                    onChange={this.onPasswordChange}
                />
                </div>

                </div>
                <div className={Classes.DIALOG_FOOTER}>
                    <div className={Classes.DIALOG_FOOTER_ACTIONS}>
                    <AnchorButton onClick={this.onLogin}
                            intent={Intent.PRIMARY}>Login
                        </AnchorButton>
                        
                    </div>
                </div>


            </Dialog>

                    
                
            )
        }

        
    }

    private handleLockClick = () => this.setState({ showPassword: !this.state.showPassword });
}
