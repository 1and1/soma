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
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

func DeleteConfigurationItem(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var (
		itemID string
		err    error
	)
	itemID = params.ByName("item")
	if _, err = uuid.FromString(itemID); err != nil {
		dispatchBadRequest(&w, err.Error())
		return
	}
	if err = deleteItem(itemID); err != nil {
		dispatchInternalServerError(&w, err.Error())
		return
	}
	dispatchNoContent(&w)
}

func deleteItem(itemID string) error {
	var (
		lookupID string
		count    int
		err      error
	)

	if err = Eye.run.get_lookup.QueryRow(itemID).Scan(&lookupID); err == sql.ErrNoRows {
		// not being able to delete what we do not have is ok
		return nil
	} else if err != nil {
		// either a real error, or what is to be deleted does not exist
		return err
	}

	if _, err = Eye.run.delete_item.Exec(itemID); err != nil {
		return err
	}

	if err = Eye.run.item_count.QueryRow(lookupID).Scan(&count); err != nil {
		return err
	}

	if count != 0 {
		return nil
	}
	_, err = Eye.run.delete_lookup.Exec(lookupID)
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
