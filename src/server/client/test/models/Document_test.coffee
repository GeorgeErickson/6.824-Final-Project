DocumentModel = require 'models/Document'

describe 'DocumentModel', ->
    beforeEach ->
        @model = new DocumentModel()

    it 'should exist', ->
        expect(@model).to.be.ok
