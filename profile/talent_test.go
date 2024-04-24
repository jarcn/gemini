package profile

import (
	"encoding/json"
	"fmt"
	"gemini/utils"
	"github.com/buger/jsonparser"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestParser(t *testing.T) {
	root := "/Users/tarzan/Desktop/ekts"
	//root := "/Users/tarzan/Documents/ekt"
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("Error reading file:", err)
			}
			parseData(content)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %v: %v\n", root, err)
	}
}

func parseData(content []byte) {
	talents := make([]Talent, 10)
	jsonparser.ArrayEach(content, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var data Talent
		json.Unmarshal(value, &data)
		formatData(&data)
		if data.Id == 758720 {
			marshal, _ := json.Marshal(data)
			fmt.Println(string(marshal))
		}
		talents = append(talents, data)
	}, "data")
	save2Es(talents)
}

func formatData(data *Talent) {
	data.BirthDate = utils.DateFormat(data.BirthDate)
	educations := data.Educations
	for i := 0; i < len(educations); i++ {
		education := &educations[i]
		education.StartDate = utils.DateFormat(education.StartDate)
		education.EndDate = utils.DateFormat(education.EndDate)
	}
	experiences := data.Experiences
	for i := 0; i < len(experiences); i++ {
		experience := &experiences[i]
		experience.StartDate = utils.DateFormat(experience.StartDate)
		experience.EndDate = utils.DateFormat(experience.EndDate)
		experience.Description = removeHTMLTags(experience.Description)
	}

}

func removeHTMLTags(text string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(text, "")
}

func save2Es(talents []Talent) {

}
