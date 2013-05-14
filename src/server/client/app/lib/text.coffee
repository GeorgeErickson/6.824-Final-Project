module.exports =
  compose: (op1, op2) ->
    cop = op1.slice()
    cop.push c for c in op2
    cop

