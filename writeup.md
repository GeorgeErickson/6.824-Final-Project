#Problem Statement
A collaborative editor is a system that allows multiple users to simultaneously modify a shared environment. When numerous users are simultaneously editing a shared environment without any sort of concurrency control, conflicts are bound to arise. We implemented a browser based collaborative text editor that utilizes operational transform to optimistically control concurrent edits.

#Design
## User Experience
A collaborative application must perform similarly to its single user counterpart or its users will be frustrated. In order to behave like a single user application a collaborative application must meet the following constraints:

1. *Latency Compensation* - A user should instantly see their modifications as they are typing. This means that local changes must be displayed before server validation.
2. *Unconstrained Interaction* - A user should be able to modify any section of a document at any time. Two users should be able to edit the same line simultaneously without conflict.
3. *Realtime Propagation* - A user should see other user’s edits as they occur (bounded by network latency).

## Concurrency Control
A concurrency control mechanism is required to allow multiple users to collaboratively edit a document. Several such methods exist, and they can be broadly classified into *pessimistic* (e.g. locking) and *optimistic* (e.g. operational transformation, differential synchronization) concurrency control algorithms. Pessimistic mechanism’s reliance on locking, and the blocking network requests that locking entails would not allow our user experience constraints (enumerated in the previous section) to be met. As such we choose to create an application based on Operational Transform (an optimistic concurrency control scheme used in Google Docs, Etherpad, and Apache Wave).

Operational Transform encodes a document as a set of operations. An *operation* is an action that can be applied to a document (e.g. insert/delete/cut/paste/retain). *Transformation* is what allows two users to modify the same section of a document by resolving conflicts to ensure consistency.

#Implementation
Our system consists of a Go server that communicates with browser based clients. Operational Transform code is implemented in Coffeescript (compiles to javascript for the browser) as well as Go for the server.
## Operational Transform
A text editor only allows insert, delete and retain operations, so we needed to implement 12 transformation functions to account for every type of conflict that might occur.

## Communication 
The browser clients communicate with the server using web-sockets, which allow for fast realtime synchronization. The server maintains a set of 'hubs' to manage incoming communications. There are two main kinds of hubs - DocumentHubs and ChatHubs. A ChatHub simply repeats every message it recieves to all subscribers. This performs the chat functionality on the site. DocumentHubs handle the incoming changes to the document. The client sends Operations, which consist of many components, to the corresponding DocumentHub. The DocumentHub then interprets the operation, checks which version it refers to, and either updates it with more recent operations or applies the corresponding transforms to the document. The server then sends the most recent operations that that client hasn't seen back to each client along with the updated version number. The DocumentHub manages a queue of incoming operations and always applies them in the order in which it receives them. If a client has been disconected for awhile or somehow falls behind in version numbers the server will either reject their ops if they are hopelessly behind (hopelessly == too much work) or update/transform their operations to make them applicable.

#Difficulties
1. **Cursor preservation** - When a clients text area is updated with modifications from the server, their cursor location is lost. It is not simple to determine where their cursor should now be located. If it were to be positioned at the same index as before, inserts could cause unintuitive results. 
2. **Web-socket Reliability** - The Web-socket spec provides no guarantee of the reliability of the data channel, this means that message queuing and handshakes had to be implemented on the application level.