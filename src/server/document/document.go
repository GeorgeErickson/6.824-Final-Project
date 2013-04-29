package document

import "time"

// ShareJS is a key-value store mapping a document's name to a versioned document. Documents must be created before they can be used.
// Each document has a type, a version, a snapshot and metadata.
// The document's type specifies a set of functions for interpreting and manipulating operations. 
// You must set a document's type when it is created. After a document's type is set, it cannot be changed.
// Each type has a unique string name -- for example, the plain text type is called 'text' and the JSON type is called 'json'.
// More types will be added in time - in particular rich text. For details of how types work, see Types.
// Document versions start at 0. The version is incremented each time an operation is applied.

// The document's snapshot is the contents of the document at some particular version. 
// For text documents, the snapshot is a string containing the document's contents. For JSON documents, the snapshot is a JSON object.
// Document metadata is a JSON object containing some extra information about the document, including a list of active sessions
// and a list of contributors. See Document Metadata for details.
// For example, a text document might have the snapshot 'abc' at version 100. An operation [{i:'d', p:3}] 
// is applied at version 100 which inserts a 'd' at the end of the document. After this operation is applied, the
// document has a snapshot at version 101 of 'abcd'.

// When a document is first created, it has:
// Version is set to 0
// Type set to the type specified in the create request
// Snapshot set to the result of a call to type.initialVersion(). For text, this is an empty string ''. For JSON, this is null.

type Document struct {
	//Type
	Type string
	
	Name string

	//Version
	Version int

	//Snapshot
	Snapshot string

	//metadata
	Metadata Metadata

}

//constructor for instances
func NewDoc(name string) Document {
  doc := Document{}
  doc.Type = "Text"
  doc.Name = name
  doc.Version = 0
  doc.Snapshot = ""
  doc.Metadata = Metadata{Creator:"", Ctime:time.Now(), Mtime:time.Now(), Sessions: map[string]Session{}}
  return doc
}

// Version: A version number counting from 0
// Snapshot: The current document contents
// Meta: (NEW) An object containing misc data about the document:
	// Arbitrary user data, set when an object is created. Editable by the server only.
	// creator: Name of the user agent which created the document, if the user agent has agent.name set.
	// ctime: Time the document was created
	// mtime: Time the last operation was applied. This is updated automatically on each client when it sees a document operation.
	// (**Removed**) contributors: List of users which have ever edited the document.
	// sessions: An object with an entry for every agent that is currently connected. Map agent.sessionId to:
		// name (optional, filled in by agent.name)
		// Cursor position(s) (type dependant - for text it'll be a number, for JSON it'll be a path)
		// Connection time maybe?
	// Any other application-specific data. This can be filled in by the auth function when a client connects. (And maybe clients should be able to edit this as well?)

type Metadata struct {
	Creator string
	Ctime time.Time
	Mtime time.Time
	Sessions map[string]Session
}

type Session struct {
	Name string
	Cursor int
}


