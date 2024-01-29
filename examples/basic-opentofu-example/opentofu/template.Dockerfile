# syntax=docker/dockerfile-upstream:1.4.0
# This is a template Dockerfile for the bundle's invocation image
# You can customize it to use different base images, install tools and copy configuration files.
#
# Porter will use it as a template and append lines to it for the mixins
# and to set the CMD appropriately for the CNAB specification.
#
# Add the following line to porter.yaml to instruct Porter to use this template
# dockerfile: template.Dockerfile

# You can control where the mixin's Dockerfile lines are inserted into this file by moving the "# PORTER_*" tokens
# another location in this file. If you remove a token, its content is appended to the end of the Dockerfile.
FROM --platform=linux/amd64 debian:stable-20230703-slim
RUN apt-get update \
    && apt-get install --no-install-recommends -y bash curl jq sshpass
COPY --link .cnab /cnab
# COPY terraform/components/terraform-provider-f5os /usr/local/share/terraform/plugins/terraform.local/local/f5os/1.3.2/linux_amd64/terraform-provider-f5os