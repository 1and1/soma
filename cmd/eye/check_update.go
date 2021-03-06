/*
Copyright (c) 2016, Jörg Pernfuß <code.jpe@gmail.com>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	"database/sql"
	"fmt"

	"github.com/1and1/soma/lib/proto"
)

func CheckUpdateOrInsertOrDelete(details *proto.Deployment) error {
	var (
		err              error
		itemID, lookupID string
		item             *ConfigurationItem
	)

	if lookupID, item, err = Itemize(details); err != nil {
		return err
	}

	fmt.Println(lookupID)
	fmt.Println(item)

	err = Eye.run.check_item.QueryRow(item.ConfigurationItemId).Scan(&itemID)
	switch details.Task {
	case "rollout":
		if err == sql.ErrNoRows {
			return addItem(item, lookupID)
		} else if err != nil {
			return err
		}
	case "deprovision":
		if err == sql.ErrNoRows {
			// nothing to do
			return nil
		} else if err != nil {
			return err
		}
	}

	if item.ConfigurationItemId.String() != itemID {
		panic(`Database corrupted`)
	}
	switch details.Task {
	case "rollout":
		return updateItem(item, lookupID)
	case "deprovision":
		return deleteItem(itemID)
	default:
		return fmt.Errorf(`Unknown Task requested`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
