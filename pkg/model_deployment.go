/*
 * Parvaeres API
 *
 * Parvaeres magic deployment API
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package parvaeres

type Deployment struct {
	Status string `json:"status,omitempty"`

	LiveUrls []string `json:"liveUrls,omitempty"`
}