package parser

import (
	"fmt"
	"log"
)

type ConcatenateProps struct {
	Tags      []string
	ExportAs  string
	Delimiter string
	AsList    bool
}

func ValidateConcatenateProps(props *map[string]interface{}) (*ConcatenateProps, error) {
	concatenateTags, ok := (*props)["TAGS"].([]string)
	if ok {
		delimiter, ok := (*props)["DELIMITER"].(string)
		if !ok {
			log.Printf("[applyFunctionMapping] 'DELIMITER' value '%s' is empty or not a string.\n", delimiter)
			return nil, fmt.Errorf("'DELIMITER' value is empty of not a string '%s'", delimiter)
		}
		exportedKeyName, ok := (*props)["EXPORTAS"].(string)
		if !ok {
			log.Printf("[applyFunctionMapping] 'EXPORTAS' value '%s' is empty or not a string.\n", exportedKeyName)
			return nil, fmt.Errorf("'EXPORTAS' value is empty of not a string '%s'", exportedKeyName)
		}
		asListFlag, ok := (*props)["ASLIST"].(bool)
		if !ok {
			log.Printf("[applyFunctionMapping] 'ASLIST' value is empty or not a boolean.\n")
			return nil, fmt.Errorf("'ASLIST' value is empty or not a boolean")
		}
		return &ConcatenateProps{Tags: concatenateTags, ExportAs: exportedKeyName, Delimiter: delimiter, AsList: asListFlag}, nil
	} else {
		log.Printf("[applyFunctionMapping] 'CONCATENATE' value '%s' is empty or not a list of strings\n", concatenateTags)
		return nil, fmt.Errorf("'CONCATENATE' is empty or not a list of strings")
	}
}
