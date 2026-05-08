/*
Copyright 2026 WspDev-Go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package dto

type M3UaAspLink struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	SctpConnID   int    `json:"sctp_conn_id"`
	M3UaAsConnID int    `json:"m3ua_as_conn_id"`
	IsEnable     bool   `json:"is_enable"`
	Description  string `json:"description"`
}
