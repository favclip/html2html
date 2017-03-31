package html2html

func TagReplacer(tag Tag, f func(tag Tag) (Token, error)) (Token, error) {
	if altToken, err := f(tag); err != nil {
		return nil, err
	} else if altToken != nil {
		return altToken, nil
	}

	for _, token := range tag.Tokens() {
		if token.Type() != TypeTagToken {
			continue
		}

		childTag := token.Tag()

		altToken, err := TagReplacer(childTag, f)
		if err != nil {
			return nil, err
		} else if altToken != nil {
			tag.ReplateChildToken(childTag, altToken)
		}
	}

	return nil, nil
}
