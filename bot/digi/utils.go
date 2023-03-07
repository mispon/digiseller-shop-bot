package digi

import "fmt"

func ProductImageUrl(id string) string {
	return fmt.Sprintf(imageUrlTempl, id)
}
