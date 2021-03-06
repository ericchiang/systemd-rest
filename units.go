/*
*  Copyright 2013 CoreOS, Inc
*  Copyright 2013 Docker Authors
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
 */

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/philips/go-systemd"
	"net/http"
)

func listHandler(w http.ResponseWriter, r *http.Request) {
	var (
		out interface{}
	)

	s := new(systemd.Systemd1)
	err := s.Connect()
	if err != nil {
		// TODO: Return 40* code
		fmt.Fprint(w, err)
	}

	out, err = s.ListUnits()

	if err != nil {
		// TODO: Return 40* code
		fmt.Fprint(w, err)
	}

	outJson, _ := json.Marshal(out)
	fmt.Println("%s\n", outJson)
	fmt.Fprint(w, "%s\n", outJson)
}

func unitHandler(w http.ResponseWriter, r *http.Request) {
	var (
		out interface{}
	)

	vars := mux.Vars(r)
	s := new(systemd.Systemd1)
	err := s.Connect()
	if err != nil {
		// TODO: Return 40* code
		w.WriteHeader(404)
		fmt.Fprint(w, err)
	}

	switch vars["method"] {
	case "start":
		out, err = s.StartUnit(vars["unit"], vars["mode"])
	case "stop":
		out, err = s.StopUnit(vars["unit"], vars["mode"])
	}

	if err != nil {
		w.WriteHeader(404)
	}

	outJson, _ := json.Marshal(out)
	fmt.Fprintf(w, "%s\n", outJson)
}

func setupUnits(r *mux.Router, o Options) {
	r.HandleFunc("", listHandler)
	r.HandleFunc("/", listHandler)
	r.HandleFunc("/{unit}/{method}/{mode}", unitHandler)

	return
}
