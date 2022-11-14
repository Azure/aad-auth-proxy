#!/bin/bash
set -e

# Note - This script used in the pipeline as inline script

echo "Publishing agent image from CDPx ACR to team's ACR"

if [ -z $AGENT_RELEASE ]; then
  echo "-e error AGENT_RELEASE shouldnt be empty. check release variables"
  exit 1
fi

echo "AGENT_RELEASE = ${AGENT_RELEASE}"

if [ -z $AGENT_IMAGE_FULL_PATH ]; then
  echo "-e error AGENT_IMAGE_FULL_PATH shouldnt be empty. check release variables"
  exit 1
fi

echo "AGENT_IMAGE_FULL_PATH = ${AGENT_IMAGE_FULL_PATH}"

if [ -z $BUILD_TAG ]; then
  echo "-e error value of BUILD_TAG shouldn't be empty. check release variables"
  exit 1
fi

echo "BUILD_TAG = ${BUILD_TAG}"

if [ -z $BUILD_ACR ]; then
  echo "-e error value of BUILD_ACR shouldn't be empty. check release variables"
  exit 1
fi

echo "BUILD_ACR = ${BUILD_ACR}"

if [ -z $BUILD_REPO_NAME ]; then
  echo "-e error value of BUILD_REPO_NAME shouldn't be empty. check release variables"
  exit 1
fi

echo "BUILD_REPO_NAME = ${BUILD_REPO_NAME}"

if [ -z $DESTINATION_ACR_NAME ]; then
  echo "-e error value of DESTINATION_ACR_NAME shouldn't be empty. check release variables"
  exit 1
fi

echo "DESTINATION_ACR_NAME = ${DESTINATION_ACR_NAME}"

#Login to az cli and authenticate to acr
echo "Login cli using managed identity"
az login --identity
if [ $? -eq 0 ]; then
  echo "Logged in successfully"
else
  echo "-e error failed to login to az with managed identity credentials"
  exit 1
fi     

echo "Pushing ${AGENT_IMAGE_FULL_PATH} to ${DESTINATION_ACR_NAME}"
az acr import --name $DESTINATION_ACR_NAME --registry $BUILD_ACR --source official/${BUILD_REPO_NAME}:${BUILD_TAG} --image $AGENT_IMAGE_FULL_PATH
if [ $? -eq 0 ]; then
  echo "Retagged and pushed image successfully"
else
  echo "-e error failed to retag and push image to destination ACR"
  exit 1
fi