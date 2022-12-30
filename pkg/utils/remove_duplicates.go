package utils

func RemoveDuplicatesUrls(urls []string) []string {
	keys := make(map[string]interface{})
	list := []string{}
	for _, entry := range urls {
		if _, exist := keys[entry]; !exist {
			keys[entry] = nil
			list = append(list, entry)
		}
	}
	return list
}
