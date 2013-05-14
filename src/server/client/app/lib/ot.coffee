class OperationalTransform
  constructor: (@options) ->
    _.extend @, Backbone.Events
    @editor = @options.editor
  
  setModel: (model) ->
    #unbind model events
    if @model
      @model.off null, @onmessage

    @model = model
    model.on 'message', @onmessage
    
    #local events
    @off()
    @on 'insert', @oninsert
    @on 'delete', @ondelete



  onmessage: (data) ->
    
  
  oninsert: (data) =>
    op =
      Insert: data.text
      Position: data.position
      Delete: ""
    attributes = _.clone @model.attributes
    attributes.OpData = [[op]]
    console.log @model.send attributes


  ondelete: (data) ->
    console.log data

module.exports = OperationalTransform
