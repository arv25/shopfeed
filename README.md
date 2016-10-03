### Build Container
docker build --rm --tag=ebs-docker-shopfeed .

### Build Deploy Artifact
zip -r ebs-docker-shopfeed.zip . --exclude=*pkg* --exclude=*src* --exclude=*.git* --exclude=*swp* --exclude=*.reql* --exclude=*.md*

###
N.B.- The embedded configuration file in the .ebextensions folder will allow nginx to keep sockets open for longer than default 60 seconds.

### Update Load Balancer Idle Timeout
Follow instructions [here](https://cloudavail.com/2015/10/18/allowing-long-idle-timeouts-when-using-aws-elasticbeanstalk-and-docker/) to update __Idle Timeout__ value in ELB configuration to be longer than the default 60 seconds.
