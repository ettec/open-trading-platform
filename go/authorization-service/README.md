# authorization-service

This service checks user permissions prior to allowing an action to proceed, for example a user must have the Trader flag to create orders.  It also contains the hooks to attach an environment specific authentication mechanism, out of the box passwords/tokens are not checked.