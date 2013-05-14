6.824-Final-Project
===================

##Setup - Client
- ```cd src/server/client``` - ☹
- ```npm install -g brunch```
- ```gem install sass```
- ```brunch b``` or ```brunch w```


##URLS
- 	 root right now: "127.0.0.1:8000"
-    "/"
-    "/ws" - websocket connection for home page - lists document names/updates
-    "/rest/documents" - list of document names ["doc1","doc2",...]
-    "/documents/(.*)" - websocket for document ops for a certain document /* = document name
-    "/chat/(.*)" - websocket for chat for a document /* = document name    