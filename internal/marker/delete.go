package marker

import "fmt"

func Delete(key string, markers map[string]string) error {

	if _, ok := markers[key]; !ok {
		return fmt.Errorf("Marker <%s> does not exist.", key)	
	}

	delete(markers, key)
	return SaveMarkers(markers)
}
