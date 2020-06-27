# \DefaultApi

All URIs are relative to *http://api.alpha.parvaeres.io/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeploymentDeploymentIdGet**](DefaultApi.md#DeploymentDeploymentIdGet) | **Get** /deployment/{deploymentId} | Get the deployment with id deploymentId
[**DeploymentPost**](DefaultApi.md#DeploymentPost) | **Post** /deployment | Create a new deployment



## DeploymentDeploymentIdGet

> Deployment DeploymentDeploymentIdGet(ctx, deploymentId).Execute()

Get the deployment with id deploymentId

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    deploymentId := "deploymentId_example" // string | the id of the deployment to retrieve

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.DeploymentDeploymentIdGet(context.Background(), deploymentId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.DeploymentDeploymentIdGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeploymentDeploymentIdGet`: Deployment
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.DeploymentDeploymentIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**deploymentId** | **string** | the id of the deployment to retrieve | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeploymentDeploymentIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Deployment**](Deployment.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeploymentPost

> DeploymentAcceptedResponse DeploymentPost(ctx).Repository(repository).Email(email).Execute()

Create a new deployment

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    repository := "repository_example" // string |  (optional)
    email := "email_example" // string |  (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.DeploymentPost(context.Background(), ).Repository(repository).Email(email).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.DeploymentPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeploymentPost`: DeploymentAcceptedResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.DeploymentPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDeploymentPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **repository** | **string** |  | 
 **email** | **string** |  | 

### Return type

[**DeploymentAcceptedResponse**](DeploymentAcceptedResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/x-www-form-urlencoded
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

