class OperationalTransform
  constructor: (@options) ->
    _.extend @, Backbone.Events
    @editor = @options.editor
  
  setModel: (model) ->
    @off()
    @stopListening()
    @model = model
    
    snapshot = @model.get 'Snapshot'
    if snapshot
      @editor.editor.setValue snapshot

    
    @listenTo model, 'message', @onmessage
    @on 'insert', @oninsert
    @on 'delete', @ondelete


    
    
    



  onmessage: (data) =>
    @quiet = true
    snapshot = @model.get 'Snapshot'
    unless snapshot
      if data.Snapshot
        @editor.editor.setValue data.Snapshot
    @quiet = false
    
  
  oninsert: (data) =>
    op =
      Insert: data.text
      Position: data.position
      Delete: ""
    attributes = _.clone @model.attributes
    attributes.OpData = [[op]]
    unless @quiet
      @model.send attributes


  ondelete: (data) ->
    console.log data
    op =
      Insert: ""
      Position: data.position
      Delete: data.text
    attributes = _.clone @model.attributes
    attributes.OpData = [[op]]
    unless @quiet
      @model.send attributes

module.exports = OperationalTransform
