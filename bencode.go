package main

import (
	"bytes"
	"fmt"
	"strconv"
)

func encodeInt(i int) []byte {
	return []byte(fmt.Sprintf("i%de", i))

}

func encodeString(s string) []byte {
	return []byte(fmt.Sprintf("%d:%s", len(s), s))
}

func encodeByte(b []byte) []byte {
	return append([]byte(fmt.Sprintf("%d:", len(b))), b...)
}

func encodeList(list []any) []byte {
	toRet := []byte("l")

	for _, item := range list {
		toRet = append(toRet, encode(item)...)
	}

	return append(toRet, byte('e'))

}

func encode(item any) []byte {
	switch item := item.(type) {
	case int:
		return encodeInt(item)
	case string:
		return encodeString(item)
	case []byte:
		return encodeByte(item)

	case []any:
		return encodeList(item)

	case map[any]any:
		return encodeDict(item)

	default:
		// Handle other types as needed
		return nil
	}
}

func encodeDict(dict map[any]any) []byte {
	toRet := []byte("d")

	for key, value := range dict {
		encoded_key := encode(key)
		encoded_value := encode(value)

		toRet = append(toRet, encoded_key...)
		toRet = append(toRet, encoded_value...)
	}

	return append(toRet, byte('e'))

}

func decodeStrOrBytes(bencoded []byte) ([]byte, int) {
	len_index := bytes.Index(bencoded, []byte{':'})
	digits_in_bencoded, error := strconv.Atoi(string(bencoded[:len_index]))
	if error != nil {
		return []byte{}, 0
	}
	return bencoded[len_index+1 : len_index+1+digits_in_bencoded], len_index + 1 + digits_in_bencoded
}

func decodeInt(bencoded []byte) (int, int) {
	first_end_occurence := bytes.Index(bencoded, []byte{'e'})
	num_as_str := string(bencoded[1:first_end_occurence])
	num, error := strconv.Atoi(num_as_str)
	if error != nil {
		return -1, 0
	}

	return num, first_end_occurence + 1
}

func decodeDict(bencoded []byte) (map[any]any, int) {
	dict := make(map[any]any)
	index := 1

	for index < len(bencoded) {
		if bencoded[index] == 'e' {
			index++
			break
		}

		decoded_key, index_to_update := decode(bencoded[index:])
		index += index_to_update
		decoded_val, index_to_update := decode(bencoded[index:])
		index += index_to_update
		bytesKey, ok := decoded_key.([]byte)
		if ok {
			decoded_key = string(bytesKey)
		}

		dict[decoded_key] = decoded_val
	}

	return dict, index
}

func decode(bencoded []byte) (any, int) {
	switch bencoded[0] {
	case 'i':
		return decodeInt(bencoded)
	case 'l':
		return decodeList(bencoded)
	case 'd':
		return decodeDict(bencoded)
	default:
		return decodeStrOrBytes(bencoded)
	}
}

func decodeList(bencoded []byte) ([]any, int) {
	ls := make([]any, 0)
	index := 1
	for index < len(bencoded) {
		current := bencoded[index]
		if current == 'e' {
			index++
			break
		}
		decoded_val, index_to_update := decode(bencoded[index:])
		index += index_to_update
		ls = append(ls, decoded_val)
	}

	return ls, index

}

// func main() {
// 	// dict := map[any]any{
// 	// 	"test": "lol",
// 	// }

// 	i := []any{"123", 19, 123, "sdfdsf"}
// 	encoded := encode(i)
// 	decoded, _ := decode(encoded)
// 	fmt.Println(decoded)

// }
