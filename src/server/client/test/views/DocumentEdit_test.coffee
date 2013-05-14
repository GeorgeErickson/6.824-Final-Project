DocumentEditView = require 'views/DocumentEdit'

describe 'DocumentEditView', ->
    beforeEach ->
        @view = new DocumentEditView()

    it 'should exist', ->
        expect(@view).to.be.ok
