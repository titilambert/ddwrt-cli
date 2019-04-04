package cmd

import (
	"bytes"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

type SubSection struct {
	Name        string
	CommandName string
	Link        string
	FormLink    string
}

type Section struct {
	Name        string
	CommandName string
	SubSections map[string]SubSection
}

func fetchTranslation() (map[string]string, error) {
	translation := make(map[string]string)

	body, err := fetchWebPage("/lang_pack/english.js")
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	lines := bytes.Split(body, []byte(";"))
	for _, line := range lines {
		if bytes.HasPrefix(line, []byte("var ")) {
			continue
		}
		key_value := bytes.Split(line, []byte("="))
		if len(key_value) != 2 {
			continue
		}
		key := string(key_value[0])
		value := string(key_value[1])
		value = strings.TrimPrefix(value, "\"")
		value = strings.TrimSuffix(value, "\"")
		translation[key] = value
	}
	return translation, nil
}

func fetchWebPage(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", "http://"+hostAddress+path, nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.SetBasicAuth(username, password)

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	return body, nil
}

func postForm(path string, data map[string]string) error {
	form := url.Values{}
	for key, value := range data {
		form.Add(key, value)
	}
	req, err := http.NewRequest("POST", "http://"+hostAddress+path,
		strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Add("Referer", "http://"+hostAddress+path)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(username, password)

	client := &http.Client{Timeout: time.Second * 10}

	_, err = client.Do(req)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	return nil
}

func getSections() (map[string]Section, error) {
	sections := make(map[string]Section)

	translations, err := fetchTranslation()
	body, err := fetchWebPage("/")
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	doc, err := htmlquery.Parse(bytes.NewReader(body))
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	// List Main Sections

	section_nodes := htmlquery.Find(doc, "//ul[@id=\"menuMainList\"]/li/a")
	for _, section_node := range section_nodes {
		section_name := _getSectionName(translations, section_node)
		commandName := _getCliCommandName(section_name)
		section := Section{Name: section_name,
			CommandName: commandName}
		section.SubSections = make(map[string]SubSection)

		section_link := htmlquery.SelectAttr(section_node, "href")

		sections[commandName] = section

		subbody, err := fetchWebPage("/" + section_link)
		if err != nil {
			log.Fatal("Error reading response. ", err)
		}
		subdoc, err := htmlquery.Parse(bytes.NewReader(subbody))
		if err != nil {
			log.Fatal("Error reading response. ", err)
		}
		// Get Form post url
		form_link_node := htmlquery.FindOne(subdoc, "//form")
		form_link := htmlquery.SelectAttr(form_link_node, "action")
		// List Sub Sections
		subsection_nodes := htmlquery.Find(subdoc, "//ul[@id=\"menuSubList\"]/li")
		for _, subsection_node := range subsection_nodes {

			subsection_a_node := htmlquery.FindOne(subsection_node, "//a")
			if subsection_a_node != nil {
				subsection_name := _getSectionName(translations, subsection_a_node)
				subCommandName := _getCliCommandName(subsection_name)
				subsection_link := htmlquery.SelectAttr(subsection_a_node, "href")
				// Get Form post URL for the current sub section
				subsectionbody, err := fetchWebPage("/" + subsection_link)
				if err != nil {
					log.Fatal("Error reading response. ", err)
				}
				subsectiondoc, err := htmlquery.Parse(bytes.NewReader(subsectionbody))
				if err != nil {
					log.Fatal("Error reading response. ", err)
				}
				subform_link_node := htmlquery.FindOne(subsectiondoc, "//form")
				subform_link := htmlquery.SelectAttr(subform_link_node, "action")
				// Create subsection object
				subsection := SubSection{Name: subsection_name,
					CommandName: subCommandName,
					Link:        subsection_link,
					FormLink:    subform_link}

				sections[commandName].SubSections[subCommandName] = subsection
			} else {
				subsection_span_node := htmlquery.FindOne(subsection_node, "//span")
				subsection_name := _getSectionName(translations, subsection_span_node)
				subCommandName := _getCliCommandName(subsection_name)
				subsection := SubSection{Name: subsection_name,
					CommandName: subCommandName,
					Link:        section_link,
					FormLink:    form_link}
				sections[commandName].SubSections[subCommandName] = subsection
			}
		}
	}
	return sections, err
}

func _getSectionName(translations map[string]string, section_node *html.Node) string {
	raw_section := htmlquery.InnerText(section_node)
	raw_section = strings.TrimPrefix(raw_section, "Capture(")
	raw_section = strings.TrimSuffix(raw_section, ")")
	section := translations[raw_section]
	return section
}

func _getCliCommandName(RawCommandName string) string {
	CommandName := strings.Replace(RawCommandName, " / ", "-", -1)
	CommandName = strings.Replace(CommandName, " ", "-", -1)
	CommandName = strings.ToLower(CommandName)
	return CommandName
}

func getWebPageData(path string, translations map[string]string) (map[string]map[string]string, error) {
	body, err := fetchWebPage(path)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}

	//data := make(map[string]map[string]string)
	data := make(map[string]map[string]string)
	doc, err := htmlquery.Parse(bytes.NewReader(body))
	fieldsets := htmlquery.Find(doc, "//fieldset")
	for _, fieldset := range fieldsets {
		section_node := htmlquery.FindOne(fieldset, "//legend")
		section := _getSectionName(translations, section_node)
		data[section] = make(map[string]string)

        // TODO handle mulple choice combo box
		// handle simple combo box
        select_nodes := htmlquery.Find(fieldset, "//select")
        for _, input := range select_nodes {
            var key = htmlquery.SelectAttr(input, "name")
            option := htmlquery.FindOne(input, "//option[@selected]")
            var value = htmlquery.SelectAttr(option, "value")
            data[section][key] = value
        }
		// Handle radio button and checkbox
		radio_nodes := htmlquery.Find(fieldset, "//input[@checked=\"checked\"]")
		for _, input := range radio_nodes {
			switch htmlquery.SelectAttr(input, "type") {
			case
				"radio",
				"checkbox":
				var key = htmlquery.SelectAttr(input, "name")
				var value = htmlquery.SelectAttr(input, "value")
				data[section][key] = value
			}
		}
		// Handle all other inputs
		input_nodes := htmlquery.Find(fieldset, "//input")
		for _, input := range input_nodes {
			switch htmlquery.SelectAttr(input, "type") {
			case
				"radio",
				"checkbox":
				continue
			}
			var key = htmlquery.SelectAttr(input, "name")
			var value = htmlquery.SelectAttr(input, "value")
			data[section][key] = value
		}
		// Handle textarea
		textarea_nodes := htmlquery.Find(fieldset, "//textarea")
		for _, input := range textarea_nodes {
			var key = htmlquery.SelectAttr(input, "name")
			raw_value := htmlquery.InnerText(input.NextSibling)
			re := regexp.MustCompile(`\( '(.*)' \)`)
			values := re.FindSubmatch([]byte(raw_value))
			value := values[1]
			data[section][key] = string(value)
		}
	}
	// Hidden inputs
	section := "hidden"
	data[section] = make(map[string]string)
	hidden_input_nodes := htmlquery.Find(doc, "//input[@type=\"hidden\"]")
	for _, input := range hidden_input_nodes {
		var key = htmlquery.SelectAttr(input, "name")
		var value = htmlquery.SelectAttr(input, "value")
		data[section][key] = value
	}
	// Footer Inputs
	/* Seems useless
	   footer_node := htmlquery.FindOne(doc, "//div[@class=\"submitFooter\"]")
	   input_nodes := htmlquery.Find(footer_node, "//input")
	   section = "footer"
	   data[section] = make(map[string]string)
	   for _, input := range input_nodes {
	           var key = htmlquery.SelectAttr(input, "name")
	           var value = htmlquery.SelectAttr(input, "value")
	           data[section][key] = value
	       }
	*/
	return data, nil
}

func toYaml(data map[string]map[string]string) ([]byte, error) {
	dataYaml, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}
	return dataYaml, nil
}

