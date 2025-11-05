// Copyright 2025 SSE Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/Smicry/sse"
)

var (
	upgrader = sse.Upgrader{}
	addr     = flag.String("addr", "localhost:8080", "http service address")
)

func echo(w http.ResponseWriter, r *http.Request) {
	f, err := upgrader.Upgrade(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("upgrade:", err)
		return
	}
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sse.WriteEvent(w, sse.Event{Data: []byte(fmt.Sprintf("Hi %s, ", string(bs)))})
	f.Flush()
	time.Sleep(1 * time.Second)
	sse.WriteEvent(w, sse.Event{Data: []byte("Happy ")})
	f.Flush()
	time.Sleep(1 * time.Second)
	sse.WriteEvent(w, sse.Event{Data: []byte("New ")})
	f.Flush()
	time.Sleep(1 * time.Second)
	sse.WriteEvent(w, sse.Event{Data: []byte("Year!")})
	f.Flush()
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "http://"+*addr+"/echo")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");

    document.getElementById("send").onclick = function(evt) {
		output.innerHTML = "";
        fetch("{{.}}",{
            method: "POST",
            headers: {
                "Content-Type": "text/plain"
            },
            body: input.value
        })
        .then(response => {
            if (!response.ok) {
                throw new Error("send failed: " + response.status);
            }
            return response.body;
        })
        .then(body => {
			const reader = body.getReader();
			const decoder = new TextDecoder('utf-8');
			function read() {
				return reader.read().then(({ done, value }) => {
        			if (done) {
          				console.log('over');
          				return;
        			}
					output.innerHTML += (decoder.decode(value)).replace("data:","")
			        read();
    			});
    		}
    		read();
        })
        .catch(error => {
            console.log("fetch error: " + error);
        });
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Send" to send a message to the server. 
You can change the message and send multiple times.
<p>
<form>
<p><input id="input" type="text" value="Nerd">
<button id="send">Send</button>
</form>
</td></tr>
<tr><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr>
</table>
</body>
</html>`))
