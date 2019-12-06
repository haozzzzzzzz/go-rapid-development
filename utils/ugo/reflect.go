package ugo

import (
	"reflect"
	"strings"
)

func StructTags(
	obj interface{},
	targetTagKeys []string,
) (
	tagsMap map[string][][]string, // target_tag_key -> tags[tag_parts]
) {
	rType := reflect.TypeOf(obj)
	numField := rType.NumField()
	tagsMap = make(map[string][][]string, 0)
	for i := 0; i < numField; i++ {
		tag := rType.Field(i).Tag

		for _, targetTagKey := range targetTagKeys {
			_, ok := tagsMap[targetTagKey]
			if !ok {
				tagsMap[targetTagKey] = make([][]string, 0)
			}

			targetTag := tag.Get(targetTagKey)
			splitParts := strings.Split(targetTag, ",")
			targetTagParts := make([]string, 0)
			for _, splitPart := range splitParts {
				splitPart = strings.TrimSpace(splitPart)
				if splitPart == "" {
					continue
				}

				targetTagParts = append(targetTagParts, splitPart)
			}

			if len(targetTagParts) == 0 {
				continue
			}

			tagsMap[targetTagKey] = append(tagsMap[targetTagKey], targetTagParts)
		}

	}

	return
}

func StructTagNames(obj interface{}, tagKey string) (names []string) {
	names = make([]string, 0)
	tagsMap := StructTags(obj, []string{tagKey})
	tags, ok := tagsMap[tagKey]
	if !ok {
		return
	}

	for _, tag := range tags {
		if len(tag) == 0 {
			continue
		}

		name := tag[0]
		if name == "-" {
			continue
		}

		names = append(names, name)
	}

	return
}

func StructJsonTagNames(obj interface{}) (names []string) {
	return StructTagNames(obj, "json")
}

func StructDbTagNames(obj interface{}) (names []string) {
	return StructTagNames(obj, "db")
}
