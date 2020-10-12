# installation

The below instructions are for linux only and assume you are using microk8s to run your kubernetes cluster, however instructions can used to install the application on any linux Kubernetes cluster which has a default storage class.  If you are installing on a non-microk8s cluster uncomment the first 3 lines of the install.sh script before running it.

Install a fresh copy of [microk8s](https://microk8s.io/)

Enable the required microk8s plugins using the following command:

`microk8s enable dns storage helm3`

Start the cluster:

`microk8s start`

Checkout the otp source code from https://github.com/ettec/open-trading-platform

Run the installation script, from the root of the checkout:

`./install/install.sh `

That's it.  After the install script completes it will inform you of the port to use to run the OTP client.  You can login using any of the following user ids, no password is required out of the box (the authentication-service has a hook for a token/password validation plugin).

**trader1** - has trading permissions and is a member of the Desk1 trading desk (users on the same desk can see and control each others orders)

**trader2** - has trading permissions and is a member of the Desk1 trading desk

**support1** - has view only permissions on the Desk1 trading desk

**traderA** - has trading permissions and is a member of the DeskA trading desk

**traderB** - has trading permission and is  a member of the DeskA trading desk

**supportA** - has view only permission on the DeskA trading desk 





