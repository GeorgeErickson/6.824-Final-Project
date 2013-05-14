DocumentListItemView = require 'views/DocumentListItem'

describe 'DocumentListItemView', ->
    beforeEach ->
        @view = new DocumentListItemView()

    it 'should exist', ->
        expect(@view).to.be.ok
