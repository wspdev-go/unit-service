/*
Copyright 2026 WspDev-Go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package dto

type SctpConn struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	IsEnable    bool   `json:"is_enable"`
	Description string `json:"description"`
}
