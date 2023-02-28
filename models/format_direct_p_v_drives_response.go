// Code generated by go-swagger; DO NOT EDIT.

// This file is part of MinIO Operator
// Copyright (c) 2023 MinIO, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// FormatDirectPVDrivesResponse format direct p v drives response
//
// swagger:model formatDirectPVDrivesResponse
type FormatDirectPVDrivesResponse struct {

	// format issues list
	FormatIssuesList []*PvFormatErrorResponse `json:"formatIssuesList"`
}

// Validate validates this format direct p v drives response
func (m *FormatDirectPVDrivesResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateFormatIssuesList(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *FormatDirectPVDrivesResponse) validateFormatIssuesList(formats strfmt.Registry) error {
	if swag.IsZero(m.FormatIssuesList) { // not required
		return nil
	}

	for i := 0; i < len(m.FormatIssuesList); i++ {
		if swag.IsZero(m.FormatIssuesList[i]) { // not required
			continue
		}

		if m.FormatIssuesList[i] != nil {
			if err := m.FormatIssuesList[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("formatIssuesList" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("formatIssuesList" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this format direct p v drives response based on the context it is used
func (m *FormatDirectPVDrivesResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateFormatIssuesList(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *FormatDirectPVDrivesResponse) contextValidateFormatIssuesList(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.FormatIssuesList); i++ {

		if m.FormatIssuesList[i] != nil {
			if err := m.FormatIssuesList[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("formatIssuesList" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("formatIssuesList" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *FormatDirectPVDrivesResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *FormatDirectPVDrivesResponse) UnmarshalBinary(b []byte) error {
	var res FormatDirectPVDrivesResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
