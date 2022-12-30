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

func GetErrorUrlsWithoutNotFound(errors map[string]string) []string {
	var urls []string
	for url, err := range errors {
		if err != "Not Found" {
			urls = append(urls, url)
		}
	}
	return urls
}
