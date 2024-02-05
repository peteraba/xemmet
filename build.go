package main

func Build(tokens []Token, num, siblingCount int) ElemList {
	if tokens == nil {
		return nil
	}

	var (
		elemList    = ElemList{}
		newElemList ElemList
	)

	for _, token := range tokens {
		switch value := token.(type) {
		case *GroupToken:
			newElemList = BuildFromGroup(value, num, siblingCount)

			elemList = append(elemList, newElemList...)

		case *TagToken:
			newElemList = BuildFromTag(value, num, siblingCount)

			elemList = append(elemList, newElemList...)
		}
	}

	return elemList
}

func BuildFromGroup(token *GroupToken, num, siblingCount int) ElemList {
	elemList := ElemList{}

	if token.GetRepeat() == 0 {
		siblingCount = token.GetRepeat()
	}

	for i := 1; i <= token.GetRepeat(); i++ {
		if token.GetRepeat() > 1 {
			num = i
		}

		elemList = append(elemList, Build(token.Children, num, siblingCount)...)
	}

	return elemList
}

func BuildFromTag(token *TagToken, num, siblingCount int) ElemList {
	elemList := ElemList{}

	if token.Repeat > 1 {
		siblingCount = token.Repeat
	}

	for i := 1; i <= token.Repeat; i++ {
		if token.Repeat > 1 {
			num = i
		}

		elem := &Elem{
			Name:         token.Name,
			ID:           token.ID,
			Classes:      token.Classes,
			Attributes:   token.Attributes,
			Text:         token.Text,
			Children:     Build(token.Children, num, siblingCount),
			Num:          num,
			SiblingCount: siblingCount,
		}

		elemList = append(elemList, elem)
	}

	return elemList
}
