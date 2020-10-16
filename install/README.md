# installation

Below are the instructions for installing OTP on both a standalone microk8s kubernetes cluster or a standard kubernetes cluster.  If you have no cluster already provisioned microk8s will probably be the easiest option to get started with.

### standard kubernetes

Point your kubectl install at the context of the cluster you want to install the application on.

Checkout the otp source code from https://github.com/ettec/open-trading-platform

Run the installation script, from the root of the checkout as follows:

`./install/install.sh -v 1.0.11 `

That's it.  After the install script completes it will inform you of the port to use to run the OTP client.  You can login using any of the [user ids](#userids) at the bottom of this README, no password is required out of the box (the authentication-service has a hook for a token/password validation plugin).  

### microk8s

Install a fresh copy of [microk8s](https://microk8s.io/)

Enable the required microk8s plugins using the following command:

`microk8s enable dns storage helm3`

Start the cluster:

`microk8s start`

Checkout the otp source code from https://github.com/ettec/open-trading-platform

Run the installation script, from the root of the checkout with the arguments as shown:

`./install/install.sh -v 1.0.11 -m `

That's it.  After the install script completes it will inform you of the port to use to run the OTP client.  You can login using any of the following user ids, no password is required out of the box (the authentication-service has a hook for a token/password validation plugin). 

### out of the box user ids <a name="userids"></a>

**trader1** - has trading permissions and is a member of the Desk1 trading desk (users on the same desk can see and control each others orders)

**trader2** - has trading permissions and is a member of the Desk1 trading desk

**support1** - has view only permissions on the Desk1 trading desk

**traderA** - has trading permissions and is a member of the DeskA trading desk

**traderB** - has trading permission and is  a member of the DeskA trading desk

**supportA** - has view only permission on the DeskA trading desk 





