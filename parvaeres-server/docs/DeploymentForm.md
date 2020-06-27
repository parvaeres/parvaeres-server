# DeploymentForm

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Repository** | Pointer to **string** |  | [optional] 
**Email** | Pointer to **string** |  | [optional] 

## Methods

### NewDeploymentForm

`func NewDeploymentForm() *DeploymentForm`

NewDeploymentForm instantiates a new DeploymentForm object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeploymentFormWithDefaults

`func NewDeploymentFormWithDefaults() *DeploymentForm`

NewDeploymentFormWithDefaults instantiates a new DeploymentForm object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRepository

`func (o *DeploymentForm) GetRepository() string`

GetRepository returns the Repository field if non-nil, zero value otherwise.

### GetRepositoryOk

`func (o *DeploymentForm) GetRepositoryOk() (*string, bool)`

GetRepositoryOk returns a tuple with the Repository field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRepository

`func (o *DeploymentForm) SetRepository(v string)`

SetRepository sets Repository field to given value.

### HasRepository

`func (o *DeploymentForm) HasRepository() bool`

HasRepository returns a boolean if a field has been set.

### GetEmail

`func (o *DeploymentForm) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *DeploymentForm) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *DeploymentForm) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *DeploymentForm) HasEmail() bool`

HasEmail returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


