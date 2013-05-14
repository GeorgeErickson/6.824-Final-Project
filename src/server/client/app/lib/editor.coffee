ot = require 'lib/ot'

class Editor extends Backbone.View
  className: 'hide'
  initialize: ->
    @ot = new ot
      editor: @

  render: ->
    @$el.append '<div id="editor"></div>'
    @editor = ace.edit("editor")

  detachEvents: =>
    @editor.removeAllListeners 'change'

  attachEvents: =>
    @detachEvents()
    @editor.on 'change', @onchange

  show: (model) ->
    @detachEvents()
    @ot.setModel model
    snapshot = model.get 'Snapshot'
    @editor.setValue snapshot
    @attachEvents()
    @$el.removeClass 'hide'

  onchange: (e, editor) =>
    data = e.data

    switch data.action
      when "insertLines", "insertText" then @ot.trigger 'insert', data
      when "removeLines", "removeText" then @ot.trigger 'remove', data

  hide: ->
    @$el.addClass 'hide'


editor = new Editor()
editor.$el.appendTo 'body'
editor.render()
module.exports = editor

