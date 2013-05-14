class OperationalTransform
  constructor: (@options) ->
    _.extend @, Backbone.Events
    @editor = @options.editor
  
  setModel: (model) ->
    @supress = false
    #unbind model events
    if @model
      @model.off null, @onmessage

    @model = model
    model.on 'message', @onmessage
    
    #local events
    @off()
    @on 'insert', @oninsert
    @on 'delete', @ondelete



  onmessage: (data) =>
    @suppress = true
    console.log data
    # if data.count is 0
    #   @editor.doc.setValue data.Snapshot

    @supress = false
    
  
  oninsert: (data) =>
    if @supress
      return
    op =
      Insert: data.text
      Position: data.position
      Delete: ""
    attributes = _.clone @model.attributes
    attributes.OpData = [[op]]
    @model.send attributes


  ondelete: (data) ->
    if @supress
      return
    op =
      Insert: ""
      Position: data.position
      Delete: data.text
    attributes = _.clone @model.attributes
    attributes.OpData = [[op]]
    @model.send attributes

module.exports = OperationalTransform
