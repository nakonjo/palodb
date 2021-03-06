package lang

import "fmt"

type Callback func(Statement, *Token) (Statement, error)

func passthrough(stmt Statement, token *Token) (Statement, error) {
  return stmt, nil
}

func createDimension(stmt Statement, token *Token) (Statement, error) {
  // fmt.Printf(">>> createDimension(%s)\n", token.Literal)
  return &CreateDimension{token.Literal, []*CreateLevel{}}, nil
}

func createFirstLevel(stmt Statement, token *Token) (Statement, error) {
  // fmt.Printf(">>> createFirstLevel(%s)\n", token.Literal)
  dim := stmt.(*CreateDimension)
  lvl := &CreateLevel{dim, token.Literal, []*CreateAttribute{}}
  dim.Levels = append(dim.Levels, lvl)
  return lvl, nil
}

func createNextLevel(stmt Statement, token *Token) (Statement, error) {
  // fmt.Printf(">>> createNextLevel(%s)\n", token.Literal)
  oth := stmt.(*CreateLevel)

  for _, tst := range oth.Dim.Levels {
    if tst.Name == token.Literal {
      return nil, &SyntaxError{
        token.Position,
        fmt.Sprintf("Redefinition of level '%s'", token.Literal),
      }
    }
  }

  lvl := &CreateLevel{oth.Dim, token.Literal, []*CreateAttribute{}}
  oth.Dim.Levels = append(oth.Dim.Levels, lvl)
  return lvl, nil
}

func createDimensionLevelsEnd(stmt Statement, token *Token) (Statement, error) {
  // fmt.Printf(">>> createDimensionLevelsEnd\n")
  lvl := stmt.(*CreateLevel)
  return lvl.Dim, nil
}

func createFirstAttribute(stmt Statement, token *Token) (Statement, error) {
  // fmt.Printf(">>> createFirstAttribute(%s)\n", token.Literal)
  lvl := stmt.(*CreateLevel)
  att := &CreateAttribute{lvl, token.Literal, "", false}
  lvl.Attributes = append(lvl.Attributes, att)
  return att, nil
}

func createNextAttribute(stmt Statement, token *Token) (Statement, error) {
  // fmt.Printf(">>> createNextAttribute(%s)\n", token.Literal)
  oth := stmt.(*CreateAttribute)

  for _, tst := range oth.Lvl.Attributes {
    if tst.Name == token.Literal {
      return nil, &SyntaxError{
        token.Position,
        fmt.Sprintf("Redefinition of attribute '%s'", token.Literal),
      }
    }
  }

  att := &CreateAttribute{oth.Lvl, token.Literal, "", false}
  oth.Lvl.Attributes = append(oth.Lvl.Attributes, att)
  return att, nil
}

func createAttributeSetType(stmt Statement, token *Token) (Statement, error) {
  // fmt.Printf(">>> createAttributeSetType(%s)\n", token.Literal)
  att := stmt.(*CreateAttribute)
  att.Dtype = token.Literal
  return att, nil
}

func createAttributeSetDefault(stmt Statement, token *Token) (Statement, error) {
  // fmt.Printf(">>> createAttributeSetDefault\n")
  att := stmt.(*CreateAttribute)
  att.Default = true
  return att, nil
}

func createLevelAttributesEnd(stmt Statement, token *Token) (Statement, error) {
  // fmt.Printf(">>> createLevelAttributesEnd\n")
  att := stmt.(*CreateAttribute)
  return att.Lvl, nil
}
