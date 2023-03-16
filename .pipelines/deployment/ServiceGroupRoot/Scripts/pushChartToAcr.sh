#!/bin/bash
set -e

# Note - This script used in the pipeline as inline script

echo "Publishing helm chart from build ACR to release ACR"

if [ -z $BUILD_TAG ]; then
  echo "-e error value of BUILD_TAG variable shouldnt be empty. Check release variables"
  exit 1
fi

echo "BUILD_TAG = ${BUILD_TAG}"

if [ -z $HELM_CHART_NAME ]; then
  echo "-e error value of HELM_CHART_NAME variable shouldnt be empty. Check release variables"
  exit 1
fi

echo "HELM_CHART_NAME = ${HELM_CHART_NAME}"

if [ -z $HELM_SEMVER ]; then
  echo "-e error value of HELM_SEMVER variable shouldnt be empty. Check release variables"
  exit 1
fi

echo "HELM_SEMVER = ${HELM_SEMVER}"

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


if [ -z $DESTINATION_CHART_REPO_NAME ]; then
  echo "-e error value of DESTINATION_CHART_REPO_NAME shouldn't be empty. check release variables"
  exit 1
fi

echo "DESTINATION_CHART_REPO_NAME = ${DESTINATION_CHART_REPO_NAME}"

echo "Done checking that all necessary variables exist."

echo "Building helm chart from image"
helm package ./aad-auth-proxy/

#Login to az cli and authenticate to acr
echo "Login cli using managed identity"
az login --identity
if [ $? -eq 0 ]; then
  echo "Logged in successfully"
else
  echo "-e error failed to login to az with managed identity credentials"
  exit 1
fi

ACCESS_TOKEN=$(az acr login --name ${DESTINATION_ACR_NAME} --expose-token --output tsv --query accessToken)
if [ $? -ne 0 ]; then 
   echo "-e error az acr login failed. Please review the Ev2 pipeline logs for more details on the error."
   exit 1
fi

echo "login to acr:${DESTINATION_ACR_NAME} using helm ..."
echo $ACCESS_TOKEN | helm registry login ${DESTINATION_ACR_NAME} -u 00000000-0000-0000-0000-000000000000 --password-stdin
if [ $? -eq 0 ]; then
  echo "login to acr:${DESTINATION_ACR_NAME} using helm completed successfully."
else
  echo "-e error login to acr:${DESTINATION_ACR_NAME} using helm failed."
  exit 1
fi 

echo "Pushing ${HELM_CHART_NAME}-${HELM_SEMVER}.tgz to oci://${DESTINATION_ACR_NAME}${DESTINATION_CHART_REPO_NAME}"
helm push ${HELM_CHART_NAME}-${HELM_SEMVER}.tgz oci://${DESTINATION_ACR_NAME}${DESTINATION_CHART_REPO_NAME}
if [ $? -eq 0 ]; then
  echo "pushing the chart to acr path: ${DESTINATION_ACR_NAME}${DESTINATION_CHART_REPO_NAME} completed successfully."
else
  echo "-e error pushing the chart to acr path:${DESTINATION_ACR_NAME}${DESTINATION_CHART_REPO_NAME} failed."
  exit 1
fi