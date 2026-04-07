targetScope = 'resourceGroup'

@description('Placeholder Bicep entry point for future infrastructure phases.')
param location string = resourceGroup().location

output deploymentLocation string = location