func dumpData(rawData map[string]map[string]string, section string, subsection string) (string, error) {
	dataYaml, err := toYaml(rawData)
	if err != nil {
		return "", err
	}
	dumpFileFolder := dumpFolder + "/" + section
	os.MkdirAll(dumpFileFolder, os.ModePerm)
	dumpFilePath := dumpFileFolder + "/" + subsection + ".yml"
	err = ioutil.WriteFile(dumpFilePath, dataYaml, 0644)
	if err != nil {
		return "", err
	}
	return dumpFilePath, nil
}

func loadData(sectionName, subSectionName string) map[string]string {
	// TO DO yaml load from file
	loadFile := path.Join(loadFolder, sectionName, subSectionName+".yml")
	raw_data, err := ioutil.ReadFile(loadFile)
	if err != nil {
		os.Exit(1)
	}
	fieldsets := make(map[string]map[string]string)
	err = yaml.Unmarshal([]byte(raw_data), &fieldsets)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	data := make(map[string]string)
	for _, inputs := range fieldsets {
		for key, value := range inputs {
			data[key] = value
		}
	}
	return data
}

func RunSection(cmd *cobra.Command, args []string) {
	translations, err := fetchTranslation()
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	cliConfig := ReadCommands()
	subSection := cliConfig.Sections[cmd.Parent().Name()].SubSections[cmd.Name()]
	if dumpFolder != "" {
		data, err := getWebPageData("/"+subSection.Link, translations)
		if err != nil {
			log.Fatal("Error reading request. ", err)
		}
		dumpFilePath, err := dumpData(data, cmd.Parent().Name(), cmd.Name())
		if err != nil {
			log.Fatal("Error reading request. ", err)
		}
		fmt.Printf("%s => %s dumped to %s\n",
			cmd.Parent().Name(), cmd.Name(), dumpFilePath)
	} else if loadFolder != "" {
		data := loadData(cmd.Parent().Name(), cmd.Name())
		postForm("/"+subSection.FormLink, data)
	} else {
		data, err := getWebPageData("/"+subSection.Link, translations)
		if err != nil {
			log.Fatal("Error reading request. ", err)
		}
		dataYaml, err := toYaml(data)
		if err != nil {
			log.Fatal("Error reading request. ", err)
		}
		fmt.Println("########################")
		fmt.Printf("# %s => %s\n", cmd.Parent().Name(), cmd.Name())
		fmt.Println("########################")
		fmt.Printf("\n%s\n", dataYaml)
	}
}

func PreRunSection(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		fmt.Print("Sections file doesn't exist, please run `init` command.\n")
		os.Exit(1)
	}
	cliConfig := ReadCommands()
	if cliConfig.Host != hostAddress {
		fmt.Print("Sections file not updated, please run `init` command.\n")
		os.Exit(0)
	}
}
