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

func GetErrorUrlsWithoutNotFound(errors map[string]string) map[string]string {
	trueErrors := make(map[string]string)
	for url, err := range errors {
		if err != "Not Found" {
			trueErrors[url] = err
		}
	}
	return trueErrors
}

func MapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
