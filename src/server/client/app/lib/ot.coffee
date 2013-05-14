class OperationalTransform
  constructor: (@options) ->
    _.extend @, Backbone.Events
    @editor = @options.editor
  
  setModel: (model) ->
    @off()
    @on 'insert', @oninsert
    @on 'delete', @ondelete
  
  oninsert: (data) ->
    console.log data

  ondelete: (data) ->
    console.log data

module.exports = OperationalTransform
