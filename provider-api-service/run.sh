#!/bin/bash
# export OPENSTACK_DNS_NAMESERVERS="8.8.8.8"
# export OPENSTACK_FAILURE_DOMAIN="compute"
# export OPENSTACK_CONTROL_PLANE_MACHINE_FLAVOR="cluster.controller"
# export OPENSTACK_NODE_MACHINE_FLAVOR="cluster.compute"
# #export OPENSTACK_IMAGE_NAME="Ubuntu-18.04-190314"
# export OPENSTACK_IMAGE_NAME="ubuntu-k8s-1.24"
# export OPENSTACK_SSH_KEY_NAME="clusterapi"
# export OPENSTACK_EXTERNAL_NETWORK_ID="147a0c31-89e4-412a-859d-7f3cec25bb6f"
# source env.rc clouds.yaml openstack-cloud
printenv
./provider
