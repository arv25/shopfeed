### Build Container
docker build --rm --tag=ebs-docker-shopfeed .

### Build Deploy Artifact
zip -r ebs-docker-shopfeed.zip . --exclude=*pkg* --exclude=*src* --exclude=*.git* --exclude=*swp* --exclude=*.reql* --exclude=*.md*
