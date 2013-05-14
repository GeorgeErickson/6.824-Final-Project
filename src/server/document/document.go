package document

import "fmt"

// ShareJS is a key-value store mapping a document's name to a versioned document. Documents must be created before they can be used.
// Each document has a type, a version, a Snapshot and metadata.
// The document's type specifies a set of functions for interpreting and manipulating operations. 
// You must set a document's type when it is created. After a document's type is set, it cannot be changed.
// Each type has a unique string name -- for example, the plain text type is called 'text' and the JSON type is called 'json'.
// More types will be added in time - in particular rich text. For details of how types work, see Types.
// Document versions start at 0. The version is incremented each time an operation is applied.

// The document's Snapshot is the contents of the document at some particular version. 
// For text documents, the Snapshot is a string containing the document's contents. For JSON documents, the Snapshot is a JSON object.
// Document metadata is a JSON object containing some extra information about the document, including a list of active sessions
// and a list of contributors. See Document Metadata for details.
// For example, a text document might have the Snapshot 'abc' at version 100. An operation [{i:'d', p:3}] 
// is applied at version 100 which inserts a 'd' at the end of the document. After this operation is applied, the
// document has a Snapshot at version 101 of 'abcd'.

// When a document is first created, it has:
// Version is set to 0
// Type set to the type specified in the create request
// Snapshot set to the result of a call to type.initialVersion(). For text, this is an empty string ''. For JSON, this is null.

type Document struct {
	//Type
	Type string
	
	Name string

	Title string

	//Version
	Version int

	//clientId
	ClientId string

	//Snapshot
	Snapshot string

	//opdata for deserializing
	//version info is important
	OpData []TextOp

}

//constructor for instances
func NewDoc(name string) Document {
  doc := Document{}
  doc.Type = "Text"
  doc.Name = name
  doc.Title = ""
  doc.Version = 0
  doc.Snapshot = ""
  doc.OpData = []TextOp{}
  return doc
}
/*
checkValidComponent = (c) ->
  throw new Error 'component missing position field' if typeof c.p != 'number'

  i_type = typeof c.i
  d_type = typeof c.d
  throw new Error 'component needs an i or d field' unless (i_type == 'string') ^ (d_type == 'string')

  throw new Error 'position cannot be negative' unless c.p >= 0

checkValidOp = (op) ->
  checkValidComponent(c) for c in op
  true

text.apply = (Snapshot, op) ->
  checkValidOp op
  for component in op
    if component.i?
      Snapshot = strInject Snapshot, component.p, component.i
    else
      deleted = Snapshot[component.p...(component.p + component.d.length)]
      throw new Error "Delete component '#{component.d}' does not match deleted text '#{deleted}'" unless component.d == deleted
      Snapshot = Snapshot[...component.p] + Snapshot[(component.p + component.d.length)..]
  
  Snapshot
*/
func (doc *Document) ApplyOp(op TextOp) bool {
	if(!doc.checkValidOp(op)){
		return false
	}
	for _, component := range op {
    	if component.Insert != "" {
    		doc.Snapshot = doc.StrInject(doc.Snapshot, component.Position, component.Insert)
    	} else{
    		if(len(doc.Snapshot) < component.Position || len(doc.Snapshot) < component.Position + len(component.Delete)){
    			fmt.Println("OP ERROR!fff")
    			return false
    		}
    		deleted := doc.Snapshot[component.Position:component.Position+len(component.Delete)]
    		if(component.Delete != deleted){
    			fmt.Println("hhhhhhhh")
    			return false
    		}
    		doc.Snapshot = doc.Snapshot[0:component.Position] + doc.Snapshot[component.Position+len(component.Delete):len(doc.Snapshot)]
    	}
	}
	return true
}

func (doc *Document) ApplyOps(op TextOp, version int) bool {
	last_op_index := doc.Version - version
	if last_op_index > 25  || doc.Version < version{
		return false
	}
	if last_op_index != 0 {
		transform_ops := doc.OpData[last_op_index-1:]
		for i := 0; i < len(transform_ops); i++ {
			top := transform_ops[i]
			op = op.transform(top)
		}
	}
	err := doc.ApplyOp(op)
	if err == false {
		return false
	}
	doc.OpData = append(doc.OpData, op)
	//don't let someone get more than 25 versions behind
	if(len(doc.OpData) > 25){
		doc.OpData = doc.OpData[1:]
	}
	doc.BumpVersion()
	return true
}

func (doc *Document) BumpVersion() {
	doc.Version++
}

func (doc *Document) checkValidOp(op TextOp) bool {
	for _, component := range op {
		if(doc.checkValidComponent(component) == false){
			return false
		}
	}
	return true
}

func (doc *Document) checkValidComponent(component Component) bool {
	if component.Position < 0 {
		return false
	}
	return true
}

func (doc *Document) StrInject(Snapshot string, position int, inserted string) string {
	return Snapshot[0:Min(position,len(Snapshot))]+inserted+Snapshot[Min(position,len(Snapshot)):len(Snapshot)]
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


