/*
 * Parvaeres API
 *
 * Parvaeres magic deployment API
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package parvaeres

type GetDeploymentLogsResponse struct {
	Message string `json:"Message,omitempty"`

	Error bool `json:"Error,omitempty"`

	Items []Logs `json:"Items,omitempty"`
}
