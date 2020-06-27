# Deployment

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | Pointer to **string** |  | [optional] 
**LiveUrls** | Pointer to **[]string** |  | [optional] 

## Methods

### NewDeployment

`func NewDeployment() *Deployment`

NewDeployment instantiates a new Deployment object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeploymentWithDefaults

`func NewDeploymentWithDefaults() *Deployment`

NewDeploymentWithDefaults instantiates a new Deployment object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *Deployment) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Deployment) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Deployment) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *Deployment) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetLiveUrls

`func (o *Deployment) GetLiveUrls() []string`

GetLiveUrls returns the LiveUrls field if non-nil, zero value otherwise.

### GetLiveUrlsOk

`func (o *Deployment) GetLiveUrlsOk() (*[]string, bool)`

GetLiveUrlsOk returns a tuple with the LiveUrls field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLiveUrls

`func (o *Deployment) SetLiveUrls(v []string)`

SetLiveUrls sets LiveUrls field to given value.

### HasLiveUrls

`func (o *Deployment) HasLiveUrls() bool`

HasLiveUrls returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


