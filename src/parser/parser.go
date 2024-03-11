package parser

import (
	"fmt"
	"log"
	"strings"

	"github.com/beevik/etree"
)

type SimpleMapping map[string]string
type FunctionMapping map[string]map[string]interface{}
type ParsedData map[string]interface{}

type XMLConfig struct {
	RootElement     string
	JobElement      string
	SimpleMapping   SimpleMapping
	FunctionMapping FunctionMapping
	Data            string
}

func applySimpleMapping(job *etree.Element, xml *XMLConfig, pd *ParsedData) error {
	// Loop through SimpleMapping exported tags (et) and find them in the job el
	for et, tagName := range xml.SimpleMapping {
		tag := job.SelectElement(tagName)
		if tag == nil {
			log.Printf("[ParseXML] '%s' not found in '<%s>' element.\n", tagName, xml.JobElement)
			return fmt.Errorf("'%s' not found in '<%s>' element", tagName, xml.JobElement)
		}
		(*pd)[et] = strings.TrimSpace(tag.Text()) // Add the exported key "et" with the value of the tag in source data
	}
	return nil
}

func applyFunctionMapping(job *etree.Element, xml *XMLConfig, pd *ParsedData) error {
	// Loop through FunctionMapping exported tags (fet)
	for fnName, props := range xml.FunctionMapping {
		switch fnName {
		case "CONCATENATE":
			cp, err := ValidateConcatenateProps(&props)
			if err != nil {
				return err
			}
			concatenatedValues := ""
			// Loop each concatenate tag
			for i, tagName := range cp.Tags {
				// Find the tag in the main wrapper element
				tag := job.SelectElement(tagName)
				if tag == nil {
					log.Printf("[applyFunctionMapping] '%s' not found in '<%s>' element.\n", tagName, xml.JobElement)
					return fmt.Errorf("'%s' not found in '<%s>' element", tagName, xml.JobElement)
				}
				// Concatenate its value
				concatenatedValues += strings.TrimSpace(tag.Text())
				// Don't add a delimiter if it's the last tag
				if i < len(cp.Tags)-1 {
					concatenatedValues += cp.Delimiter
				}
			}
			// Add concatenated value to parsed data.
			(*pd)[cp.ExportAs] = concatenatedValues
		}
	}
	return nil
}

func ParseXML(xml XMLConfig) ([]ParsedData, error) {
	/*
		This will be the XML or JSON parser.
		Main functionality is to parse the data from the importer into Go structs.
	*/
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml.Data); err != nil {
		log.Printf("[ParseXML] Error parsing XML =%s\n", err.Error())
		return nil, err
	}
	pdList := []ParsedData{}
	// Find root el
	root := doc.SelectElement(xml.RootElement)
	if root == nil {
		log.Printf("[ParseXML] Root element '%s' not found.\n", xml.RootElement)
		return nil, fmt.Errorf("root element '%s' not found", xml.RootElement)
	}
	// Find "job" element (corresponds to 1 job data entry)
	jobElement := root.SelectElements(xml.JobElement)
	if jobElement == nil {
		log.Printf("[ParseXML] Job element '%s' not found.\n", xml.JobElement)
		return nil, fmt.Errorf("job element '%s' not found", xml.JobElement)
	}
	// Validate all jobs contain SimpleMapping/FunctionMapping tags and return error if not
	// Loop through each "job" element
	for jobIteration, job := range jobElement {
		log.Printf("[ParseXML] Processing job iteration %d\n", jobIteration)
		// Parsed data struct
		pd := ParsedData{}
		// Apply simple mapping
		smErr := applySimpleMapping(job, &xml, &pd)
		if smErr != nil {
			return nil, smErr
		}
		// Apply function mapping
		fmErr := applyFunctionMapping(job, &xml, &pd)
		if fmErr != nil {
			return nil, fmErr
		}
		// Append processed (pd) data to list of processed data (pdList)
		pdList = append(pdList, pd)
	}
	return pdList, nil
}

// #TODO handle JSON parsing from source
