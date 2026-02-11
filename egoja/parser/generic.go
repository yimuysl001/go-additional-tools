package eparser

import "strings"

type GenericTokenParser struct {
	open    string
	close   string
	skipStr bool
}

func NewGenericTokenParser(open, close string, skipStr bool) *GenericTokenParser {
	return &GenericTokenParser{
		open:    open,
		close:   close,
		skipStr: skipStr,
	}
}

func (gtp *GenericTokenParser) Parse(source string, handler func(string) string) string {
	stream := NewCharacterStream(source)
	var builder strings.Builder

	for stream.hasMore() {
		builder.WriteString(gtp.parseStream(stream, "", handler))
	}

	return strings.TrimSpace(builder.String())
}

func (gtp *GenericTokenParser) parseStream(stream *CharacterStream, expect string, handler func(string) string) string {
	var builder strings.Builder

	for stream.hasMore() {
		if expect != "" && stream.match(expect, true) {
			return strings.TrimSpace(builder.String())
		}

		if stream.match(gtp.open, true) {
			value := handler(gtp.parseStream(stream, gtp.close, handler))
			if value != "" {
				builder.WriteString(value)
			}
		} else {
			ch := stream.consume()
			builder.WriteRune(ch)
			if gtp.skipStr && ch == '\'' {
				builder.WriteString(gtp.consumeUntil(stream, "'"))
			} else if gtp.skipStr && ch == '"' {
				builder.WriteString(gtp.consumeUntil(stream, "\""))
			} else if ch == '{' {
				builder.WriteString(gtp.parseStream(stream, "}", handler))
				builder.WriteString("}")
			}
		}
	}

	return strings.TrimSpace(builder.String())

}

func (gtp *GenericTokenParser) consumeUntil(stream *CharacterStream, str string) string {
	start := stream.getPosition()

	for stream.hasMore() {
		if stream.match("\\", true) {
			stream.consume()
		}

		if stream.match(str, true) {
			break
		}

		stream.consume()
	}

	return stream.substring(start, stream.getPosition())
}
