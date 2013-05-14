ot = require 'lib/ot'

class Editor extends Backbone.View
  className: 'hide'
  initialize: ->
    @ot = new ot
      editor: @

  render: ->
    @$el.append '<div id="editor"></div>'
    @editor = ace.edit("editor")
    @doc = @editor.getSession().getDocument()
    @doc.setNewLineMode 'unix'

  getPosition: (range) ->
    lines = @doc.getLines 0, range.start.row
      
    offset = 0

    for line, i in lines
      offset += if i < range.start.row
        line.length
      else
        range.start.column

    # Add the row number to include newlines.
    offset + range.start.row

  detachEvents: =>
    @editor.removeAllListeners 'change'

  attachEvents: =>
    @detachEvents()
    @editor.on 'change', @onchange

  show: (model) ->
    @detachEvents()
    @ot.setModel model


    @attachEvents()
    @$el.removeClass 'hide'

  onchange: (e, editor) =>
    data = e.data

    # For insertLines and removeLines
    if data.lines
      data.text = data.lines.join('\n') + '\n'

    data.position = @getPosition data.range

    switch data.action
      when "insertLines", "insertText" then @ot.trigger 'insert', data
      when "removeLines", "removeText" then @ot.trigger 'remove', data

  hide: ->
    @$el.addClass 'hide'


editor = new Editor()
editor.$el.appendTo 'body'
editor.render()
module.exports = editor

